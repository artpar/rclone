package gendocs

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/artpar/rclone/cmd"
	"github.com/artpar/rclone/lib/file"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

func init() {
	cmd.Root.AddCommand(commandDefinition)
}

// define things which go into the frontmatter
type frontmatter struct {
	Date        string
	Title       string
	Description string
	Slug        string
	URL         string
	Source      string
}

var frontmatterTemplate = template.Must(template.New("frontmatter").Parse(`---
title: "{{ .Title }}"
description: "{{ .Description }}"
slug: {{ .Slug }}
url: {{ .URL }}
# autogenerated - DO NOT EDIT, instead edit the source code in {{ .Source }} and as part of making a release run "make commanddocs"
---
`))

var commandDefinition = &cobra.Command{
	Use:   "gendocs output_directory",
	Short: `Output markdown docs for rclone to the directory supplied.`,
	Long: `
This produces markdown docs for the rclone commands to the directory
supplied.  These are in a format suitable for hugo to render into the
rclone.org website.`,
	RunE: func(command *cobra.Command, args []string) error {
		cmd.CheckArgs(1, 1, command, args)
		now := time.Now().Format(time.RFC3339)

		// Create the directory structure
		root := args[0]
		out := filepath.Join(root, "commands")
		err := file.MkdirAll(out, 0777)
		if err != nil {
			return err
		}

		// Write the flags page
		var buf bytes.Buffer
		cmd.Root.SetOutput(&buf)
		cmd.Root.SetArgs([]string{"help", "flags"})
		cmd.GeneratingDocs = true
		err = cmd.Root.Execute()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(root, "flags.md"), buf.Bytes(), 0777)
		if err != nil {
			return err
		}

		// Look up name => description for prepender
		var description = map[string]string{}
		var addDescription func(root *cobra.Command)
		addDescription = func(root *cobra.Command) {
			name := strings.Replace(root.CommandPath(), " ", "_", -1) + ".md"
			description[name] = root.Short
			for _, c := range root.Commands() {
				addDescription(c)
			}
		}
		addDescription(cmd.Root)

		// markup for the docs files
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			data := frontmatter{
				Date:        now,
				Title:       strings.Replace(base, "_", " ", -1),
				Description: description[name],
				Slug:        base,
				URL:         "/commands/" + strings.ToLower(base) + "/",
				Source:      strings.Replace(strings.Replace(base, "rclone", "cmd", -1), "_", "/", -1) + "/",
			}
			var buf bytes.Buffer
			err := frontmatterTemplate.Execute(&buf, data)
			if err != nil {
				log.Errorf("Failed to render frontmatter template: %v", err)
			}
			return buf.String()
		}
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		// Hide all of the root entries flags
		cmd.Root.Flags().VisitAll(func(flag *pflag.Flag) {
			flag.Hidden = true
		})
		err = doc.GenMarkdownTreeCustom(cmd.Root, out, prepender, linkHandler)
		if err != nil {
			return err
		}

		var outdentTitle = regexp.MustCompile(`(?m)^#(#+)`)

		// Munge the files to add a link to the global flags page
		err = filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				doc := string(b)
				doc = strings.Replace(doc, "\n### SEE ALSO", `
See the [global flags page](/flags/) for global options not listed here.

### SEE ALSO`, 1)
				// outdent all the titles by one
				doc = outdentTitle.ReplaceAllString(doc, `$1`)
				err = ioutil.WriteFile(path, []byte(doc), 0777)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}
