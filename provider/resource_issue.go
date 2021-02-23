package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	IssID            = "id"
	IssProjectID     = "project_id"
	IssTrackerID     = "tracker_id"
	IssSubject       = "subject"
	IssDescription   = "description"
	IssParentIssueID = "parent_issue_id"
	IssCreatedOn     = "created_on"
	IssUpdatedOn     = "updated_on"
)

type IssueClient interface {
	CreateIssue(ctx context.Context, issue *redmine.Issue) (*redmine.Issue, error)
	ReadIssue(ctx context.Context, id string) (*redmine.Issue, error)
	UpdateIssue(ctx context.Context, issue *redmine.Issue) (*redmine.Issue, error)
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

	diagRead := resourceIssueRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceIssueDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
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

	return issue
}
