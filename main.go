package main

import (
	"context"

	"github.com/thecodeteam/gocsi"

	"github.com/yunify/kubernetes-qingstor-csi/provider"
	"github.com/yunify/kubernetes-qingstor-csi/service"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	gocsi.Run(
		context.Background(),
		service.Name,
		"Qingcloud CSI plugin",
		"CSI plugin for qingcloud storage resources",
		provider.New())
}
