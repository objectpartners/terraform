package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccGithubTeam_basic(t *testing.T) {
	var team github.Team

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubTeamDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGithubTeamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubTeamExists("github_team.foo", &team),
				),
			},
		},
	})
}

func testAccCheckGithubTeamExists(n string, team *github.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Team ID is set")
		}

		conn := testAccProvider.Meta().(*GithubClient).client
		githubTeam, _, err := conn.Organizations.GetTeam(toGithubId(rs.Primary.ID))
		if err != nil {
			return err
		}
		*team = *githubTeam
		return nil
	}
}

func testAccCheckGithubTeamDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*GithubClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_team" {
			continue
		}

		team, resp, err := conn.Organizations.GetTeam(toGithubId(rs.Primary.ID))
		if err == nil {
			if team != nil &&
				fromGithubId(team.ID) == rs.Primary.ID {
				return fmt.Errorf("Team still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

const testAccGithubTeamConfig = `
resource "github_team" "foo" {
	name = "foo"
	description = "Terraform acc test group"
}
`
