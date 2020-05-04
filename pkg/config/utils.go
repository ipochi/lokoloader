package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"os"
	"path/filepath"
)

func isDir(name string) (bool, error) {
	s, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	return s.IsDir(), nil
}

func getFilesFromPath(path, extension string) ([]string, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	isDir, err := isDir(path)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Could not determine if %s is a directory", path),
			Detail:   err.Error(),
		})

		return nil, diags
	}

	if !isDir {
		//TODO: check extension of the file before returning
		// throw error if not matching.
		return []string{path}, hcl.Diagnostics{}
	}

	var paths []string

	if isDir {
		globPattern := filepath.Join(path, extension)
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Could not match pattern with the extension %s", extension),
				Detail:   err.Error(),
			})

			return nil, diags
		}

		//TODO: check if any pattern matches the glob pattern
		// throw error if not.
		paths = append(paths, matches...)
	}

	return paths, hcl.Diagnostics{}
}
