// Package cmd implements the rclone command
//
// It is in a sub package so it's internals can be reused elsewhere
package cmd

// FIXME only attach the remote flags when using a remote???
// would probably mean bringing all the flags in to here? Or define some flagsets in fs...

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fs/accounting"
	"github.com/artpar/rclone/fs/cache"
	"github.com/artpar/rclone/fs/config/configfile"
	"github.com/artpar/rclone/fs/config/configflags"
	"github.com/artpar/rclone/fs/config/flags"
	"github.com/artpar/rclone/fs/filter"
	"github.com/artpar/rclone/fs/fserrors"
	"github.com/artpar/rclone/fs/fspath"
	fslog "github.com/artpar/rclone/fs/log"
	"github.com/artpar/rclone/fs/rc"
	"github.com/artpar/rclone/fs/rc/rcserver"
	fssync "github.com/artpar/rclone/fs/sync"
	"github.com/artpar/rclone/lib/atexit"
	"github.com/artpar/rclone/lib/buildinfo"
	"github.com/artpar/rclone/lib/exitcode"
	"github.com/artpar/rclone/lib/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Globals
var (
	// Flags
	cpuProfile    = flags.StringP("cpuprofile", "", "", "Write cpu profile to file", "Debugging")
	memProfile    = flags.StringP("memprofile", "", "", "Write memory profile to file", "Debugging")
	statsInterval = flags.DurationP("stats", "", time.Minute*1, "Interval between printing stats, e.g. 500ms, 60s, 5m (0 to disable)", "Logging")
	version       bool
	// Errors
	errorCommandNotFound    = errors.New("command not found")
	errorNotEnoughArguments = errors.New("not enough arguments")
	errorTooManyArguments   = errors.New("too many arguments")
)

// ShowVersion prints the version to stdout
func ShowVersion() {
	osVersion, osKernel := buildinfo.GetOSVersion()
	if osVersion == "" {
		osVersion = "unknown"
	}
	if osKernel == "" {
		osKernel = "unknown"
	}

	linking, tagString := buildinfo.GetLinkingAndTags()

	arch := buildinfo.GetArch()

	fmt.Printf("rclone %s\n", fs.Version)
	fmt.Printf("- os/version: %s\n", osVersion)
	fmt.Printf("- os/kernel: %s\n", osKernel)
	fmt.Printf("- os/type: %s\n", runtime.GOOS)
	fmt.Printf("- os/arch: %s\n", arch)
	fmt.Printf("- go/version: %s\n", runtime.Version())
	fmt.Printf("- go/linking: %s\n", linking)
	fmt.Printf("- go/tags: %s\n", tagString)
}

// NewFsFile creates an Fs from a name but may point to a file.
//
// It returns a string with the file name if points to a file
// otherwise "".
func NewFsFile(remote string) (fs.Fs, string) {
	ctx := context.Background()
	_, fsPath, err := fspath.SplitFs(remote)
	if err != nil {
		err = fs.CountError(ctx, err)
		//fs.Fatalf(nil, "Failed to create file system for %q: %v", remote, err)
	}
	f, err := cache.Get(ctx, remote)
	switch err {
	case fs.ErrorIsFile:
		cache.Pin(f) // pin indefinitely since it was on the CLI
		return f, path.Base(fsPath)
	case nil:
		cache.Pin(f) // pin indefinitely since it was on the CLI
		return f, ""
	default:
		err = fs.CountError(ctx, err)
		//fs.Fatalf(nil, "Failed to create file system for %q: %v", remote, err)
	}
	return nil, ""
}

// newFsFileAddFilter creates an src Fs from a name
//
// This works the same as NewFsFile however it adds filters to the Fs
// to limit it to a single file if the remote pointed to a file.
func newFsFileAddFilter(remote string) (fs.Fs, string) {
	ctx := context.Background()
	fi := filter.GetConfig(ctx)
	f, fileName := NewFsFile(remote)
	if fileName != "" {
		if !fi.InActive() {
			err := fmt.Errorf("can't limit to single files when using filters: %v", remote)
			err = fs.CountError(ctx, err)
			fmt.Printf("%v", err.Error())
		}
		// Limit transfers to this file
		err := fi.AddFile(fileName)
		if err != nil {
			err = fs.CountError(ctx, err)
			fmt.Printf("Failed to limit to single file %q: %v", remote, err)

		}
	}
	return f, fileName
}

// NewFsSrc creates a new src fs from the arguments.
//
// The source can be a file or a directory - if a file then it will
// limit the Fs to a single file.
func NewFsSrc(args []string) fs.Fs {
	fsrc, _ := newFsFileAddFilter(args[0])
	return fsrc
}

