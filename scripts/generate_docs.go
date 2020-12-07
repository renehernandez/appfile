package main

import (
	"log"

	"github.com/renehernandez/appfile/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(cmd, "docs/cmd")
	if err != nil {
		log.Fatal(err)
	}
}
