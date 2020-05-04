package config

import (
	"github.com/hashicorp/hcl/v2"
)

type Component struct {
	Name  string
	block *hcl.Block
}

type Components map[string]*Component
