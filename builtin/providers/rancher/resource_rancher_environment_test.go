package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rancher/go-rancher/client"
)

func TestAccRancherEnvironment_Basic(t *testing.T) {
	var env client.Project
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckRancherEnvironmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foobar", &env),
					testAccCheckRancherEnvironmentAttributes(&env),
					resource.TestCheckResourceAttr(
						"rancher_environment.foobar", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"rancher_environment.foobar", "description", "Terraform Test Environment"),
					resource.TestCheckResourceAttr(
						"rancher_environment.foobar", "engine", "cattle"),
				),
			},
		},
	})
}

func testAccCheckRancherEnvironmentDestroy(s *terraform.State) error {
	rancher, _ := testAccProvider.Meta().(*ClientProvider).client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_environment" {
			continue
		}

		env, _ := rancher.Project.ById(rs.Primary.ID)
		if env != nil && env.State != "removed" {
			return fmt.Errorf("Environment[%s] still exists with state %s", env.Id, env.State)
		}

	}
	return nil
}

func testAccCheckRancherEnvironmentExists(n string, env *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rancher, _ := testAccProvider.Meta().(*ClientProvider).client()
		environment, err := rancher.Project.ById(rs.Primary.ID)
		if err != nil {
			return err
		}
		*env = *environment
		return nil
	}
}

func testAccCheckRancherEnvironmentAttributes(env *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

const testAccCheckRancherEnvironmentBasic = `
resource "rancher_environment" "foobar" {
  name = "foobar"
  description = "Terraform Test Environment"
}`
