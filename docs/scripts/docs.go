package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/renehernandez/appfile/cmd"
	"github.com/spf13/cobra"
)

func generateCommandsReference(buffer *bytes.Buffer) {
	cmd := cmd.NewRootCmd()

	buffer.WriteString("# CLI Reference\n\n")
	buffer.WriteString("This is a reference for the `appfile` CLI,")
	buffer.WriteString(" which enables you to manage deployments to DigitalOcean App Platform.\n\n")
	buffer.WriteString("## Command List\n\n")
	buffer.WriteString("The following is a complete list of the commands provided by `appfile`.\n\n")
	buffer.WriteString("Command | Description \n")
	buffer.WriteString("---- | ---- \n")

	for _, c := range cmd.Commands() {
		buffer.WriteString(fmt.Sprintf("[%s](#appfile-%s) | %s\n", c.Name(), c.Name(), c.Short))
	}

	buffer.WriteString("\n")
}

func generateCommandsHelp(buffer *bytes.Buffer) {
	cmd := cmd.NewRootCmd()
	cmd.DisableAutoGenTag = true

	GenMarkdown(cmd, buffer)

	for _, c := range cmd.Commands() {
		c.DisableAutoGenTag = true
		GenMarkdown(c, buffer)
	}
}

// GenMarkdown creates markdown output.
func GenMarkdown(cmd *cobra.Command, w io.Writer) error {
	return GenMarkdownCustom(cmd, w, func(s string) string { return s })
}

// GenMarkdownCustom creates custom markdown output.
func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}

	if !cmd.DisableAutoGenTag {
		buf.WriteString("###### Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "\n")
	}
	_, err := buf.WriteTo(w)
	return err
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

func generateCliReference() {
	buffer := &bytes.Buffer{}
	generateCommandsReference(buffer)
	generateCommandsHelp(buffer)

	f, err := os.Create("./docs/cli_reference.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = buffer.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
}

func generateRedirect() {
	version := os.Getenv("APPFILE_VERSION")

	content, err := ioutil.ReadFile("./docs/scripts/index.html.gotmpl")

	if err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir("./docs/redirect", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./docs/redirect/index.html")

	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.New("redirect").Parse(string(content)))

	tmpl.Execute(f, version)
}

func main() {
	generateCliReference()
	generateRedirect()
}
