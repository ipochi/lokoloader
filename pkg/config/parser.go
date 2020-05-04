package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/hclparse"
	clusterconfig "github.com/ipochi/lokoloader/pkg/cluster/config"
	//"github.com/ipochi/lokoloader/pkg/platform"
)

var configSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "variable", LabelNames: []string{"name"}},
		{Type: "backend", LabelNames: []string{"name"}},
		{Type: "cluster", LabelNames: []string{"name"}},
		{Type: "components", LabelNames: []string{"name"}},
	},
}

const (
	configFileExt = ".lokocfg"
	varFileExt    = ".vars"
)

type Parser struct {
	*hclparse.Parser
}

func NewParser() *Parser {
	return &Parser{
		Parser: hclparse.NewParser(),
	}
}

func Parse(path, varFiles string) (*clusterconfig.LokomotiveConfig, hcl.Diagnostics) {
	p := NewParser()
	cfg, diags := p.parse(path, varFiles)
	if diags.HasErrors() {
		return nil, diags
	}

	lokocfg, diags := p.getLokomotiveConfig(cfg)
	if diags.HasErrors() {
		return nil, diags
	}

	return lokocfg, diags
}

func (p *Parser) parse(configFilePath, varFilePath string) (*HCLLokomotiveConfig, hcl.Diagnostics) {

	var files []*hcl.File
	var diags hcl.Diagnostics

	hclFiles, diags := getFilesFromPath(configFilePath, configFileExt)
	if diags.HasErrors() {
		return nil, diags
	}

	for _, filename := range hclFiles {
		f, moreDiags := p.ParseHCLFile(filename)
		diags = append(diags, moreDiags...)
		files = append(files, f)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	cfg := &HCLLokomotiveConfig{}

	// First load variables block so that they can be
	// referenced later in the parsing of backend,
	// cluster,component blocks.
	for _, file := range files {
		diags = append(diags, cfg.decodeVariables(file)...)
	}

	// Load and parse var files.
	hclVarFiles, diags := getFilesFromPath(varFilePath, varFileExt)
	if diags.HasErrors() {
		return nil, diags
	}
	var varFiles []*hcl.File

	for _, filename := range hclVarFiles {
		f, moreDiags := p.ParseHCLFile(filename)
		diags = append(diags, moreDiags...)
		varFiles = append(varFiles, f)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	diags = append(diags, cfg.collectVariableValues(varFiles)...)
	if diags.HasErrors() {
		return nil, diags
	}

	// Decode the rest of the blocks in lokocfg.
	for _, file := range files {
		diags = append(diags, p.decodeConfig(file, cfg)...)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return cfg, diags
}

func (p *Parser) decodeConfig(file *hcl.File, cfg *HCLLokomotiveConfig) hcl.Diagnostics {
	var diags hcl.Diagnostics

	body := dynblock.Expand(file.Body, cfg.EvalContext())
	content, moreDiags := body.Content(configSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case "backend":
			backend, moreDiags := p.decodeBackend(block)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}

			cfg.Backend = backend

		case "cluster":
			cluster, moreDiags := p.decodeCluster(block)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}

			cfg.Cluster = cluster
		}
	}

	return diags
}

func (p *Parser) getLokomotiveConfig(cfg *HCLLokomotiveConfig) (
	*clusterconfig.LokomotiveConfig,
	hcl.Diagnostics) {
	var diags hcl.Diagnostics

	lokocfg := &clusterconfig.LokomotiveConfig{}

	if cfg.Cluster != nil {
		platform, moreDiags := p.decodeToPlatform(cfg.Cluster, cfg.EvalContext())
		diags = append(diags, moreDiags...)
		lokocfg.Platform = platform
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return lokocfg, diags
}
