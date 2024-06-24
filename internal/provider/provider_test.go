package provider

import (

	// "github.com/hashicorp/terraform-plugin-framework/providerserver"
	// "github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
// const (
// 	providerConfig = `
// 	provider "firecrest" {
// 	client_id     = "firecrest-ntafta-coder"
// 	client_secret = "D1wLfcA3BfVzxYA7eJ7AivIEklWNTH3C"
// 	}`
// )


// var (
// 	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
// 		"firecrest": providerserver.NewProtocol6WithError(New("test")()),
// 	}
// )