// newFsDir creates an Fs from a name
//
// This must point to a directory
func newFsDir(remote string) fs.Fs {
	ctx := context.Background()
	f, err := cache.Get(ctx, remote)
	if err != nil {
		err = fs.CountError(ctx, err)
		fmt.Printf("Failed to create file system for %q: %v", remote, err)
	}
	cache.Pin(f) // pin indefinitely since it was on the CLI
	return f
}

// NewFsDir creates a new Fs from the arguments
//
// The argument must point a directory
func NewFsDir(args []string) fs.Fs {
	fdst := newFsDir(args[0])
	return fdst
}

// NewFsSrcDst creates a new src and dst fs from the arguments
func NewFsSrcDst(args []string) (fs.Fs, fs.Fs) {
	fsrc, _ := newFsFileAddFilter(args[0])
	fdst := newFsDir(args[1])
	return fsrc, fdst
}

// NewFsSrcFileDst creates a new src and dst fs from the arguments
//
// The source may be a file, in which case the source Fs and file name is returned
func NewFsSrcFileDst(args []string) (fsrc fs.Fs, srcFileName string, fdst fs.Fs) {
	fsrc, srcFileName = NewFsFile(args[0])
	fdst = newFsDir(args[1])
	return fsrc, srcFileName, fdst
}

// NewFsSrcDstFiles creates a new src and dst fs from the arguments
// If src is a file then srcFileName and dstFileName will be non-empty
func NewFsSrcDstFiles(args []string) (fsrc fs.Fs, srcFileName string, fdst fs.Fs, dstFileName string) {
	ctx := context.Background()
	fsrc, srcFileName = newFsFileAddFilter(args[0])
	// If copying a file...
	dstRemote := args[1]
	// If file exists then srcFileName != "", however if the file
	// doesn't exist then we assume it is a directory...
	if srcFileName != "" {
		var err error
		dstRemote, dstFileName, err = fspath.Split(dstRemote)
		if err != nil {
			fmt.Printf("Parsing %q failed: %v", args[1], err)
		}
		if dstRemote == "" {
			dstRemote = "."
		}
		if dstFileName == "" {
			fmt.Printf("%q is a directory", args[1])
		}
	}
	fdst, err := cache.Get(ctx, dstRemote)
	switch err {
	case fs.ErrorIsFile:
		_ = fs.CountError(ctx, err)
		fmt.Printf("Source doesn't exist or is a directory and destination is a file")
	case nil:
	default:
		_ = fs.CountError(ctx, err)
		fmt.Printf("Failed to create file system for destination %q: %v", dstRemote, err)
	}
	cache.Pin(fdst) // pin indefinitely since it was on the CLI
	return
}

// NewFsDstFile creates a new dst fs with a destination file name from the arguments
func NewFsDstFile(args []string) (fdst fs.Fs, dstFileName string) {
	dstRemote, dstFileName, err := fspath.Split(args[0])
	if err != nil {
		fmt.Printf("Parsing %q failed: %v", args[0], err)
	}
	if dstRemote == "" {
		dstRemote = "."
	}
	if dstFileName == "" {
		fmt.Printf("%q is a directory", args[0])
	}
	fdst = newFsDir(dstRemote)
	return
}

// ShowStats returns true if the user added a `--stats` flag to the command line.
//
// This is called by Run to override the default value of the
// showStats passed in.
func ShowStats() bool {
	statsIntervalFlag := pflag.Lookup("stats")
	return statsIntervalFlag != nil && statsIntervalFlag.Changed
}

