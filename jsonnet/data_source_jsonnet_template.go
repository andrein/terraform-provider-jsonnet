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
