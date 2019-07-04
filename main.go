package main

import (
	"flag"

	"github.com/yolean/k8s-endpoints-across-clusters/pkg/cluster"
)

func main() {
	var configsdir string

	flag.StringVar(&configsdir, "configsdir", "/etc/kubeconfigs", "Directory that contains the kubeconfig files")
	flag.Parse()

	kubeconfig := configsdir + "/cluster1"

	cluster1 := cluster.New(kubeconfig)
	cluster1.Start()

	// Wait forever
	select {}
}
