package jsonnet

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"reflect"
	"testing"
)

func TestAccJsonnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonnetTemplate_basic,
				Check:  testCheckResourceJsonAttr("data.jsonnet_template.test", "json", `{"result": "test"}`),
			},
		},
	})
}

func TestAccJsonnet_jpath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonnetTemplate_jpath,
				Check:  testCheckResourceJsonAttr("data.jsonnet_template.test", "json", `{"result": "test"}`),
			},
		},
	})
}

func testCheckResourceJsonAttr(name, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual := s.RootModule().Resources[name].Primary.Attributes[key]
		expected := value

		var expectedJSONAsInterface, actualJSONAsInterface interface{}

		if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
			return fmt.Errorf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error())
		}

		if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
			return fmt.Errorf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error())
		}

		if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
			return fmt.Errorf("%s: Attribute '%s': \n"+
				"expected: %s\n"+
				"actual  : %s", name, key, expected, actual)
		}

		return nil
	}
}

var testAccJsonnetTemplate_basic = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  local result = "test";

  {
	result: result
  }
  EOF
}
`

var testAccJsonnetTemplate_jpath = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  local test = import "test.libsonnet";
  local result = "test";
  
  {} + test.withResult(result)
  EOF
  jpath = ["test-fixtures"]
}
`
