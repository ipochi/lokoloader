package main

import (
	"fmt"
	"github.com/ipochi/lokoloader/pkg/config"
)

func main() {

	lokocfg, diags := config.Parse("test.lokocfg", "test.vars")
	if diags.HasErrors() {
		for _, diag := range diags {
			fmt.Println("diag ---", diag.Error())
		}

		fmt.Println("thats it ... no lokocfg")
		return
	}

	lokocfg.Platform.GetData()
}
