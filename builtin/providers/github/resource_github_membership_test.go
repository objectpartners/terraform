package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccGithubMembership_basic(t *testing.T) {
	var membership github.Membership

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubMembershipDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGithubMembershipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubMembershipExists("github_membership.test_org_membership", &membership),
					testAccCheckGithubMembershipRoleState("github_membership.test_org_membership", &membership),
				),
			},
		},
	})
}

func testAccCheckGithubMembershipDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*GithubClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_membership" {
			continue
		}
		o, u := parseMembershipId(rs.Primary.ID)

		membership, resp, err := conn.Organizations.GetOrgMembership(u, o)

		if err == nil {
			if membership != nil &&
				buildMembershipId(membership.Organization.Login, membership.User.Login) == rs.Primary.ID {
				return fmt.Errorf("Organization membership still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGithubMembershipExists(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		conn := testAccProvider.Meta().(*GithubClient).client
		o, u := parseMembershipId(rs.Primary.ID)

		githubMembership, _, err := conn.Organizations.GetOrgMembership(u, o)
		if err != nil {
			return err
		}
		*membership = *githubMembership
		return nil
	}
}

func testAccCheckGithubMembershipRoleState(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		conn := testAccProvider.Meta().(*GithubClient).client
		o, u := parseMembershipId(rs.Primary.ID)

		githubMembership, _, err := conn.Organizations.GetOrgMembership(u, o)
		if err != nil {
			return err
		}

		resourceRole := membership.Role
		actualRole := githubMembership.Role

		if *resourceRole != *actualRole {
			return fmt.Errorf("Membership role %v in resource does match actual state of %v", *resourceRole, *actualRole)
		}
		return nil
	}
}

func TestAccResourceGithubMembership_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "invalid",
			ErrCount: 1,
		},
		{
			Value:    "member",
			ErrCount: 0,
		},
		{
			Value:    "admin",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateRoleValue(tc.Value, "github_membership")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected github_membership to trigger a validation error")
		}
	}
}

const testAccGithubMembershipConfig = `
resource "github_membership" "test_org_membership" {
	username = "TerraformDummyUser"
	role = "member"
}
`
