package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Client interface {
	CreateProject(ctx context.Context, project *redmine.Project) error
	ReadProject(ctx context.Context, id string) (*redmine.Project, error)
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_URL", "http://localhost:8080/"),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_USERNAME", "admin"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_PASSWORD", "admin"),
			},
			"skip_cert_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_SKIP_CERT_VERIFY", false),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"redmine_project": resourceProject(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"redmine_issue_statuses": dataSourceIssueStatuses(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skipVerify := d.Get("skip_cert_verify").(bool)
	apiKey := d.Get("api_key").(string)

	var url string

	uVal, ok := d.GetOk("url")
	if ok {
		url = uVal.(string)
	}

	client := redmine.NewClient(redmine.Config{
		URL:            url,
		Username:       username,
		Password:       password,
		SkipCertVerify: skipVerify,
		APIKey:         apiKey,
	})

	return client, nil
}