// Run the function with stats and retries if required
func Run(Retry bool, showStats bool, cmd *cobra.Command, f func() error) {
	ctx := context.Background()
	ci := fs.GetConfig(ctx)
	var cmdErr error
	stopStats := func() {}
	if !showStats && ShowStats() {
		showStats = true
	}
	if ci.Progress {
		stopStats = startProgress()
	} else if showStats {
		stopStats = StartStats()
	}
	SigInfoHandler()
	for try := 1; try <= ci.Retries; try++ {
		cmdErr = f()
		cmdErr = fs.CountError(ctx, cmdErr)
		lastErr := accounting.GlobalStats().GetLastError()
		if cmdErr == nil {
			cmdErr = lastErr
		}
		if !Retry || !accounting.GlobalStats().Errored() {
			if try > 1 {
				fs.Errorf(nil, "Attempt %d/%d succeeded", try, ci.Retries)
			}
			break
		}
		if accounting.GlobalStats().HadFatalError() {
			fmt.Printf("Fatal error received - not attempting retries")
			break
		}
		if accounting.GlobalStats().Errored() && !accounting.GlobalStats().HadRetryError() {
			fs.Errorf(nil, "Can't retry any of the errors - not attempting retries")
			break
		}
		if retryAfter := accounting.GlobalStats().RetryAfter(); !retryAfter.IsZero() {
			d := time.Until(retryAfter)
			if d > 0 {
				fs.Logf(nil, "Received retry after error - sleeping until %s (%v)", retryAfter.Format(time.RFC3339Nano), d)
				time.Sleep(d)
			}
		}
		if lastErr != nil {
			fs.Errorf(nil, "Attempt %d/%d failed with %d errors and: %v", try, ci.Retries, accounting.GlobalStats().GetErrors(), lastErr)
		} else {
			fs.Errorf(nil, "Attempt %d/%d failed with %d errors", try, ci.Retries, accounting.GlobalStats().GetErrors())
		}
		if try < ci.Retries {
			accounting.GlobalStats().ResetErrors()
		}
		if ci.RetriesInterval > 0 {
			time.Sleep(ci.RetriesInterval)
		}
	}
	stopStats()
	if showStats && (accounting.GlobalStats().Errored() || *statsInterval > 0) {
		accounting.GlobalStats().Log()
	}
	fs.Debugf(nil, "%d go routines active\n", runtime.NumGoroutine())

	if ci.Progress && ci.ProgressTerminalTitle {
		// Clear terminal title
		terminal.WriteTerminalTitle("")
	}

	// dump all running go-routines
	if ci.Dump&fs.DumpGoRoutines != 0 {
		err := pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		if err != nil {
			fs.Errorf(nil, "Failed to dump goroutines: %v", err)
		}
	}

	// dump open files
	if ci.Dump&fs.DumpOpenFiles != 0 {
		c := exec.Command("lsof", "-p", strconv.Itoa(os.Getpid()))
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			fs.Errorf(nil, "Failed to list open files: %v", err)
		}
	}

	// clear cache and shutdown backends
	cache.Clear()
	if lastErr := accounting.GlobalStats().GetLastError(); cmdErr == nil {
		cmdErr = lastErr
	}

	// Log the final error message and exit
	if cmdErr != nil {
		nerrs := accounting.GlobalStats().GetErrors()
		if nerrs <= 1 {
			fs.Logf(nil, "Failed to %s: %v", cmd.Name(), cmdErr)
		} else {
			fs.Logf(nil, "Failed to %s with %d errors: last error was: %v", cmd.Name(), nerrs, cmdErr)
		}
	}
	resolveExitCode(cmdErr)
}

// CheckArgs checks there are enough arguments and prints a message if not
func CheckArgs(MinArgs, MaxArgs int, cmd *cobra.Command, args []string) {
	if len(args) < MinArgs {
		_ = cmd.Usage()
		_, _ = fmt.Fprintf(os.Stderr, "Command %s needs %d arguments minimum: you provided %d non flag arguments: %q\n", cmd.Name(), MinArgs, len(args), args)
		resolveExitCode(errorNotEnoughArguments)
	} else if len(args) > MaxArgs {
		_ = cmd.Usage()
		_, _ = fmt.Fprintf(os.Stderr, "Command %s needs %d arguments maximum: you provided %d non flag arguments: %q\n", cmd.Name(), MaxArgs, len(args), args)
		resolveExitCode(errorTooManyArguments)
	}
}

// StartStats prints the stats every statsInterval
//
// It returns a func which should be called to stop the stats.
func StartStats() func() {
	if *statsInterval <= 0 {
		return func() {}
	}
	stopStats := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(*statsInterval)
		for {
			select {
			case <-ticker.C:
				accounting.GlobalStats().Log()
			case <-stopStats:
				ticker.Stop()
				return
			}
		}
	}()
	return func() {
		close(stopStats)
		wg.Wait()
	}
}

