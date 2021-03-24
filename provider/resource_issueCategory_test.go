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

const (
	testIssueCategoryTFResourceType = "redmine_issue_category"
	testIssueCategoryTFResourceName = "test_issue_category1"
	testIssueCategoryTFResource     = testIssueCategoryTFResourceType + "." + testIssueCategoryTFResourceName
)

const (
	issCatKeyID        = "id"
	issCatKeyProjectID = "project_id"
	issCatKeyName      = "name"
)

func TestAccIssueCategoryCreate_basic(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"
	tfProjectAndIssueCategoryBlocks := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueCategoryAsHCL(testIssueCategoryTFResourceName, projectResourceIDReference, "category name")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueCategoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueCategoryBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyName, "category name"),
				),
			},
		},
	})
}

func TestAccIssueCategoryCreate_multipleIssueCategoriesToTheSameProject(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"

	tfProjectAndIssueCategoryBlocks := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueCategoryAsHCL(testIssueCategoryTFResourceName, projectResourceIDReference, "category 1") + "\n" +
		issueCategoryAsHCL("another_issue_category", projectResourceIDReference, "category 2")
	testIssueCategoryTFRessourceName2 := testIssueCategoryTFResourceType + "." + "another_issue_category"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueCategoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueCategoryBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyName, "category 1"),
					// check 2nd issue category
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueCategoryTFRessourceName2, issCatKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFRessourceName2, issCatKeyName, "category 2"),
				),
			},
		},
	})
}

func TestAccIssueCategoryUpdate_categoryValuesChanged(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"
	tfProjectAndIssueCategoryBlocksCreation := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueCategoryAsHCL(testIssueCategoryTFResourceName, projectResourceIDReference, "original category")
	tfProjectAndIssueCategoryBlocksChanged := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueCategoryAsHCL(testIssueCategoryTFResourceName, projectResourceIDReference, "Booyaka! Renamed!")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueCategoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueCategoryBlocksCreation,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyName, "original category"),
				),
			},
			{
				Config: tfProjectAndIssueCategoryBlocksChanged,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyName, "Booyaka! Renamed!"),
				),
			},
		},
	})
}

func testAccCheckIssueCategoryDestroy(s *terraform.State) error {
	cli := testAccProvider.Meta().(*redmine.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testIssueCategoryTFResourceType {
			continue
		}

		issueCategoryID := rs.Primary.ID

		// when
		issue, err := cli.ReadIssueCategory(context.Background(), issueCategoryID)

		// then
		if err == nil {
			if issue.ID != "" {
				return fmt.Errorf("issue category (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if !strings.Contains(err.Error(), "issue category (id: "+issueCategoryID+") was not found") {
			return err
		}
	}

	return testAccCheckProjectDestroy(s)
}

func issueCategoryAsHCL(tfName, projectID string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
  project_id = %s
  name = "%s"
}`, testIssueCategoryTFResourceType, tfName,
		projectID, name)
}
