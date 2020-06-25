# terraform-provider-jsonnet

Render Jsonnet templates in terraform.

## Example
```hcl-terraform
data "jsonnet_template" "pipeline" {
  jsonnet = file("pipeline.jsonnet")
  jpath = [
    "vendor"
  ]
}

output "rendered_json" {
  value = data.jsonnet_template.pipeline.json
}
```

## Installation

#### Build from source

```shell script
$ git clone https://github.com/andrein/terraform-provider-jsonnet.git
$ make install
```

#### Use released binaries

Download the latest [release](https://github.com/andrein/terraform-provider-jsonnet/releases/) to `~/.terraform.d/plugins`. See the [terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for more details.
