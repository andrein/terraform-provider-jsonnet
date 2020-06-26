package jsonnet

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)
import "github.com/google/go-jsonnet"

func dataSourceJsonnetTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJsonnetTemplateRead,

		Schema: map[string]*schema.Schema{
			"jsonnet": {
				Type:        schema.TypeString,
				Description: "The Jsonnet input",
				Required:    true,
			},
			"jpath": {
				Type:        schema.TypeList,
				Description: "The Jsonnet additional library search dir",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ext_var": {
				Type:        schema.TypeMap,
				Description: "A map of Jsonnet external variables, as strings",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ext_code": {
				Type:        schema.TypeMap,
				Description: "A map of Jsonnet external variables, as code",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tla_var": {
				Type:        schema.TypeMap,
				Description: "A map of Jsonnet top-level arguments, as strings",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tla_code": {
				Type:        schema.TypeMap,
				Description: "A map of Jsonnet top-level arguments, as code",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"json": {
				Type:        schema.TypeString,
				Description: "The JSON output",
				Computed:    true,
			},
		},
	}
}

func dataSourceJsonnetTemplateRead(d *schema.ResourceData, _ interface{}) error {
	vm := jsonnet.MakeVM()

	if jpath := d.Get("jpath").([]interface{}); jpath != nil {
		importer := jsonnet.FileImporter{JPaths: ExpandJPath(jpath)}
		vm.Importer(&importer)
	}

	if extVar := d.Get("ext_var").(map[string]interface{}); extVar != nil {
		for key, val := range extVar {
			vm.ExtVar(key, val.(string))
		}
	}

	if extCode := d.Get("ext_code").(map[string]interface{}); extCode != nil {
		for key, val := range extCode {
			vm.ExtCode(key, val.(string))
		}
	}

	if tlaVar := d.Get("tla_var").(map[string]interface{}); tlaVar != nil {
		for key, val := range tlaVar {
			vm.TLAVar(key, val.(string))
		}
	}

	if tlaCode := d.Get("tla_code").(map[string]interface{}); tlaCode != nil {
		for key, val := range tlaCode {
			vm.TLACode(key, val.(string))
		}
	}

	json, err := vm.EvaluateSnippet("input", d.Get("jsonnet").(string))
	if err != nil {
		return err
	}

	if err = d.Set("json", json); err != nil {
		return err
	}

	sha := sha256.Sum256([]byte(json))
	d.SetId(hex.EncodeToString(sha[:]))

	return nil
}
