package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/umich-vci/terraform-provider-rhsm/rhsm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rhsm.Provider})
}
