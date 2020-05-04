package config

import (
	"github.com/hashicorp/hcl/v2"
)

type BackendBlock struct {
	Name  string
	block *hcl.Block
}

func (p *Parser) decodeBackend(block *hcl.Block) (*BackendBlock, hcl.Diagnostics) {

	return nil, hcl.Diagnostics{}
}
