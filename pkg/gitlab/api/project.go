package api

import (
	"github.com/xanzy/go-gitlab"
)

func ProjectList(client *gitlab.Client, options *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := client.Projects.ListProjects(options)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func ProjectGet(client *gitlab.Client, id interface{}) (*gitlab.Project, error) {
	project, _, err := client.Projects.GetProject(id, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	return project, nil
}
