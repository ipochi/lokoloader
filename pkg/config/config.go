package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
	"io/ioutil"
)

type HCLLokomotiveConfig struct {
	Backend    *BackendBlock
	Cluster    *ClusterBlock
	Components Components
	Variables  Variables
}

func (l *HCLLokomotiveConfig) decodeVariables(file *hcl.File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, moreDiags := file.Body.Content(configSchema)
	diags = append(diags, moreDiags...)

	variables := make(map[string]*Variable)
	l.Variables = variables

	for _, block := range content.Blocks {
		switch block.Type {
		case "variable":
			fmt.Println("Entered here ... ")
			moreDiags := l.Variables.decodeVariableBlock(block, nil)
			diags = append(diags, moreDiags...)
		}
	}

	return diags
}

func (l *HCLLokomotiveConfig) collectVariableValues(files []*hcl.File) hcl.Diagnostics {
	var diags hcl.Diagnostics
	variables := l.Variables

	fmt.Println("Got these as variables --- ", variables)
	for _, file := range files {
		attrs, moreDiags := file.Body.JustAttributes()
		diags = append(diags, moreDiags...)

		for name, attr := range attrs {
			variable, found := variables[name]
			if !found {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Undefined variable",
					Detail: fmt.Sprintf("A %q variable was set but was "+
						"not found in known variables", name),
					Context: attr.Range.Ptr(),
				})
				continue
			}

			val, moreDiags := attr.Expr.Value(nil)
			diags = append(diags, moreDiags...)

			if variable.Type != cty.NilType {
				var err error
				val, err = convert.Convert(val, variable.Type)
				if err != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid value for variable",
						Detail:   fmt.Sprintf("The value for %s is not compatible with the variable's type constraint: %s.", name, err),
						Subject:  attr.Expr.Range().Ptr(),
					})
					val = cty.DynamicVal
				}
			}

			variable.Value = val
		}
	}

	return diags
}

func (l *HCLLokomotiveConfig) EvalContext() *hcl.EvalContext {
	variables := l.Variables.Values()
	evalContext := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(variables),
		},
		Functions: map[string]function.Function{
			"pathexpand": evalFuncPathExpand(),
			"file":       evalFuncFile(),
		},
	}

	return evalContext
}

func evalFuncPathExpand() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: cty.String,
			}},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			expandedPath, err := homedir.Expand(args[0].AsString())
			return cty.StringVal(expandedPath), err
		},
	})
}

func evalFuncFile() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: cty.String,
			}},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			filePath := args[0].AsString()
			content, err := ioutil.ReadFile(filePath)
			return cty.StringVal(string(content)), err
		},
	})
}
