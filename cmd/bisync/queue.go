package bisync

import (
	"context"
	"fmt"

	"github.com/artpar/rclone/cmd/bisync/bilib"
	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fs/filter"
	"github.com/artpar/rclone/fs/operations"
	"github.com/artpar/rclone/fs/sync"
)

func (b *bisyncRun) fastCopy(ctx context.Context, fsrc, fdst fs.Fs, files bilib.Names, queueName string) error {
	if err := b.saveQueue(files, queueName); err != nil {
		return err
	}

	ctxCopy, filterCopy := filter.AddConfig(b.opt.setDryRun(ctx))
	for _, file := range files.ToList() {
		if err := filterCopy.AddFile(file); err != nil {
			return err
		}
	}

	return sync.CopyDir(ctxCopy, fdst, fsrc, false)
}

func (b *bisyncRun) fastDelete(ctx context.Context, f fs.Fs, files bilib.Names, queueName string) error {
	if err := b.saveQueue(files, queueName); err != nil {
		return err
	}

	transfers := fs.GetConfig(ctx).Transfers
	ctxRun := b.opt.setDryRun(ctx)

	objChan := make(fs.ObjectsChan, transfers)
	errChan := make(chan error, 1)
	go func() {
		errChan <- operations.DeleteFiles(ctxRun, objChan)
	}()
	err := operations.ListFn(ctxRun, f, func(obj fs.Object) {
		remote := obj.Remote()
		if files.Has(remote) {
			objChan <- obj
		}
	})
	close(objChan)
	opErr := <-errChan
	if err == nil {
		err = opErr
	}
	return err
}

func (b *bisyncRun) saveQueue(files bilib.Names, jobName string) error {
	if !b.opt.SaveQueues {
		return nil
	}
	queueFile := fmt.Sprintf("%s.%s.que", b.basePath, jobName)
	return files.Save(queueFile)
}
