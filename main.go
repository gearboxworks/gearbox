package main

import (
	"fmt"
	"gearbox/cmd"
	"gearbox/gearbox"
	"os"
)

//go:generate go-bindata -dev -o gearbox/asset.go -pkg gearbox admin/...

func main() {
	gearbox.WriteAdminAssetsToWebRoot(gearbox.Host)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
