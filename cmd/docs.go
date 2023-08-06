package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"io"
	"log"
	"n0rdy.me/remindme/common"
	"os"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for remindme command",
	Long:  "Generate documentation for remindme command.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("docs command: called")

		dir, err := getDir(cmd)
		if err != nil {
			return err
		}
		return generateDocs(os.Stdout, dir)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for the generated documentation")
}

func getDir(cmd *cobra.Command) (string, error) {
	dir, err := cmd.Flags().GetString(common.DirFlag)
	if err != nil {
		log.Println("docs command: error while parsing flag: "+common.DirFlag, err)
		return "", common.ErrWrongFormattedStringFlag(common.DirFlag)
	}

	if dir == "" {
		if dir, err = os.MkdirTemp("", "remindme"); err != nil {
			log.Println("docs command: error while creating temp dir", err)
			return "", common.ErrDocsCmdOnDirCreation
		}
	}
	return dir, nil
}

func generateDocs(out io.Writer, dir string) error {
	if err := doc.GenMarkdownTree(rootCmd, dir); err != nil {
		log.Println("docs command: error while generating docs", err)
		return common.ErrDocsCmdOnDocsGeneration
	}
	_, err := fmt.Fprintf(out, "Documentation successfully created in %s\n", dir)
	if err != nil {
		log.Println("docs command: error while writing to output", err)
	}
	return nil
}
