package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jtopjian/terraform-provisioner-waffles/waffles"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: func() terraform.ResourceProvisioner {
			return new(waffles.ResourceProvisioner)
		},
	})
}
