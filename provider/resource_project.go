package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inherit_members": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tracker_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"enabled_module_names": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"boards", "calendar", "documents", "files", "gantt", "issue_tracking", "news", "repository", "time_tracking", "wiki",
					},
						false),
				},
				Optional: true,
			},
		},
	}
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	projectID := d.Get("id").(string)

	client := i.(Client)
	project, err := client.ReadProject(ctx, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.GetID())

	return diags
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	project := projectFromState(d)

	err := client.CreateProject(ctx, project)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.GetID())
	return resourceProjectRead(ctx, d, i)
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func projectSetToState(project redmine.Project, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("id", project.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identifier", project.Identifier); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", project.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_public", project.IsPublic); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_id", project.ParentID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inherit_members", project.InheritMembers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tracker_ids", project.TrackerIDs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled_module_names", project.EnabledModuleNames); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func projectFromState(d *schema.ResourceData) *redmine.Project {
	project := &redmine.Project{}
	println(1)
	project.ID, _ = strconv.Atoi(d.Get("id").(string))
	println(12)
	project.Name = d.Get("name").(string)
	println(13)
	project.Identifier = d.Get("identifier").(string)
	println(14)
	project.Description = d.Get("description").(string)
	println(15)
	project.IsPublic = d.Get("is_public").(bool)
	println(16)
	project.ParentID = d.Get("parent_id").(string)
	println(17)
	project.InheritMembers = d.Get("inherit_members").(bool)
	println(18, d.Get("tracker_ids").([]interface{}))
	project.TrackerIDs = toStringSlice(d.Get("tracker_ids").([]interface{}))
	println(19)
	project.EnabledModuleNames = toStringSlice(d.Get("enabled_module_names").([]interface{}))
	println(10)

	return project
}

func toStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for _, item := range slice {
		resultItem := item.(string)
		result = append(result, resultItem)
	}

	return result
}
