package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

const (
	IssCatID        = "id"
	IssCatProjectID = "project_id"
	IssCatName      = "name"
)

// IssueCategoryClient provides methods for reading and modifying Redmine IssueCategories.
type IssueCategoryClient interface {
	// CreateIssueCategory creates an IssueCategory.
	CreateIssueCategory(ctx context.Context, IssueCategory *redmine.IssueCategory) (*redmine.IssueCategory, error)
	// ReadIssueCategory reads an IssueCategory identified by the id. The id must not be empty string or "0".
	ReadIssueCategory(ctx context.Context, id string) (*redmine.IssueCategory, error)
	// UpdateIssueCategory updates an existing IssueCategory.
	UpdateIssueCategory(ctx context.Context, IssueCategory *redmine.IssueCategory) (*redmine.IssueCategory, error)
	// DeleteIssueCategory deletes an IssueCategory identified by the id. The id must not be empty string or "0".
	DeleteIssueCategory(ctx context.Context, id string) error
}

func resourceIssueCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIssueCategoryCreate,
		ReadContext:   resourceIssueCategoryRead,
		UpdateContext: resourceIssueCategoryUpdate,
		DeleteContext: resourceIssueCategoryDelete,
		Schema: map[string]*schema.Schema{
			IssCatID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			IssCatProjectID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			IssCatName: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIssueCategoryRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	IssueCategoryID := d.Get(IssID).(string)

	client := i.(IssueCategoryClient)
	IssueCategory, err := client.ReadIssueCategory(ctx, IssueCategoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("IssueCategory read id %s, project %d", IssueCategory.ID, IssueCategory.ProjectID)

	return IssueCategorySetToState(IssueCategory, d)
}

func resourceIssueCategoryCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueCategoryClient)

	IssueCategory := IssueCategoryFromState(d)

	createdIssueCategory, err := client.CreateIssueCategory(ctx, IssueCategory)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdIssueCategory.ID)

	log.Printf("IssueCategory create id %s, project %d", IssueCategory.ID, IssueCategory.ProjectID)

	diagRead := resourceIssueCategoryRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceIssueCategoryUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueCategoryClient)

	IssueCategory := IssueCategoryFromState(d)

	_, err := client.UpdateIssueCategory(ctx, IssueCategory)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("IssueCategory update id %s, project %d", IssueCategory.ID, IssueCategory.ProjectID)

	diagRead := resourceIssueCategoryRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceIssueCategoryDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueCategoryClient)

	IssueCategoryID := d.Id()
	err := client.DeleteIssueCategory(ctx, IssueCategoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("IssueCategory delete id %s", IssueCategoryID)

	return diags
}

func IssueCategorySetToState(IssueCategory *redmine.IssueCategory, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(IssueCategory.ID)
	if err := d.Set(IssCatID, IssueCategory.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssCatProjectID, IssueCategory.ProjectID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssCatName, IssueCategory.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func IssueCategoryFromState(d *schema.ResourceData) *redmine.IssueCategory {
	IssueCategory := &redmine.IssueCategory{}
	IssueCategory.ProjectID, _ = d.Get(IssCatProjectID).(int)
	IssueCategory.Name = d.Get(IssCatName).(string)

	IssueCategoryID := d.Id()
	if IssueCategoryID != "" && IssueCategoryID != "0" {
		IssueCategory.ID = IssueCategoryID
	}

	return IssueCategory
}
