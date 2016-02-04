package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccGithubTeamMembership_basic(t *testing.T) {
	var membership github.Membership

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubTeamMembershipDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGithubTeamMembershipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubTeamMembershipExists("github_team_membership.test_team_membership", &membership),
					testAccCheckGithubTeamMembershipRoleState("github_team_membership.test_team_membership", &membership),
				),
			},
		},
	})
}

func testAccCheckGithubTeamMembershipDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*GithubClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_team_membership" {
			continue
		}

		t, u := parseTwoPartId(rs.Primary.ID)
		membership, resp, err := conn.Organizations.GetTeamMembership(toGithubId(t), u)
		if err == nil {
			if membership != nil {
				return fmt.Errorf("Team membership still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGithubTeamMembershipExists(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No team membership ID is set")
		}

		conn := testAccProvider.Meta().(*GithubClient).client
		t, u := parseTwoPartId(rs.Primary.ID)

		teamMembership, _, err := conn.Organizations.GetTeamMembership(toGithubId(t), u)

		if err != nil {
			return err
		}
		*membership = *teamMembership
		return nil
	}
}

func testAccCheckGithubTeamMembershipRoleState(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No team membership ID is set")
		}

		conn := testAccProvider.Meta().(*GithubClient).client
		t, u := parseTwoPartId(rs.Primary.ID)

		teamMembership, _, err := conn.Organizations.GetTeamMembership(toGithubId(t), u)
		if err != nil {
			return err
		}

		resourceRole := membership.Role
		actualRole := teamMembership.Role

		if *resourceRole != *actualRole {
			return fmt.Errorf("Team membership role %v in resource does match actual state of %v", *resourceRole, *actualRole)
		}
		return nil
	}
}

const testAccGithubTeamMembershipConfig = `
resource "github_membership" "test_org_membership" {
	username = "TerraformDummyUser"
	role = "member"
}

resource "github_team" "test_team" {
	name = "foo"
	description = "Terraform acc test group"
}

resource "github_team_membership" "test_team_membership" {
	team_id = "${github_team.test_team.id}"
	username = "TerraformDummyUser"
	role = "member"
}
`
