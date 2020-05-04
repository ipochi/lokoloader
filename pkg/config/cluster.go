package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/ipochi/lokoloader/pkg/platform"
	"github.com/ipochi/lokoloader/pkg/platform/packet"
	//	"github.com/zclconf/go-cty/cty"
)

var clusterSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "worker_pool", LabelNames: []string{"name"}},
	},
}

type WorkerPoolBlock struct {
	Name string `hcl:"name,label"`
	HCL2Ref
}

type ClusterBlock struct {
	Name             string
	WorkerPoolBlocks []*WorkerPoolBlock
	HCL2Ref
}

func (p *Parser) decodeWorkerPool(block *hcl.Block) (*WorkerPoolBlock, hcl.Diagnostics) {
	var b struct {
		Rest hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)
	if diags.HasErrors() {
		return nil, diags
	}

	wp := &WorkerPoolBlock{
		Name:    block.Labels[0],
		HCL2Ref: newHCL2Ref(block, b.Rest),
	}

	return wp, diags
}

func (p *Parser) decodeCluster(block *hcl.Block) (*ClusterBlock, hcl.Diagnostics) {
	type wp struct {
		Name string   `hcl:"name,label"`
		Rest hcl.Body `hcl:",remain"`
	}

	var b struct {
		Workerpool []wp     `hcl:"worker_pool,block"`
		Rest       hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)
	if diags.HasErrors() {
		return nil, diags
	}

	cluster := &ClusterBlock{
		Name:    block.Labels[0],
		HCL2Ref: newHCL2Ref(block, b.Rest),
	}

	fmt.Println("How many wpools --- ", len(b.Workerpool))
	for _, w := range b.Workerpool {
		diags := gohcl.DecodeBody(w.Rest, nil, &w)
		if diags.HasErrors() {
			return nil, diags
		}
		wp := &WorkerPoolBlock{
			Name:    w.Name,
			HCL2Ref: newHCL2Ref(block, w.Rest),
		}
		cluster.WorkerPoolBlocks = append(cluster.WorkerPoolBlocks, wp)
	}
	//content, moreDiags := b.Rest.Content(clusterSchema)
	//diags = append(diags, moreDiags...)

	//for _, block := range content.Blocks {
	//	switch block.Type {
	//	case "worker_pool":
	//		wp, moreDiags := p.decodeWorkerPool(block)
	//		diags = append(diags, moreDiags...)
	//		if moreDiags.HasErrors() {
	//			continue
	//		}
	//		cluster.WorkerPoolBlocks = append(cluster.WorkerPoolBlocks, wp)

	//	}
	//}

	return cluster, diags
}

func (p *Parser) decodeToPlatform(cb *ClusterBlock, ctx *hcl.EvalContext) (platform.Platform, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var pform platform.Platform
	pform = packet.NewPacketConfig()
	flatPlatformConfig, moreDiags := decodeHCL2Spec(cb.HCL2Ref.Rest, ctx, pform)
	diags = append(diags, moreDiags...)
	err := pform.LoadConfig(flatPlatformConfig)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Error in LoadConfig from Decode",
			Detail:   err.Error(),
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return pform, diags
}
