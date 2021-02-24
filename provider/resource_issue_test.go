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
	issKeyCreatedOn     = "created_on"
	issKeyUpdatedOn     = "updated_on"
)

func TestAccIssue_createBasic(t *testing.T) {
	projectResourceIDReference := testProjectTFResource + ".id"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicProjectWithDescription("testproject", "project", "a project") + "\n" +
					issueAsJSON(testIssueTFResourceName, projectResourceIDReference, 2, "issue subject", "This is an example issue"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyProjectID, "1"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyTrackerID, "2"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeySubject, "issue subject"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example issue"),
					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentIssueID, "0"),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
				),
			},
		},
	})
}

// func TestAccIssue_createMultipleIssues(t *testing.T) {
// 	const project2Name = "project2"
// 	const project2TFResource = testIssueTFResourceType + "." + project2Name
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckIssueDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: basicIssueWithDescription(issValueIdentifier, issValueName, "This is an example project") + "\n" +
// 						issueAsJSON(project2Name, "anotherident", "Another project", "Yet another project",
// 							"https://www.example.com/", true, false),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIdentifier, issValueIdentifier),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyName, issValueName),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyParentID, ""),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "This is an example project"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyHomepage, issValueHomepage),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIsPublic, "false"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyInheritMembers, "true"),
// 					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyCreatedOn),
// 					resource.TestCheckResourceAttrSet(testIssueTFResource, issKeyUpdatedOn),
// 					// check 2nd project
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyIdentifier, "anotherident"),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyName, "Another project"),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyParentID, ""),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyDescription, "Yet another project"),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyHomepage, "https://www.example.com/"),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyIsPublic, "true"),
// 					resource.TestCheckResourceAttr(project2TFResource, issKeyInheritMembers, "false"),
// 					resource.TestCheckResourceAttrSet(project2TFResource, issKeyCreatedOn),
// 					resource.TestCheckResourceAttrSet(project2TFResource, issKeyUpdatedOn),
// 				),
// 			},
// 		},
// 	})
// }
//
// func TestAccIssueUpdate(t *testing.T) {
// 	createdOn := "updated during 1. step"
// 	updatedOn := "updated during 1. and 2. step"
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckIssueDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: basicIssueWithDescription(issValueIdentifier, issValueName, "This is an example project"),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(testIssueTFResource, "id", "1"),
// 					func(state *terraform.State) error {
// 						rs := state.RootModule().Resources[testIssueTFResource]
// 						is := rs.Primary
// 						createdOn = is.Attributes[issKeyCreatedOn]
// 						updatedOn = is.Attributes[issKeyUpdatedOn]
//
// 						time.Sleep(2 + time.Second) // avoid flaky tests because Redmine's timestamp precision is seconds
// 						return nil
// 					},
// 				),
// 			},
// 			{
// 				Config: issueAsJSON(testIssueTFResourceName, issValueIdentifier, "nameChange", "descriptionChange",
// 					"homepageChange", true, false),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyID, "1"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIdentifier, issValueIdentifier),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyName, "nameChange"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "descriptionChange"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyHomepage, "homepageChange"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIsPublic, "true"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyInheritMembers, "false"),
// 					func(state *terraform.State) error {
// 						rs := state.RootModule().Resources[testIssueTFResource]
// 						is := rs.Primary
//
// 						actualCreatedOn := is.Attributes[issKeyCreatedOn]
// 						if err := assertNotEqual(issKeyCreatedOn, createdOn, actualCreatedOn); err != nil {
// 							return err
// 						}
// 						createdOn = is.Attributes[issKeyCreatedOn]
//
// 						actualUpdatedOn := is.Attributes[issKeyUpdatedOn]
// 						if err := assertEqual(issKeyUpdatedOn, updatedOn, actualUpdatedOn); err != nil {
// 							return err
// 						}
// 						updatedOn = is.Attributes[issKeyUpdatedOn]
//
// 						time.Sleep(2 + time.Second) // avoid flaky tests because Redmine's timestamp precision is seconds
// 						return nil
// 					},
// 				),
// 			},
// 			{
// 				Config: issueAsJSON(testIssueTFResourceName, issValueIdentifier, "nameChange2", "descriptionChange2",
// 					"homepageChange2", false, true),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyID, "1"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIdentifier, issValueIdentifier),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyName, "nameChange2"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyDescription, "descriptionChange2"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyHomepage, "homepageChange2"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyIsPublic, "false"),
// 					resource.TestCheckResourceAttr(testIssueTFResource, issKeyInheritMembers, "true"),
// 					func(state *terraform.State) error {
// 						rs := state.RootModule().Resources[testIssueTFResource]
// 						is := rs.Primary
//
// 						actualCreatedOn := is.Attributes[issKeyCreatedOn]
// 						if err := assertNotEqual(issKeyCreatedOn, createdOn, actualCreatedOn); err != nil {
// 							return err
// 						}
//
// 						actualUpdatedOn := is.Attributes[issKeyUpdatedOn]
// 						if err := assertEqual(issKeyUpdatedOn, updatedOn, actualUpdatedOn); err != nil {
// 							return err
// 						}
// 						return nil
// 					},
// 				),
// 			},
// 		},
// 	})
// }

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

func issueAsJSON(tfName, projectID string, trackerID int, subject, description string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
  project_id = %s
  tracker_id = %d
  subject = "%s"
  description = "%s"
}`, testIssueTFResourceType, tfName,
		projectID, trackerID, subject, description)
}
