package count

import (
	"fmt"

	cmdutil "github.com/lshgdut/repoctl/pkg/cmd/util"
	_gitlab "github.com/lshgdut/repoctl/pkg/gitlab"
	"github.com/lshgdut/repoctl/pkg/gitlab/api"

	// api "github.com/lshgdut/repoctl/pkg/gitlab/api"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CountOptions struct {
	client *gitlab.Client
	dryRun bool

	args         []string
	repositories []_gitlab.Repository

	BranchName   string
	RepoListFile string
}

func NewCountOptions() *CountOptions {
	return &CountOptions{
		dryRun: false,
	}
}

func NewCountCmd() *cobra.Command {
	o := NewCountOptions()

	var cmd = &cobra.Command{
		Use:                   "count --branch <branch> -f <repolist.txt>",
		DisableFlagsInUseLine: true,
		Short:                 "show branch ahead and behind counts for each repository between main branch and release branch",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.BranchName, "branch", "b", "", "release branch name, e.g. release-1.0.0")
	cmd.Flags().StringVarP(&o.RepoListFile, "repolist-file", "f", "", "repolist file path, e.g. repolist.txt")

	return cmd
}

func (o *CountOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	// Load config
	client, err := _gitlab.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create gitlab client: %v", err)
	}
	o.client = client

	// Check if repolist file is provided
	if o.RepoListFile == "" {
		return fmt.Errorf("repolist file is required")
	}

	repos, err := _gitlab.LoadRepositories(o.RepoListFile)
	if err != nil {
		return fmt.Errorf("failed to load repositories from repolist file: %v", err)
	}
	o.repositories = repos

	return nil
}

func (o *CountOptions) Validate() error {
	// Check if version is provided
	if o.BranchName == "" {
		return fmt.Errorf("branch name is required")
	}

	// Check if version is valid

	// check if repositories are valid
	if len(o.repositories) == 0 {
		return fmt.Errorf("no repositories found in repolist file")
	}

	return nil
}

func (o *CountOptions) Run() error {
	return _gitlab.IterateRepositories(o.repositories, countCallback, o)
}

func countCallback(repo _gitlab.Repository, options interface{}) error {
	o := options.(*CountOptions)

	compareAhead, _ := api.RepositoryCompare(o.client, repo.Pid, o.BranchName, _gitlab.DefaultRefName)
	ahead := len(compareAhead.Commits)
	compareBehind, _ := api.RepositoryCompare(o.client, repo.Pid, _gitlab.DefaultRefName, o.BranchName)
	behind := len(compareBehind.Commits)
	fmt.Printf("Repository(%-24s):\t%s<->%s\t%d|%d\n", repo.Pid, _gitlab.DefaultRefName, o.BranchName, ahead, behind)
	return nil
}
