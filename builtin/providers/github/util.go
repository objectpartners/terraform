package github

import (
	"strconv"
)

func toGithubId(id string) int {
	githubId, _ := strconv.Atoi(id)
	return githubId
}

func fromGithubId(id *int) string {
	return strconv.Itoa(*id)
}
