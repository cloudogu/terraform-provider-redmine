package provider

import (
	"context"
	"fmt"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

const (
	testProjectTFResourceType = "redmine_project"
	testProjectTFResourceName = "testproject"
	testProjectTFResource     = testProjectTFResourceType + "." + testProjectTFResourceName
)

const (
	prjValueIdentifier = "exampleproject"
	prjValueName       = "Example Project"
	prjValueHomepage   = "https://cloudogu.com/"

	prjKeyID             = "id"
	prjKeyIdentifier     = "identifier"
	prjKeyName           = "name"
	prjKeyParentID       = "parent_id"
	prjKeyDescription    = "description"
	prjKeyHomepage       = "homepage"
	prjKeyIsPublic       = "is_public"
	prjKeyInheritMembers = "inherit_members"
	prjKeyCreatedOn      = "created_on"
	prjKeyUpdatedOn      = "updated_on"
)

func TestAccProject_createBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicProjectWithDescription(prjValueIdentifier, prjValueName, "This is an example project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyID, "1"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIdentifier, prjValueIdentifier),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyName, prjValueName),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyParentID, ""),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyDescription, "This is an example project"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyHomepage, prjValueHomepage),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIsPublic, "false"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyInheritMembers, "true"),
					resource.TestCheckResourceAttrSet(testProjectTFResource, prjKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testProjectTFResource, prjKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccProject_createMultipleProjects(t *testing.T) {
	const project2Name = "project2"
	const project2TFResource = testProjectTFResourceType + "." + project2Name
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicProjectWithDescription(prjValueIdentifier, prjValueName, "This is an example project") + "\n" +
					genericProjectAsJSON(project2Name, "anotherident", "Another project", "Yet another project",
						"https://www.example.com/", true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIdentifier, prjValueIdentifier),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyName, prjValueName),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyParentID, ""),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyDescription, "This is an example project"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyHomepage, prjValueHomepage),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIsPublic, "false"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyInheritMembers, "true"),
					resource.TestCheckResourceAttrSet(testProjectTFResource, prjKeyCreatedOn),
					resource.TestCheckResourceAttrSet(testProjectTFResource, prjKeyUpdatedOn),
					// check 2nd project
					resource.TestCheckResourceAttr(project2TFResource, prjKeyIdentifier, "anotherident"),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyName, "Another project"),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyParentID, ""),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyDescription, "Yet another project"),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyHomepage, "https://www.example.com/"),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyIsPublic, "true"),
					resource.TestCheckResourceAttr(project2TFResource, prjKeyInheritMembers, "false"),
					resource.TestCheckResourceAttrSet(project2TFResource, prjKeyCreatedOn),
					resource.TestCheckResourceAttrSet(project2TFResource, prjKeyUpdatedOn),
				),
			},
		},
	})
}

func TestAccProjectUpdate(t *testing.T) {
	createdOn := "updated during 1. step"
	updatedOn := "updated during 1. and 2. step"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicProjectWithDescription(prjValueIdentifier, prjValueName, "This is an example project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testProjectTFResource, "id", "1"),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testProjectTFResource]
						is := rs.Primary
						createdOn = is.Attributes[prjKeyCreatedOn]
						updatedOn = is.Attributes[prjKeyUpdatedOn]

						time.Sleep(2 + time.Second) // avoid flaky tests because Redmine's timestamp precision is seconds
						return nil
					},
				),
			},
			{
				Config: genericProjectAsJSON(testProjectTFResourceName, prjValueIdentifier, "nameChange", "descriptionChange",
					"homepageChange", true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyID, "1"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIdentifier, prjValueIdentifier),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyName, "nameChange"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyDescription, "descriptionChange"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyHomepage, "homepageChange"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIsPublic, "true"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyInheritMembers, "false"),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testProjectTFResource]
						is := rs.Primary

						actualCreatedOn := is.Attributes[prjKeyCreatedOn]
						if err := assertNotEqual(prjKeyCreatedOn, createdOn, actualCreatedOn); err != nil {
							return err
						}
						createdOn = is.Attributes[prjKeyCreatedOn]

						actualUpdatedOn := is.Attributes[prjKeyUpdatedOn]
						if err := assertEqual(prjKeyUpdatedOn, updatedOn, actualUpdatedOn); err != nil {
							return err
						}
						updatedOn = is.Attributes[prjKeyUpdatedOn]

						time.Sleep(2 + time.Second) // avoid flaky tests because Redmine's timestamp precision is seconds
						return nil
					},
				),
			},
			{
				Config: genericProjectAsJSON(testProjectTFResourceName, prjValueIdentifier, "nameChange2", "descriptionChange2",
					"homepageChange2", false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyID, "1"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIdentifier, prjValueIdentifier),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyName, "nameChange2"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyDescription, "descriptionChange2"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyHomepage, "homepageChange2"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyIsPublic, "false"),
					resource.TestCheckResourceAttr(testProjectTFResource, prjKeyInheritMembers, "true"),
					func(state *terraform.State) error {
						rs := state.RootModule().Resources[testProjectTFResource]
						is := rs.Primary

						actualCreatedOn := is.Attributes[prjKeyCreatedOn]
						if err := assertNotEqual(prjKeyCreatedOn, createdOn, actualCreatedOn); err != nil {
							return err
						}

						actualUpdatedOn := is.Attributes[prjKeyUpdatedOn]
						if err := assertEqual(prjKeyUpdatedOn, updatedOn, actualUpdatedOn); err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	cli := testAccProvider.Meta().(*redmine.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testProjectTFResourceType {
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

		if !strings.Contains(err.Error(), "project (id: "+projectID+") was not found") {
			return err
		}
	}

	return nil
}

func basicProjectWithDescription(identifier, name, description string) string {
	const isPublic = false
	const inheritMembers = true

	return genericProjectAsJSON(testProjectTFResourceName, identifier, name, description, prjValueHomepage, isPublic, inheritMembers)
}

func genericProjectAsJSON(tfName, identifier, name, description, homepage string, isPublic, inheritMembers bool) string {
	return fmt.Sprintf(`resource "%s" "%s" {
  identifier = "%s"
  name = "%s"
  description = "%s"
  homepage = "%s"
  is_public = %t
  inherit_members = %t
}`, testProjectTFResourceType, tfName,
		identifier, name, description, homepage, isPublic, inheritMembers)
}

func assertEqual(resourceField string, expected, actual interface{}) error {
	if assert.ObjectsAreEqual(expected, actual) {
		return fmt.Errorf("Not equal: %s\n"+
			"expected: %s\n"+
			"actual  : %s", resourceField, expected, actual)
	}
	return nil
}

func assertNotEqual(resourceField string, expected, actual interface{}) error {
	if !assert.ObjectsAreEqual(expected, actual) {
		return fmt.Errorf("equal: %s\n"+
			"expected: %s\n"+
			"actual  : %s", resourceField, expected, actual)
	}
	return nil
}
