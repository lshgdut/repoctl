package tag

import (
	"fmt"

	cmdutil "github.com/lshgdut/repoctl/pkg/cmd/util"
	_gitlab "github.com/lshgdut/repoctl/pkg/gitlab"
	"github.com/lshgdut/repoctl/pkg/gitlab/api"
	"k8s.io/klog/v2"

	// api "github.com/lshgdut/repoctl/pkg/gitlab/api"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type TagOptions struct {
	client *gitlab.Client
	dryRun bool

	args         []string
	repositories []_gitlab.Repository

	Version      string
	RepoListFile string

	tagName string
	tagFrom string
}

func NewTagOptions() *TagOptions {
	return &TagOptions{
		dryRun: false,
	}
}

func NewTagCmd() *cobra.Command {
	o := NewTagOptions()

	var cmd = &cobra.Command{
		Use:                   "tag --version <version> -f <repolist-file>",
		DisableFlagsInUseLine: true,
		Short:                 "tag release branch to <version>",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.Version, "version", "v", "", "release version, e.g. v1.0.0")
	cmd.Flags().StringVarP(&o.RepoListFile, "repolist-file", "f", "", "repolist file path, e.g. repolist.txt")

	cmdutil.AddDryRunFlag(cmd)
	return cmd
}

func (o *TagOptions) Complete(cmd *cobra.Command, args []string) error {
	// var err error
	var dryRunFlag = cmdutil.GetFlagBool(cmd, "dry-run")
	o.dryRun = dryRunFlag

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

	o.Version = cmdutil.TrimVersion(o.Version)
	o.tagFrom = cmdutil.GetReleaseBranchName(o.Version)
	o.tagName = cmdutil.GetReleaseTagName(o.Version)
	return nil
}

func (o *TagOptions) Validate() error {
	// Check if version is provided
	if o.Version == "" {
		return fmt.Errorf("version is required")
	}

	// Check if version is valid

	// check if repositories are valid
	if len(o.repositories) == 0 {
		return fmt.Errorf("no repositories found in repolist file")
	}

	return nil
}

func (o *TagOptions) Run() error {
	return _gitlab.IterateRepositories(o.repositories, tagCallback, o)
}

func tagCallback(repo _gitlab.Repository, options interface{}) error {
	o := options.(*TagOptions)

	klog.Infof("Taging release repository: %s", repo.Pid)
	if !o.dryRun {
		_, err := api.TagCreate(o.client, repo.Pid, &gitlab.CreateTagOptions{
			TagName: gitlab.Ptr(o.tagName),
			Ref:     gitlab.Ptr(o.tagFrom),
			Message: gitlab.Ptr(fmt.Sprintf("Release %s", o.tagName)),
		})
		if err != nil {
			return fmt.Errorf("failed to create tag: %v", err)
		}
		klog.Infof("Tag created(%-26s): %s->%s", repo.Pid, o.tagFrom, o.tagName)
	}

	return nil
}
