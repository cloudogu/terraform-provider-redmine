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
	testVersionTFResourceType = "redmine_version"
	testVersionTFResourceName = "test_version1"
	testVersionTFResource     = testVersionTFResourceType + "." + testVersionTFResourceName
)

const (
	verKeyID          = "id"
	verKeyProjectID   = "project_id"
	verKeyName        = "name"
	verKeyDescription = "description"
	verKeyStatus      = "status"
	verKeyDueDate     = "due_date"
	verKeyCreatedOn   = "created_on"
	verKeyUpdatedOn   = "updated_on"
)

var projectResourceBlock = basicProjectWithDescription("testproject", "project", "a project")

const projectResourceIDReference = testProjectTFResource + ".id"

func TestAccVersionCreate_basic(t *testing.T) {
	tfProjectAndVersionBlocks := projectResourceBlock + "\n" +
		VersionAsHCL(testVersionTFResourceName, projectResourceIDReference, "Sprint 1", "desc", "open", "")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndVersionBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyID, "1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyName, "Sprint 1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDescription, "desc"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyStatus, "open"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDueDate, ""),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccVersionCreate_multipleIssueCategoriesToTheSameProject(t *testing.T) {
	tfProjectAndVersionBlocks := projectResourceBlock + "\n" +
		VersionAsHCL(testVersionTFResourceName, projectResourceIDReference, "Sprint 1", "desc", "closed", "2021-03-01") + "\n" +
		VersionAsHCL("another_version", projectResourceIDReference, "Sprint 2", "desc2", "locked", "") + "\n" +
		VersionAsHCL("yet_another_version", projectResourceIDReference, "Sprint 3", "desc3", "open", "")
	testVersionTFRessourceName2 := testVersionTFResourceType + "." + "another_version"
	testVersionTFRessourceName3 := testVersionTFResourceType + "." + "yet_another_version"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndVersionBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyName, "Sprint 1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDescription, "desc"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyStatus, "closed"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDueDate, "2021-03-01"),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyUpdatedOn),
					// check 2nd version
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testVersionTFRessourceName2, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName2, verKeyName, "Sprint 2"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName2, verKeyDescription, "desc2"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName2, verKeyStatus, "locked"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName2, verKeyDueDate, ""),
					resource.TestCheckResourceAttrSet(testVersionTFRessourceName2, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFRessourceName2, verKeyUpdatedOn),
					// check 3rd version
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testVersionTFRessourceName3, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName3, verKeyName, "Sprint 3"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName3, verKeyDescription, "desc3"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName3, verKeyStatus, "open"),
					resource.TestCheckResourceAttr(testVersionTFRessourceName3, verKeyDueDate, ""),
					resource.TestCheckResourceAttrSet(testVersionTFRessourceName3, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFRessourceName3, verKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccVersionUpdate_versionValuesChanged(t *testing.T) {
	tfProjectAndVersionBlocksCreation := projectResourceBlock + "\n" +
		VersionAsHCL(testVersionTFResourceName, projectResourceIDReference, "Sprint", "desc", "open", "")
	tfProjectAndVersionBlocksChanged := projectResourceBlock + "\n" +
		VersionAsHCL(testVersionTFResourceName, projectResourceIDReference, "Shazam! Renamed!", "Booyaka! Renamed!", "closed", "2021-03-01")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndVersionBlocksCreation,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyName, "Sprint"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDescription, "desc"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyStatus, "open"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDueDate, ""),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyUpdatedOn),
				),
			},
			{
				Config: tfProjectAndVersionBlocksChanged,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyName, "Shazam! Renamed!"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDescription, "Booyaka! Renamed!"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyStatus, "closed"),
					resource.TestCheckResourceAttr(testVersionTFResource, verKeyDueDate, "2021-03-01"),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testVersionTFResource, verKeyUpdatedOn),
				),
			},
		},
	})
}

func testAccCheckVersionDestroy(s *terraform.State) error {
	cli := testAccProvider.Meta().(*redmine.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testVersionTFResourceType {
			continue
		}

		VersionID := rs.Primary.ID

		// when
		issue, err := cli.ReadVersion(context.Background(), VersionID)

		// then
		if err == nil {
			if issue.ID != "" {
				return fmt.Errorf("version (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if !strings.Contains(err.Error(), "version (id: "+VersionID+") was not found") {
			return err
		}
	}

	return testAccCheckProjectDestroy(s)
}

func VersionAsHCL(tfName, projectID string, name, description, status, dueDate string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
  project_id = %s
  name = "%s"
	description = "%s"
	status = "%s"
	due_date = "%s"
}`, testVersionTFResourceType, tfName,
		projectID, name, description, status, dueDate)
}