// initConfig is run by cobra after initialising the flags
func initConfig() {
	// Set the global options from the flags
	err := fs.GlobalOptionsInit()
	if err != nil {
		fmt.Printf("Failed to initialise global options: %v", err)
	}

	ctx := context.Background()
	ci := fs.GetConfig(ctx)

	// Start the logger
	fslog.InitLogging()

	// Finish parsing any command line flags
	configflags.SetFlags(ci)

	// Load the config
	configfile.Install()

	// Start accounting
	accounting.Start(ctx)

	// Configure console
	if ci.NoConsole {
		// Hide the console window
		terminal.HideConsole()
	} else {
		// Enable color support on stdout if possible.
		// This enables virtual terminal processing on Windows 10,
		// adding native support for ANSI/VT100 escape sequences.
		terminal.EnableColorsStdout()
	}

	// Write the args for debug purposes
	fs.Debugf("rclone", "Version %q starting with parameters %q", fs.Version, os.Args)

	// Inform user about systemd log support now that we have a logger
	if fslog.Opt.LogSystemdSupport {
		fs.Debugf("rclone", "systemd logging support activated")
	}

	// Start the remote control server if configured
	_, err = rcserver.Start(ctx, &rc.Opt)
	if err != nil {
		fs.Fatalf(nil, "Failed to start remote control: %v", err)
	}

	// Start the metrics server if configured and not running the "rc" command
	if os.Args[1] != "rc" {
		_, err = rcserver.MetricsStart(ctx, &rc.Opt)
		if err != nil {
			fs.Fatalf(nil, "Failed to start metrics server: %v", err)
		}
	}

	// Setup CPU profiling if desired
	if *cpuProfile != "" {
		fs.Infof(nil, "Creating CPU profile %q\n", *cpuProfile)
		f, err := os.Create(*cpuProfile)
		if err != nil {
			err = fs.CountError(ctx, err)
			fs.Fatal(nil, fmt.Sprint(err))
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			err = fs.CountError(ctx, err)
			fs.Fatal(nil, fmt.Sprint(err))
		}
		atexit.Register(func() {
			pprof.StopCPUProfile()
			err := f.Close()
			if err != nil {
				err = fs.CountError(ctx, err)
				fs.Fatal(nil, fmt.Sprint(err))
			}
		})
	}

	// Setup memory profiling if desired
	if *memProfile != "" {
		atexit.Register(func() {
			fs.Infof(nil, "Saving Memory profile %q\n", *memProfile)
			f, err := os.Create(*memProfile)
			if err != nil {
				err = fs.CountError(ctx, err)
				fs.Fatal(nil, fmt.Sprint(err))
			}
			err = pprof.WriteHeapProfile(f)
			if err != nil {
				err = fs.CountError(ctx, err)
				fs.Fatal(nil, fmt.Sprint(err))
			}
			err = f.Close()
			if err != nil {
				err = fs.CountError(ctx, err)
				fs.Fatal(nil, fmt.Sprint(err))
			}
		})
	}
}

func resolveExitCode(err error) {
	ctx := context.Background()
	ci := fs.GetConfig(ctx)
	atexit.Run()
	if err == nil {
		if ci.ErrorOnNoTransfer {
			if accounting.GlobalStats().GetTransfers() == 0 {
				fmt.Printf("exitcode - %v", exitcode.NoFilesTransferred)
			}
		}
		//fmt.Printf("exitcode - %v", exitcode.Success)
	}

	switch {
	case errors.Is(err, fs.ErrorDirNotFound):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.DirNotFound)
	case errors.Is(err, fs.ErrorObjectNotFound):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.FileNotFound)
	case errors.Is(err, accounting.ErrorMaxTransferLimitReached):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.TransferExceeded)
	case errors.Is(err, fssync.ErrorMaxDurationReached):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.DurationExceeded)
	case fserrors.ShouldRetry(err):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.RetryError)
	case fserrors.IsNoRetryError(err), fserrors.IsNoLowLevelRetryError(err):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.NoRetryError)
	case fserrors.IsFatalError(err):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.FatalError)
	case errors.Is(err, errorCommandNotFound), errors.Is(err, errorNotEnoughArguments), errors.Is(err, errorTooManyArguments):
		fmt.Printf("rclone session exitcode - %v\n", exitcode.UsageError)
	default:
		fmt.Printf("rclone session exitcode - %v\n", exitcode.UncategorizedError)
	}
}

var backendFlags map[string]struct{}

// AddBackendFlags creates flags for all the backend options
func AddBackendFlags() {
	backendFlags = map[string]struct{}{}
	for _, fsInfo := range fs.Registry {
		flags.AddFlagsFromOptions(pflag.CommandLine, fsInfo.Prefix, fsInfo.Options)
		// Store the backend flag names for the help generator
		for i := range fsInfo.Options {
			opt := &fsInfo.Options[i]
			name := opt.FlagName(fsInfo.Prefix)
			backendFlags[name] = struct{}{}
		}
	}
}

// Main runs rclone interpreting flags and commands out of os.Args
func Main() {
	setupRootCommand(Root)
	AddBackendFlags()
	if err := Root.Execute(); err != nil {
		if strings.HasPrefix(err.Error(), "unknown command") && selfupdateEnabled {
			Root.PrintErrf("You could use '%s selfupdate' to get latest features.\n\n", Root.CommandPath())
		}
		fs.Logf(nil, "Fatal error: %v", err)
		fmt.Printf("%v", exitcode.UsageError)
	}
}
