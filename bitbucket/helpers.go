package bitbucket

import (
	"fmt"
)

func baseConfigForRepositoryBasedTests(projectKey string) string {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "repo"
		}
	`, projectKey, projectKey)

	return config
}
