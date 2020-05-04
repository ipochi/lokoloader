package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type Variable struct {
	Name         string
	Description  string
	Type         cty.Type
	DefaultValue cty.Value
	Value        cty.Value
}

type Variables map[string]*Variable

func (v Variables) Values() map[string]cty.Value {
	res := map[string]cty.Value{}

	for key, val := range v {
		res[key] = val.Value
	}

	return res
}

func (v *Variables) decodeVariableBlock(block *hcl.Block, ctx *hcl.EvalContext) hcl.Diagnostics {
	if (*v) == nil {
		(*v) = Variables{}
	}

	if _, found := (*v)[block.Labels[0]]; found {
		return []*hcl.Diagnostic{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate variable",
			Detail:   "Duplicate " + block.Labels[0] + " variable definition found.",
			Context:  block.DefRange.Ptr(),
		}}
	}

	var b struct {
		Description string   `hcl:"description,optional"`
		Rest        hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	if diags.HasErrors() {
		return diags
	}

	name := block.Labels[0]

	res := &Variable{
		Name:        name,
		Description: b.Description,
	}

	attrs, moreDiags := b.Rest.JustAttributes()
	diags = append(diags, moreDiags...)

	if t, ok := attrs["type"]; ok {
		delete(attrs, "type")
		tp, moreDiags := typeexpr.Type(t.Expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return diags
		}

		res.Type = tp
	}

	if def, ok := attrs["default"]; ok {
		delete(attrs, "default")
		defaultValue, moreDiags := def.Expr.Value(ctx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return diags
		}

		if res.Type != cty.NilType {
			var err error
			defaultValue, err = convert.Convert(defaultValue, res.Type)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("This default value is not compatible with the variable's type constraint: %s.", err),
					Subject:  def.Expr.Range().Ptr(),
				})
				defaultValue = cty.DynamicVal
			}
		}

		res.DefaultValue = defaultValue
		// It's possible no type attribute was assigned so lets make
		// sure we have a valid type otherwise there will be issues parsing the value.
		if res.Type == cty.NilType {
			res.Type = res.DefaultValue.Type()
		}
	}

	(*v)[name] = res
	return diags
}
