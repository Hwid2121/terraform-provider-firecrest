package provider

// import (
// 	"testing"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// )

// func TestAccFireDataSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAcctestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// 				data "firecrest_fire` "test" {
// 					name = "example"
// 				},
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.firecrest_fire.test", "name", "example"),
// 					resource.TestCheckModuleResourceAttrSet("data.firecrest_fire.test", "token"),
// 					resource.TestCheckResourceAttr("data.firecrest_fire.test", "id", "placeholder"),
// 				),
// 			},
// 		},
// 	})
// }