package provider

import (
	"context"
	"fmt"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
	"time"
)

const (
	testIssueTFResourceType = "redmine_issue"
	testIssueTFResourceName = "testissue1"
	testIssueTFResource     = testIssueTFResourceType + "." + testIssueTFResourceName
)

const (
	issKeyID            = "id"
	issKeyProjectID     = "project_id"
	issKeyTrackerID     = "tracker_id"
	issKeySubject       = "subject"
	issKeyDescription   = "description"
	issKeyParentIssueID = "parent_issue_id"
	issKeyPriorityID    = "priority_id"
	issKeyCategoryID    = "category_id"
	issKeyCreatedOn     = "created_on"
	issKeyUpdatedOn     = "updated_on"
)

func TestAccIssueCreate_basic(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"
	tfProjectAndIssueBlocks := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue", 2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccIssueCreate_multipleIssuesToTheSameProject(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"

	tfProjectAndIssueBlocks := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue", 2) + "\n" +
		issueAsHCL("another_issue", projectResourceIDReference, 1, "issue subject2", "This is an example issue2", 5)
	testIssueTFRessourceName2 := testIssueTFResourceType + "." + "another_issue"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueBlocks,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
					// check 2nd issue
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyTrackerID, "1"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeySubject, "issue subject2"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyDescription, "This is an example issue2"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyPriorityID, "5"),
					resource.TestCheckResourceAttr(testIssueTFRessourceName2, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFRessourceName2, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFRessourceName2, issKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccIssueUpdate_issueValuesChanged(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"
	tfProjectAndIssueBlocksCreation := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue", 2)
	tfProjectAndIssueBlocksChanged := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceIDReference, 1, "subjectChanged", "descriptionChanged", 5)

	createdOn := "updated during 1. step"
	updatedOn := "updated during 1. and 2. step"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueBlocksCreation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueTFResource, "id", "1"),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testIssueTFResource]
						is := rs.Primary
						createdOn = is.Attributes[issKeyCreatedOn]
						updatedOn = is.Attributes[issKeyUpdatedOn]

						time.Sleep(2 + time.Second) // avoid flaky tests because Redmine's timestamp precision is seconds
						return nil
					},
				),
			},
			{
				Config: tfProjectAndIssueBlocksChanged,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "subjectChanged"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "descriptionChanged"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "5"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testIssueTFResource]
						is := rs.Primary

						actualCreatedOn := is.Attributes[issKeyCreatedOn]
						if err := assertEqual(issKeyCreatedOn, createdOn, actualCreatedOn); err != nil {
							return err
						}

						actualUpdatedOn := is.Attributes[issKeyUpdatedOn]
						if err := assertNotEqual(issKeyUpdatedOn, updatedOn, actualUpdatedOn); err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccIssueUpdate_movesIssueToAnotherProject(t *testing.T) {
	projectResourceID1Reference := testProjectTFResource + ".id"
	projectResourceID2Reference := testProjectTFResourceType + ".project2.id"

	tfProjectAndIssueBlocksCreated := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		genericProjectAsHCL("project2", "anotherproject", "target project for moved issues", "desc", "", false, false) + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceID1Reference, 2, "issue subject", "This is an example issue", 2)

	tfProjectAndIssueBlocksMovedIssue := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		genericProjectAsHCL("project2", "anotherproject", "target project for moved issues", "desc", "", false, false) + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceID2Reference, 2, "issue subject", "This is an example issue", 2)

	projectIDFirstRun := "updated in 1. step"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectAndIssueBlocksCreated,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testIssueTFResource]
						is := rs.Primary

						projectIDFirstRun = is.Attributes[issKeyProjectID]
						return nil
					},
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			}, {
				Config: tfProjectAndIssueBlocksMovedIssue,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testIssueTFResource]
						is := rs.Primary

						// do not test id equality here because creation sequence is not guaranteed
						actualProjectID := is.Attributes[issKeyProjectID]
						if err := assertNotEqual(issKeyProjectID, projectIDFirstRun, actualProjectID); err != nil {
							return err
						}
						return nil
					},
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccIssueCreate_addAndRemoveIssueCategory(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"

	tfProjectIssueBlocksWithoutIssueCategory := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueAsHCL(testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue", 2)

	issueWithCategory := fmt.Sprintf(`resource "%s" "%s" {
  project_id = %s
  tracker_id = %d
  subject = "%s"
  description = "%s"
  priority_id = "%d"
  category_id = %s.id
}`, testIssueTFResourceType, testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue", 2, testIssueCategoryTFResource)

	tfProjectIssueBlocksWithIssueCategory := basicProjectWithDescription("testproject", "project", "a project") + "\n" +
		issueWithCategory + "\n" +
		issueCategoryAsHCL(testIssueCategoryTFResourceName, projectResourceIDReference, "Bananas and other tropical fruits")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfProjectIssueBlocksWithoutIssueCategory,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
			{
				// add issue category
				Config: tfProjectIssueBlocksWithIssueCategory,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueCategoryTFResource, issCatKeyName, "Bananas and other tropical fruits"),
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "1"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
			{
				// remove issue category again
				Config: tfProjectIssueBlocksWithoutIssueCategory,
				Check: resource.ComposeAggregateTestCheckFunc(
					// do not test id's here because creation sequence is not guaranteed
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyCategoryID, "0"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyPriorityID, "2"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
		},
	})
}

func testAccCheckIssueDestroy(s *terraform.State) error {
	cli := testAccProvider.Meta().(*redmine.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testIssueTFResourceType {
			continue
		}

		issueID := rs.Primary.ID

		// when
		issue, err := cli.ReadIssue(context.Background(), issueID)

		// then
		if err == nil {
			if issue.ID != "" {
				return fmt.Errorf("issue (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if !strings.Contains(err.Error(), "issue (id: "+issueID+") was not found") {
			return err
		}
	}

	return testAccCheckProjectDestroy(s)
}

func issueAsHCL(tfName, projectID string, trackerID int, subject, description string, prioID int) string {
	return fmt.Sprintf(`resource "%s" "%s" {
  project_id = %s
  tracker_id = %d
  subject = "%s"
  description = "%s"
  priority_id = "%d"
}`, testIssueTFResourceType, tfName,
		projectID, trackerID, subject, description, prioID)
}
