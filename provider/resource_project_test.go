package provider

import (
	"context"
	"fmt"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicProjectWithDescription("exampleproject", "Example Project", "This is an example project"),
				Check:  resource.TestCheckResourceAttr("redmine_project.testproject", "id", "1"),
			},
		},
	})
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	cli := testAccProvider.Meta().(*redmine.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redmine_project" {
			continue
		}

		projectID := rs.Primary.ID

		// when
		project, err := cli.ReadProject(context.Background(), projectID)

		// then
		if err == nil {
			if project.ID != "" {
				return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		// If the error is equivalent to 404 not found, the widget is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "asdf") {
			return err
		}
	}

	return nil
}

func basicProjectWithDescription(identifier, name, description string) string {
	return fmt.Sprintf(`resource "redmine_project" "testproject1" {
  identifier = "%s"
  name = "%s"
  description = "%s"
  homepage = "https://cloudogu.com/"
  is_public = false
  inherit_members = true
}`, identifier, name, description)
}
