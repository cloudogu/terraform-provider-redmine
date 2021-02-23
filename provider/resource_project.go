package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	PrjID             = "id"
	PrjName           = "name"
	PrjIdentifier     = "identifier"
	PrjDescription    = "description"
	PrjHomepage       = "homepage"
	PrjIsPublic       = "is_public"
	PrjParentID       = "parent_id"
	PrjInheritMembers = "inherit_members"
	PrjCreatedOn      = "created_on"
	PrjUpdatedOn      = "updated_on"
)

type ProjectClient interface {
	CreateProject(ctx context.Context, project *redmine.Project) (*redmine.Project, error)
	ReadProject(ctx context.Context, identifier string) (*redmine.Project, error)
	UpdateProject(ctx context.Context, project *redmine.Project) (*redmine.Project, error)
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			PrjID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PrjName: {
				Type:     schema.TypeString,
				Required: true,
			},
			PrjIdentifier: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			PrjDescription: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			PrjHomepage: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			PrjIsPublic: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			PrjParentID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			PrjInheritMembers: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			PrjCreatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			PrjUpdatedOn: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	projectID := d.Get(PrjID).(string)

	client := i.(ProjectClient)
	project, err := client.ReadProject(ctx, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	return projectSetToState(project, d)
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(ProjectClient)

	project := projectFromState(d)

	createdProject, err := client.CreateProject(ctx, project)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdProject.ID)

	diagRead := resourceProjectRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(ProjectClient)

	project := projectFromState(d)

	_, err := client.UpdateProject(ctx, project)
	if err != nil {
		return diag.FromErr(err)
	}

	diagRead := resourceProjectRead(ctx, d, i)
	diags = append(diags, diagRead...)

	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func projectSetToState(project *redmine.Project, d *schema.ResourceData) diag.Diagnostics {
	d.SetId(project.ID)
	if err := d.Set(PrjName, project.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjIdentifier, project.Identifier); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjDescription, project.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjIsPublic, project.IsPublic); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjParentID, project.ParentID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjInheritMembers, project.InheritMembers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjCreatedOn, project.CreatedOn); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(PrjUpdatedOn, project.UpdatedOn); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func projectFromState(d *schema.ResourceData) *redmine.Project {
	project := &redmine.Project{}
	project.ID, _ = d.Get(PrjID).(string)
	project.Name = d.Get(PrjName).(string)
	project.Identifier = d.Get(PrjIdentifier).(string)
	project.Description = d.Get(PrjDescription).(string)
	project.Homepage = d.Get(PrjHomepage).(string)
	project.IsPublic = d.Get(PrjIsPublic).(bool)
	project.ParentID = d.Get(PrjParentID).(string)
	project.InheritMembers = d.Get(PrjInheritMembers).(bool)
	project.CreatedOn = d.Get(PrjCreatedOn).(string)
	project.UpdatedOn = d.Get(PrjUpdatedOn).(string)

	return project
}
