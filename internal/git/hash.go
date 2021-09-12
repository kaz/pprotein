package git

import "github.com/go-git/go-git/v5"

func GetCommitHash(repoDir string) string {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return ""
	}
	ref, err := repo.Head()
	if err != nil {
		return ""
	}
	return ref.Hash().String()
}
