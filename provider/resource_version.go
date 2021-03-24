package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
)

const (
	VerID          = "id"
	VerProjectID   = "project_id"
	VerName        = "name"
	VerDescription = "description"
	VerStatus      = "status"
	VerDueDate     = "due_date"
	VerCreatedOn   = "created_on"
	VerUpdatedOn   = "updated_on"
)

var dueDateYYYYMMDDRegexp, _ = regexp.Compile(`^(\d{4}-\d{2}-\d{2})?$`)

// VersionClient provides methods for reading and modifying Redmine versions.
type VersionClient interface {
	// CreateVersion creates an Version.
	CreateVersion(ctx context.Context, Version *redmine.Version) (*redmine.Version, error)
	// ReadVersion reads an Version identified by the id. The id must not be empty string or "0".
	ReadVersion(ctx context.Context, id string) (*redmine.Version, error)
	// UpdateVersion updates an existing Version.
	UpdateVersion(ctx context.Context, Version *redmine.Version) (*redmine.Version, error)
	// DeleteVersion deletes an Version identified by the id. The id must not be empty string or "0".
	DeleteVersion(ctx context.Context, id string) error
}

func resourceVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVersionCreate,
		ReadContext:   resourceVersionRead,
		UpdateContext: resourceVersionUpdate,
		DeleteContext: resourceVersionDelete,
		Schema: map[string]*schema.Schema{
			VerID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			VerProjectID: {
				Type:     schema.TypeInt,
				Required: true,
			},
			VerName: {
				Type:     schema.TypeString,
				Required: true,
			},
			VerDescription: {
				Type:     schema.TypeString,
				Required: true,
			},
			VerStatus: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"open", "locked", "closed"}, false),
				Default:      "open",
			},
			VerDueDate: {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringMatch(dueDateYYYYMMDDRegexp, "invalid due date found; expected either empty string or formatted date (YYYY-MM-DD)"),
				Optional:     true,
			},
			VerCreatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			VerUpdatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVersionRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	VersionID := d.Get(IssID).(string)

	client := i.(VersionClient)
	Version, err := client.ReadVersion(ctx, VersionID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Version read id %s, project %d", Version.ID, Version.ProjectID)

	return VersionSetToState(Version, d)
}

func resourceVersionCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(VersionClient)

	Version := VersionFromState(d)

	createdVersion, err := client.CreateVersion(ctx, Version)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdVersion.ID)

	log.Printf("Version create id %s, project %d", Version.ID, Version.ProjectID)

	diagRead := resourceVersionRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceVersionUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(VersionClient)

	Version := VersionFromState(d)

	_, err := client.UpdateVersion(ctx, Version)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Version update id %s, project %d", Version.ID, Version.ProjectID)

	diagRead := resourceVersionRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceVersionDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(VersionClient)

	VersionID := d.Id()
	err := client.DeleteVersion(ctx, VersionID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Version delete id %s", VersionID)

	return diags
}

func VersionSetToState(Version *redmine.Version, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(Version.ID)
	if err := d.Set(VerID, Version.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerProjectID, Version.ProjectID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerName, Version.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerDescription, Version.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerStatus, Version.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerDueDate, Version.DueDate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerCreatedOn, Version.CreatedOn); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(VerUpdatedOn, Version.UpdatedOn); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func VersionFromState(d *schema.ResourceData) *redmine.Version {
	Version := &redmine.Version{}
	Version.ProjectID, _ = d.Get(VerProjectID).(int)
	Version.Name = d.Get(VerName).(string)
	Version.Description = d.Get(VerDescription).(string)
	Version.Status = d.Get(VerStatus).(string)
	Version.DueDate = d.Get(VerDueDate).(string)
	Version.CreatedOn = d.Get(VerCreatedOn).(string)
	Version.UpdatedOn = d.Get(VerUpdatedOn).(string)

	VersionID := d.Id()
	if VersionID != "" && VersionID != "0" {
		Version.ID = VersionID
	}

	return Version
}
