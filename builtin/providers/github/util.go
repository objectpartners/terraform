package github

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
)

func toGithubId(id string) int {
	githubId, _ := strconv.Atoi(id)
	return githubId
}

func fromGithubId(id *int) string {
	return strconv.Itoa(*id)
}

func validateRoleValueFunc(roles []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range roles {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid Github role type for %s", value, k))
		}
		return
	}
}

// return the pieces of id `a:b` as a, b
func parseTwoPartId(id string) (string, string) {
	parts := strings.SplitN(id, ":", 2)
	return parts[0], parts[1]
}

// format the strings into an id `a:b`
func buildTwoPartId(a, b *string) string {
	return fmt.Sprintf("%s:%s", *a, *b)
}
