package release

import (
	"fmt"

	cmdutil "github.com/lshgdut/repoctl/pkg/cmd/util"
	"github.com/lshgdut/repoctl/pkg/gitlab/api"
	"github.com/xanzy/go-gitlab"
)

func unreleaseRepoBranch(client *gitlab.Client, pid string, version string) error {
	release := cmdutil.GetReleaseBranchName(version)
	unprotectedRepositoryBranch(client, pid, release)
	deleteRepositoryBranch(client, pid, release)
	return nil
}

func unprotectedRepositoryBranch(client *gitlab.Client, pid, branch string) error {
	client.ProtectedBranches.UnprotectRepositoryBranches(pid, branch)
	return nil
}

func deleteRepositoryBranch(client *gitlab.Client, pid, release string) error {
	_, err := client.Branches.DeleteBranch(pid, release)
	return err
}

func releaseRepoBranch(client *gitlab.Client, pid string, version, ref string) error {
	release := cmdutil.GetReleaseBranchName(version)
	pbs, _, _ := client.ProtectedBranches.ListProtectedBranches(pid, &gitlab.ListProtectedBranchesOptions{})
	for _, pb := range pbs {
		if pb.Name == release {
			fmt.Println("protected release branch already exists", pid, release)
			return nil
		}
	}

	if err := createReleaseBranchIfNotExist(client, pid, release, ref); err != nil {
		fmt.Println("Failed to create release branch", pid, release, err)
		return err
	} else {
		fmt.Println("Created release branch", pid, release, ref)
	}

	if err := protectedRepositoryBranch(client, pid, release); err != nil {
		fmt.Println("Failed to protect release branch", pid, release, err)
		return err
	} else {
		fmt.Println("Protected release branch", pid, release, ref)
	}
	return nil
}

func createReleaseBranchIfNotExist(client *gitlab.Client, pid, release, ref string) error {

	// Check if the release branch already exists
	if api.BranchExists(client, pid, release) {
		fmt.Println("release branch already exists", pid, release)
		return nil
	}

	// Create the release branch
	_, err_create := api.BranchCreate(client, pid, &gitlab.CreateBranchOptions{
		Branch: gitlab.Ptr(release),
		Ref:    gitlab.Ptr(ref),
	})
	return err_create
}

func protectedRepositoryBranch(client *gitlab.Client, pid, branch string) error {

	// Unprotect if the protected branch already exists
	unprotectedRepositoryBranch(client, pid, branch)

	// Protect the release branch
	_, _, err := client.ProtectedBranches.ProtectRepositoryBranches(pid, &gitlab.ProtectRepositoryBranchesOptions{
		Name:             &branch,
		PushAccessLevel:  gitlab.Ptr(gitlab.NoPermissions),
		MergeAccessLevel: gitlab.Ptr(gitlab.NoPermissions),
	})

	if err != nil {
		fmt.Println("Failed to protect release branch", pid, branch, err)
		return err
	}

	return nil
}
