package provider

import (
	"context"
	"github.com/cloudogu/terraform-provider-redmine/redmine"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDMINE_URL", "http://localhost:3000/"),
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
			"redmine_project":        resourceProject(),
			"redmine_issue":          resourceIssue(),
			"redmine_issue_category": resourceIssueCategory(),
			"redmine_version":        resourceVersion(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skipVerify := d.Get("skip_cert_verify").(bool)

	var url string

	uVal, ok := d.GetOk("url")
	if ok {
		url = uVal.(string)
	}

	client, err := redmine.NewClient(redmine.Config{
		URL:            url,
		Username:       username,
		Password:       password,
		SkipCertVerify: skipVerify,
	})

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
