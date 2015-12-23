package waffles

import (
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceProvisioner_impl(t *testing.T) {
	var _ terraform.ResourceProvisioner = new(ResourceProvisioner)
}

/*
func TestResourceProvider_Apply(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"host":           "www.example.com",
		"role":           "test",
		"site_directory": "~/waffles",
	})

	output := new(terraform.MockUIOutput)
	p := new(ResourceProvisioner)
	if err := p.Apply(output, nil, c); err != nil {
		t.Fatalf("err: %v", err)
	}
}
*/

func TestResourceProvider_Validate_good(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"host":           "www.example.com",
		"role":           "test",
		"site_directory": "~/waffles",
	})

	p := new(ResourceProvisioner)
	warn, errs := p.Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}
}

func TestResourceProvider_Validate_missing(t *testing.T) {
	c := testConfig(t, map[string]interface{}{})
	p := new(ResourceProvisioner)
	warn, errs := p.Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) == 0 {
		t.Fatalf("Should have errors")
	}
}

func testConfig(
	t *testing.T,
	c map[string]interface{}) *terraform.ResourceConfig {
	r, err := config.NewRawConfig(c)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	return terraform.NewResourceConfig(r)
}
