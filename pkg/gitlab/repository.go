package gitlab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Repo represents a repository with a PID and Ref.
type Repository struct {
	Pid string
	Ref string
}

type RepositoryCallback func(repo Repository, options interface{}) error

const (
	DefaultRefName = "main"
)

// readRepos reads the repository data from a file and returns a slice of Repo.
func LoadRepositories(filename string) ([]Repository, error) {
	var repos []Repository

	// Check if repolist file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			repos = append(repos, Repository{
				Pid: fields[0],
				Ref: fields[1],
			})
		} else {
			repos = append(repos, Repository{
				Pid: fields[0],
				Ref: DefaultRefName,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func ValidateRepositories(repos []Repository) error {
	for _, repo := range repos {
		fmt.Printf("Repo: %s, Ref: %s\n", repo.Pid, repo.Ref)
	}
	return nil
}

// iterateRepositories iterates over a slice of Repository and calls a callback function for each repository.
func IterateRepositories(repos []Repository, callback RepositoryCallback, options interface{}) error {
	for _, repo := range repos {
		if err := callback(repo, options); err != nil {
			return err
		}
	}
	return nil
}
