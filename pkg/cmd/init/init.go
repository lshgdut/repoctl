package init

import (
	cmdutil "github.com/lshgdut/repoctl/pkg/cmd/util"
	"github.com/lshgdut/repoctl/pkg/config"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	config.Config
	dryRun bool

	args []string
}

func NewInitOptions() *InitOptions {
	return &InitOptions{
		Config: config.Config{
			GitlabToken: "",
			GitlabUrl:   "",
		},
		dryRun: false,
	}
}

func NewInitCmd() *cobra.Command {
	o := NewInitOptions()

	var cmd = &cobra.Command{
		Use:                   "init --token=<token> --url=<url>",
		DisableFlagsInUseLine: true,
		Short:                 "initialize gitlab token",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	o.AddInitFlags(cmd)
	// cmd.Flags().StringVar(&o.Token, "token", o.Token, "gitlab token")
	// cmd.Flags().StringVar(&o.GitlabUrl, "url", o.GitlabUrl, "gitlab url")

	cmdutil.AddDryRunFlag(cmd)
	return cmd
}

func (o *InitOptions) AddInitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.GitlabToken, "token", o.GitlabToken, "gitlab token")
	cmd.Flags().StringVar(&o.GitlabUrl, "url", o.GitlabUrl, "gitlab url")
}

func (o *InitOptions) Complete(cmd *cobra.Command, args []string) error {
	// var err error
	var dryRunFlag = cmdutil.GetFlagBool(cmd, "dry-run")
	o.dryRun = dryRunFlag

	o.args = args
	return nil
}

func (o *InitOptions) Validate() error {
	return nil
}

func (o *InitOptions) Run() error {
	gitlabConfig := config.Config{
		GitlabToken: o.GitlabToken,
		GitlabUrl:   o.GitlabUrl,
	}

	if err := gitlabConfig.Save(o.dryRun); err != nil {
		return err
	}

	return nil
}
