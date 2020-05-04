package platform

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// a struct (or type) implementing HCL2Speccer is a type that can tell it's own
// hcl2 conf/layout.
type HCL2Speccer interface {
	// ConfigSpec should return the hcl object spec used to configure the
	// builder. It will be used to tell the HCL parsing library how to
	// validate/configure a configuration.
	ConfigSpec() hcldec.ObjectSpec
}

type Platform interface {
	HCL2Speccer
	LoadConfig(cty.Value) error
	GetData()
}

type Metadata struct {
	AssetDir    string `mapstructure:"asset_dir"`
	ClusterName string `mapstructure:"cluster_name"`
}

type ClusterConfig struct {
	ClusterDomainSuffix      string `mapstructure:"cluster_domain_suffix,optional"`
	EnableAggregation        bool   `mapstructure:"enable_aggregation,optional"`
	CertsValidityPeriodHours int    `mapstructure:"certs_validity_period_hours,optional"`
}

type FlatcarConfig struct {
	Channel string `mapstructure:"os_channel,optional"`
	Version string `mapstructure:"os_version,optional"`
}

type NetworkConfig struct {
	NetworkMTU      string `mapstructure:"network_mtu,optional"`
	PodCIDR         string `mapstructure:"pod_cidr,optional"`
	ServiceCIDR     string `mapstructure:"service_cidr,optional"`
	EnableReporting bool   `mapstructure:"enable_reporting,optional"`
}
