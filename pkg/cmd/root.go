package cmd

import (
	"fmt"
	"os"

	"github.com/lshgdut/repoctl/pkg/cmd/count"
	init_cmd "github.com/lshgdut/repoctl/pkg/cmd/init"
	"github.com/lshgdut/repoctl/pkg/cmd/project"
	"github.com/lshgdut/repoctl/pkg/cmd/release"
	"github.com/lshgdut/repoctl/pkg/cmd/tag"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repoctl",
	Short: "Repositories management tool for gitlab",
	Long:  `A tool to manage repositories in gitlab`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(init_cmd.NewInitCmd())
	rootCmd.AddCommand(project.NewProjectCmd())
	rootCmd.AddCommand(release.NewReleaseCmd())
	rootCmd.AddCommand(tag.NewTagCmd())
	rootCmd.AddCommand(count.NewCountCmd())
}
