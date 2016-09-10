package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rancher/go-rancher/client"
)

func TestAccRancherStack_Basic(t *testing.T) {
	var stack client.Environment
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherStackDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckRancherStackBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherStackExists("rancher_stack.foo", &stack),
					testAccCheckRancherStackAttributes(&stack),
				),
			},
		},
	})
}

func testAccCheckRancherStackDestroy(s *terraform.State) error {
	rancher, _ := testAccProvider.Meta().(*ClientProvider).client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_stack" {
			continue
		}

		stack, _ := rancher.Environment.ById(rs.Primary.ID)
		if stack != nil && stack.State != "removed" {
			return fmt.Errorf("Stack[%s] still exists with state %s", stack.Id, stack.State)
		}
	}
	return nil
}

func testAccCheckRancherStackExists(n string, stack *client.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rancher, _ := testAccProvider.Meta().(*ClientProvider).client()
		t, err := rancher.Environment.ById(rs.Primary.ID)
		if err != nil {
			return err
		}
		*stack = *t
		return nil
	}
}

func testAccCheckRancherStackAttributes(token *client.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// if token.Command == "" || token.Image == "" || token.Token == "" {
		// 	return fmt.Errorf("RegistrationToken does not contain computed fields")
		// }
		return nil
	}
}

const testAccCheckRancherStackBasic = `
resource "rancher_environment" "foobar" {
  name = "foobar"
}

resource "rancher_stack" "foo" {
  environment_id = "${rancher_environment.foobar.id}"
  name = "foo"
  description = "Terraform acceptance test stack"
}`
