package api

import "github.com/xanzy/go-gitlab"

func RepositoryCompare(client *gitlab.Client, pid interface{}, from, to string) (*gitlab.Compare, error) {
	comparison, _, err := client.Repositories.Compare(pid, &gitlab.CompareOptions{
		From: &from,
		To:   &to,
	})
	return comparison, err
}
