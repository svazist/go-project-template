package main

import (
	"fmt"
	"github.com/svazist/go-project-template/cmd"
	"github.com/svazist/go-project-template/metrics"
)

var (
	Version   = "dev"
	Build     = "none"
	BuildDate = "unknown"
)

func main() {

	metrics.Version = Version
	metrics.Build = Build
	metrics.BuildDate = BuildDate

	cmd.Version = fmt.Sprintf("%s, commit %s, built at %s", Version, Build, BuildDate)
	fmt.Printf("Development version:\n Version: %v, commit: %v, built at: %v\n\n", Version, Build, BuildDate)
	cmd.Execute()
}
