package api

import "github.com/xanzy/go-gitlab"

func BranchList(client *gitlab.Client, pid interface{}, options *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error) {
	branches, _, err := client.Branches.ListBranches(pid, options)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func BranchCreate(client *gitlab.Client, pid interface{}, options *gitlab.CreateBranchOptions) (*gitlab.Branch, error) {
	branch, _, err := client.Branches.CreateBranch(pid, options)

	return branch, err
}

func BranchDelete(client *gitlab.Client, pid interface{}, branch string) error {
	_, err := client.Branches.DeleteBranch(pid, branch)

	return err
}

func BranchExists(client *gitlab.Client, pid interface{}, branch string) bool {
	b, _, err := client.Branches.GetBranch(pid, branch)
	if b != nil && err == nil {
		return true
	} else {
		return false
	}
}
