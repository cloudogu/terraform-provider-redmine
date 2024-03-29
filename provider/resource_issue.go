package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

const (
	IssID            = "id"
	IssProjectID     = "project_id"
	IssTrackerID     = "tracker_id"
	IssSubject       = "subject"
	IssDescription   = "description"
	IssParentIssueID = "parent_issue_id"
	IssPriorityID    = "priority_id"
	IssCategoryID    = "category_id"
	IssCreatedOn     = "created_on"
	IssUpdatedOn     = "updated_on"
)

// IssueClient provides methods for reading and modifying Redmine issues.
type IssueClient interface {
	// CreateIssue creates an issue.
	CreateIssue(ctx context.Context, issue *redmine.Issue) (*redmine.Issue, error)
	// ReadIssue reads an issue identified by the id. The id must not be empty string or "0".
	ReadIssue(ctx context.Context, id string) (*redmine.Issue, error)
	// UpdateIssue updates an existing issue.
	UpdateIssue(ctx context.Context, issue *redmine.Issue) (*redmine.Issue, error)
	// DeleteIssue deletes an issue identified by the id. The id must not be empty string or "0".
	DeleteIssue(ctx context.Context, id string) error
}

func resourceIssue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIssueCreate,
		ReadContext:   resourceIssueRead,
		UpdateContext: resourceIssueUpdate,
		DeleteContext: resourceIssueDelete,
		Schema: map[string]*schema.Schema{
			IssID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			IssProjectID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			IssTrackerID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			IssSubject: {
				Type:     schema.TypeString,
				Required: true,
			},
			IssDescription: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			IssParentIssueID: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			IssPriorityID: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			IssCategoryID: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			IssCreatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			IssUpdatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceIssueRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	issueID := d.Get(IssID).(string)

	client := i.(IssueClient)
	issue, err := client.ReadIssue(ctx, issueID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("issue read id %s, project %d", issue.ID, issue.ProjectID)

	return issueSetToState(issue, d)
}

func resourceIssueCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueClient)

	issue := issueFromState(d)

	createdIssue, err := client.CreateIssue(ctx, issue)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdIssue.ID)

	log.Printf("issue create id %s, project %d", issue.ID, issue.ProjectID)

	diagRead := resourceIssueRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceIssueUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueClient)

	issue := issueFromState(d)

	_, err := client.UpdateIssue(ctx, issue)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("issue update id %s, project %d", issue.ID, issue.ProjectID)

	diagRead := resourceIssueRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceIssueDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(IssueClient)

	issueID := d.Id()
	err := client.DeleteIssue(ctx, issueID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("issue delete id %s", issueID)

	return diags
}

func issueSetToState(issue *redmine.Issue, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(issue.ID)
	if err := d.Set(IssProjectID, issue.ProjectID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssTrackerID, issue.TrackerID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssParentIssueID, issue.ParentIssueID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssSubject, issue.Subject); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssDescription, issue.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssPriorityID, issue.PriorityID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssCategoryID, issue.CategoryID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssCreatedOn, issue.CreatedOn); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(IssUpdatedOn, issue.UpdatedOn); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func issueFromState(d *schema.ResourceData) *redmine.Issue {
	issue := &redmine.Issue{}
	issue.ProjectID, _ = d.Get(IssProjectID).(int)
	issue.TrackerID = d.Get(IssTrackerID).(int)
	issue.Subject = d.Get(IssSubject).(string)
	issue.Description = d.Get(IssDescription).(string)
	issue.CreatedOn = d.Get(IssCreatedOn).(string)
	issue.UpdatedOn = d.Get(IssUpdatedOn).(string)

	issueID := d.Id()
	if issueID != "" && issueID != "0" {
		issue.ID = issueID
	}

	issuePriorityID := d.Get(IssPriorityID).(int)
	if issuePriorityID != 0 {
		issue.PriorityID = issuePriorityID
	}

	issueCategoryID := d.Get(IssCategoryID).(int)
	if issueCategoryID != 0 {
		issue.CategoryID = issueCategoryID
	}

	return issue
}
