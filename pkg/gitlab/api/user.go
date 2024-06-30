package api

import "github.com/xanzy/go-gitlab"

func UserList(client *gitlab.Client) ([]*gitlab.User, error) {
	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	return users, nil
}
