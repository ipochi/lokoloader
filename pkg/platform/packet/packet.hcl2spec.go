// Code generated by "mapstructure-to-hcl2 -type config,flatcar,workerPool"; DO NOT EDIT.
package packet

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// Flatconfig is an auto-generated flat version of config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type Flatconfig struct {
	Name                     *string          `cty:"name"`
	AssetDir                 *string          `mapstructure:"asset_dir" cty:"asset_dir"`
	ClusterName              *string          `mapstructure:"cluster_name" cty:"cluster_name"`
	ManagementCIDRs          []string         `mapstructure:"management_cidrs" cty:"management_cidrs"`
	NodePrivateCIDR          *string          `mapstructure:"node_private_cidr" cty:"node_private_cidr"`
	Channel                  *string          `mapstructure:"os_channel,optional" cty:"os_channel"`
	Version                  *string          `mapstructure:"os_version,optional" cty:"os_version"`
	OSArch                   *string          `mapstructure:"os_arch,optional" cty:"os_arch"`
	IPXEScriptURL            *string          `mapstructure:"ipxe_script_url,optional" cty:"ipxe_script_url"`
	ClusterDomainSuffix      *string          `mapstructure:"cluster_domain_suffix,optional" cty:"cluster_domain_suffix"`
	EnableAggregation        *bool            `mapstructure:"enable_aggregation,optional" cty:"enable_aggregation"`
	CertsValidityPeriodHours *int             `mapstructure:"certs_validity_period_hours,optional" cty:"certs_validity_period_hours"`
	WorkerPools              []FlatworkerPool `mapstructure:"worker_pool" cty:"worker_pool"`
	AuthToken                *string          `mapstructure:"auth_token" cty:"auth_token"`
	Facility                 *string          `mapstructure:"facility" cty:"facility"`
}

// FlatMapstructure returns a new Flatconfig.
// Flatconfig is an auto-generated flat version of config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(Flatconfig)
}

// HCL2Spec returns the hcl spec of a config.
// This spec is used by HCL to read the fields of config.
// The decoded values from this spec will then be applied to a Flatconfig.
func (*Flatconfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"name":                        &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"asset_dir":                   &hcldec.AttrSpec{Name: "asset_dir", Type: cty.String, Required: false},
		"cluster_name":                &hcldec.AttrSpec{Name: "cluster_name", Type: cty.String, Required: false},
		"management_cidrs":            &hcldec.AttrSpec{Name: "management_cidrs", Type: cty.List(cty.String), Required: false},
		"node_private_cidr":           &hcldec.AttrSpec{Name: "node_private_cidr", Type: cty.String, Required: false},
		"os_channel":                  &hcldec.AttrSpec{Name: "os_channel", Type: cty.String, Required: false},
		"os_version":                  &hcldec.AttrSpec{Name: "os_version", Type: cty.String, Required: false},
		"os_arch":                     &hcldec.AttrSpec{Name: "os_arch", Type: cty.String, Required: false},
		"ipxe_script_url":             &hcldec.AttrSpec{Name: "ipxe_script_url", Type: cty.String, Required: false},
		"cluster_domain_suffix":       &hcldec.AttrSpec{Name: "cluster_domain_suffix", Type: cty.String, Required: false},
		"enable_aggregation":          &hcldec.AttrSpec{Name: "enable_aggregation", Type: cty.Bool, Required: false},
		"certs_validity_period_hours": &hcldec.AttrSpec{Name: "certs_validity_period_hours", Type: cty.Number, Required: false},
		"worker_pool":                 &hcldec.BlockListSpec{TypeName: "worker_pool", Nested: hcldec.ObjectSpec((*FlatworkerPool)(nil).HCL2Spec())},
		"auth_token":                  &hcldec.AttrSpec{Name: "auth_token", Type: cty.String, Required: false},
		"facility":                    &hcldec.AttrSpec{Name: "facility", Type: cty.String, Required: false},
	}
	return s
}

// Flatflatcar is an auto-generated flat version of flatcar.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type Flatflatcar struct {
	Channel       *string `mapstructure:"os_channel,optional" cty:"os_channel"`
	Version       *string `mapstructure:"os_version,optional" cty:"os_version"`
	OSArch        *string `mapstructure:"os_arch,optional" cty:"os_arch"`
	IPXEScriptURL *string `mapstructure:"ipxe_script_url,optional" cty:"ipxe_script_url"`
}

// FlatMapstructure returns a new Flatflatcar.
// Flatflatcar is an auto-generated flat version of flatcar.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*flatcar) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(Flatflatcar)
}

// HCL2Spec returns the hcl spec of a flatcar.
// This spec is used by HCL to read the fields of flatcar.
// The decoded values from this spec will then be applied to a Flatflatcar.
func (*Flatflatcar) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"os_channel":      &hcldec.AttrSpec{Name: "os_channel", Type: cty.String, Required: false},
		"os_version":      &hcldec.AttrSpec{Name: "os_version", Type: cty.String, Required: false},
		"os_arch":         &hcldec.AttrSpec{Name: "os_arch", Type: cty.String, Required: false},
		"ipxe_script_url": &hcldec.AttrSpec{Name: "ipxe_script_url", Type: cty.String, Required: false},
	}
	return s
}

// FlatworkerPool is an auto-generated flat version of workerPool.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatworkerPool struct {
	Name  *string `cty:"name"`
	Count *int    `mapstructure:"count" cty:"count"`
}

// FlatMapstructure returns a new FlatworkerPool.
// FlatworkerPool is an auto-generated flat version of workerPool.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*workerPool) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatworkerPool)
}

// HCL2Spec returns the hcl spec of a workerPool.
// This spec is used by HCL to read the fields of workerPool.
// The decoded values from this spec will then be applied to a FlatworkerPool.
func (*FlatworkerPool) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"name":  &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"count": &hcldec.AttrSpec{Name: "count", Type: cty.Number, Required: false},
	}
	return s
}