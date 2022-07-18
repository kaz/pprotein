package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type (
	RepositoryInfo struct {
		Ref     string
		Hash    string
		Author  string
		Message string
		Remote  string
	}
)

func GetInfo(repoDir string) (*RepositoryInfo, error) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repo: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get head: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("failed to get remote: %w", err)
	}

	remoteUrl := ""
	if len(remote.Config().URLs) > 0 {
		remoteUrl = remote.Config().URLs[0]
	}

	return &RepositoryInfo{
		Ref:     ref.Name().String(),
		Hash:    commit.Hash.String(),
		Author:  commit.Author.String(),
		Message: commit.Message,
		Remote:  remoteUrl,
	}, nil
}
