//go:generate mapstructure-to-hcl2 -type config,flatcar,workerPool

package packet

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/ipochi/lokoloader/pkg/platform"
	"github.com/mitchellh/mapstructure"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"sort"
)

type network struct {
	//platform.NetworkConfig
	ManagementCIDRs []string `mapstructure:"management_cidrs"`
	NodePrivateCIDR string   `mapstructure:"node_private_cidr"`
}

type flatcar struct {
	platform.FlatcarConfig `mapstructure:",squash"`
	OSArch                 string `mapstructure:"os_arch,optional"`
	IPXEScriptURL          string `mapstructure:"ipxe_script_url,optional"`
}
type workerPool struct {
	Name    string
	Count   int `mapstructure:"count"`
	flatcar `mapstructure:",squash"`
}

type config struct {
	Name                   string
	platform.Metadata      `mapstructure:",squash"`
	Network                network `mapstructure:",squash"`
	Flatcar                flatcar `mapstructure:",squash"`
	platform.ClusterConfig `mapstructure:",squash"`
	WorkerPools            []workerPool `mapstructure:"worker_pool"`
	AuthToken              string       `mapstructure:"auth_token"`
	Facility               string       `mapstructure:"facility"`
}

func NewWorkerPool() *workerPool {
	return &workerPool{}
}
func NewPacketConfig() *config {
	return &config{}
}

func (c *config) GetData() {
	fmt.Println("Platform Name - ", c.Name)
	fmt.Println("Platform AssetDir - ", c.AssetDir)
	fmt.Println("Platform ClusterName - ", c.ClusterName)
	fmt.Println("Platform OSChannel- ", c.Flatcar.Channel)
	fmt.Println("Platform OSVersion - ", c.Flatcar.Version)
	fmt.Println("Platform OSArch - ", c.Flatcar.OSArch)
	fmt.Println("Platform DomainSuffix - ", c.ClusterDomainSuffix)
	fmt.Println("Platform AuthToken - ", c.AuthToken)
	fmt.Println("Platform Facility - ", c.Facility)

	fmt.Println("No Of Worker Pools :", len(c.WorkerPools))

	for _, wp := range c.WorkerPools {
		fmt.Println("WP Name - ", wp.Name)
		fmt.Println("WP Count - ", wp.Count)
		fmt.Println("WP OSVersion - ", wp.Version)
		fmt.Println("WP OSChanne; - ", wp.Channel)
	}
}

func (c *config) ConfigSpec() hcldec.ObjectSpec     { return c.FlatMapstructure().HCL2Spec() }
func (w *workerPool) ConfigSpec() hcldec.ObjectSpec { return w.FlatMapstructure().HCL2Spec() }

func (c *config) LoadConfig(raws cty.Value) error {
	err := Decode(c, raws)
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes the configuration into the target and optionally
// automatically interpolates all the configuration as it goes.
func Decode(target interface{}, raws ...interface{}) error {
	for i, raw := range raws {
		// check for cty values and transform them to json then to a
		// map[string]interface{} so that mapstructure can do its thing.
		cval, ok := raw.(cty.Value)
		if !ok {
			continue
		}
		type flatConfigurer interface {
			FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec }
		}
		ctarget := target.(flatConfigurer)
		flatCfg := ctarget.FlatMapstructure()
		err := gocty.FromCtyValue(cval, flatCfg)
		if err != nil {
			switch err := err.(type) {
			case cty.PathError:
				return fmt.Errorf("%v: %v", err, err.Path)
			}
			return err
		}
		b, err := ctyjson.SimpleJSONValue{Value: cval}.MarshalJSON()
		if err != nil {
			return err
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(b, &raw); err != nil {
			return err
		}
		raws[i] = raw
	}

	// Build our decoder
	var md mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           target,
		Metadata:         &md,
		WeaklyTypedInput: true,
	})

	if err != nil {
		return err
	}
	for _, raw := range raws {
		if err := decoder.Decode(raw); err != nil {

			fmt.Println("ot here ,,,, Hello", err)
			return err
		}
	}

	// If we have unused keys, it is an error
	if len(md.Unused) > 0 {
		sort.Strings(md.Unused)
		unusedStr := ""
		for _, unused := range md.Unused {
			unusedStr = unusedStr + fmt.Sprintf("Unused string ---- ", unused)
		}

		if unusedStr != "" {
			return fmt.Errorf("Gound unused strings ----\n%s", unusedStr)
		}
	}

	return nil
}
