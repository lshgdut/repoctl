package api

import "github.com/xanzy/go-gitlab"

func TagList(client *gitlab.Client, pid interface{}, options *gitlab.ListTagsOptions) ([]*gitlab.Tag, error) {
	tags, _, err := client.Tags.ListTags(pid, options)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func TagCreate(client *gitlab.Client, pid interface{}, options *gitlab.CreateTagOptions) (*gitlab.Tag, error) {
	tag, _, err := client.Tags.CreateTag(pid, options)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func TagDelete(client *gitlab.Client, pid interface{}, tag string) error {
	_, err := client.Tags.DeleteTag(pid, tag)
	return err
}
