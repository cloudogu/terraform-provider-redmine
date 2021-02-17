package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"scm": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("REDMINE_URL"); err == "" {
		t.Fatal("REDMINE_URL must be set for acceptance tests")
	}
	if err := os.Getenv("REDMINE_USERNAME"); err == "" {
		t.Fatal("REDMINE_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("REDMINE_PASSWORD"); err == "" {
		t.Fatal("REDMINE_PASSWORD must be set for acceptance tests")
	}
}
