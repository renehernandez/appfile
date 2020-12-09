package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/renehernandez/appfile/cmd"
	"github.com/spf13/cobra/doc"
)

func generateCommandsReference(filepath string) error {
	cmd := cmd.NewRootCmd()

	buf := bytes.Buffer{}
	buf.WriteString("# CLI Reference for appfile\n")
	sb := strings.Builder{}
	sb.WriteString("This is a reference for the `appfile` CLI,")
	sb.WriteString(" which enables you to manage deployments to DigitalOcean App Platform.\n")
	sb.WriteString("## Command List\n")
	sb.WriteString("The following is a complete list of the commands provided by `appfile`.\n")
	sb.WriteString("Command | Description\n")
	sb.WriteString("- | -\n")

	for _, c := range cmd.Commands() {
		sb.WriteString(fmt.Sprintf("[%s](appfile_%s.md) | %s\n", c.Name(), c.Name(), c.Short))
	}

	buf.WriteString(sb.String())

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = buf.WriteTo(f)
	return err
}

func main() {
	cmd := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(cmd, "docs/cmd")
	if err != nil {
		log.Fatal(err)
	}

	err = generateCommandsReference("docs/cmd/reference.md")
	if err != nil {
		log.Fatal(err)
	}
}
