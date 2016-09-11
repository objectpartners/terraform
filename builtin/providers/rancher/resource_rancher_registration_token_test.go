package rancher

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rancher/go-rancher/client"
)

func TestAccRancherRegistrationToken_Basic(t *testing.T) {
	var token client.RegistrationToken
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherRegistrationTokenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckRancherRegistrationTokenBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherRegistrationTokenExists("rancher_registration_token.foo", &token),
					testAccCheckRancherRegistrationTokenAttributes(&token),
					resource.TestCheckResourceAttr(
						"rancher_registration_token.foo", "description", "created"),
				),
			},
			resource.TestStep{
				Config: testAccCheckRancherRegistrationTokenUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherRegistrationTokenExists("rancher_registration_token.foo", &token),
					testAccCheckRancherRegistrationTokenAttributes(&token),
					resource.TestCheckResourceAttr(
						"rancher_registration_token.foo", "description", "updated"),
				),
			},
		},
	})
}

func testAccCheckRancherRegistrationTokenDestroy(s *terraform.State) error {
	rancher, _ := testAccProvider.Meta().(*ClientProvider).client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_registration_token" {
			continue
		}

		token, _ := rancher.RegistrationToken.ById(rs.Primary.ID)
		if token != nil && token.State != "removed" {
			return fmt.Errorf("RegistrationToken[%s] still exists with state %s", token.Id, token.State)
		}
	}
	return nil
}

func testAccCheckRancherRegistrationTokenExists(n string, token *client.RegistrationToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		rancher, _ := testAccProvider.Meta().(*ClientProvider).client()
		t, err := rancher.RegistrationToken.ById(rs.Primary.ID)
		if err != nil {
			return err
		}
		*token = *t
		return nil
	}
}

func testAccCheckRancherRegistrationTokenAttributes(token *client.RegistrationToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if token.Command == "" || token.Image == "" || token.Token == "" {
			return errors.New("RegistrationToken does not contain computed fields")
		}
		return nil
	}
}

const testAccCheckRancherRegistrationTokenBasic = `
resource "rancher_environment" "foobar" {
  name = "foobar"
}

resource "rancher_registration_token" "foo" {
  environment_id = "${rancher_environment.foobar.id}"
	description = "created"
}`

const testAccCheckRancherRegistrationTokenUpdate = `
resource "rancher_environment" "foobar" {
  name = "foobar"
}

resource "rancher_registration_token" "foo" {
  environment_id = "${rancher_environment.foobar.id}"
	description = "updated"
}`
