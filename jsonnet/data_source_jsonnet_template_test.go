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

func TestAccJsonnet_ext(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonnetTemplate_ext_var,
				Check:  testCheckResourceJsonAttr("data.jsonnet_template.test", "json", `{"result": "test"}`),
			},
			{
				Config: testAccJsonnetTemplate_ext_code,
				Check:  testCheckResourceJsonAttr("data.jsonnet_template.test", "json", `{"result": "test"}`),
			},
		},
	})
}

func TestAccJsonnet_tla(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonnetTemplate_tla_var,
				Check:  testCheckResourceJsonAttr("data.jsonnet_template.test", "json", `{"result": "test"}`),
			},
			{
				Config: testAccJsonnetTemplate_tla_code,
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

var testAccJsonnetTemplate_ext_var = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  { result: std.extVar('result') }
  EOF
  ext_var = {
	result = "test"
  }
}
`

var testAccJsonnetTemplate_ext_code = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  { result: std.extVar('result') }
  EOF
  ext_code = {
	result = "'test'"
  }
}
`

var testAccJsonnetTemplate_tla_var = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  function(result) {
	result: result
  }
  EOF
  tla_var = {
	result = "test"
  }
}
`

var testAccJsonnetTemplate_tla_code = `
data "jsonnet_template" "test" {
  jsonnet = <<-EOF
  function(result) {
	result: result
  }
  EOF
  tla_code = {
	result = "'test'"
  }
}
`
