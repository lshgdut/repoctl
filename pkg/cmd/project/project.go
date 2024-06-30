package project

import (
	"fmt"

	cmdutil "github.com/lshgdut/repoctl/pkg/cmd/util"
	_gitlab "github.com/lshgdut/repoctl/pkg/gitlab"
	"github.com/lshgdut/repoctl/pkg/gitlab/api"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ProjectOptions struct {
	client *gitlab.Client

	args []string
}

func NewProjectOptions() *ProjectOptions {
	return &ProjectOptions{}
}

func NewProjectCmd() *cobra.Command {
	o := NewProjectOptions()

	var cmd = &cobra.Command{
		Use:                   "projects",
		DisableFlagsInUseLine: true,
		Short:                 "list projects",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

func (o *ProjectOptions) Complete(cmd *cobra.Command, args []string) error {
	// var err error
	o.args = args

	// Load config
	client, err := _gitlab.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create gitlab client: %v", err)
	}
	o.client = client

	return nil
}

func (o *ProjectOptions) Validate() error {
	return nil
}

func (o *ProjectOptions) Run() error {
	projects, err := api.ProjectList(o.client, &gitlab.ListProjectsOptions{})

	if err != nil {
		return fmt.Errorf("failed to list projects: %v", err)
	}

	for _, project := range projects {
		fmt.Println(project.PathWithNamespace)
	}

	return nil
}
