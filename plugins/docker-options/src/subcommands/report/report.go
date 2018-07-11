package main

import (
	"os"

	dockeroptions "../../"
	"github.com/dokku/dokku/plugins/common"
)

// cmd-docker-options-report "$@"

func main() {
	cmd := os.Args[1]
	app := os.Args[2]
	flag := os.Args[3]
	if common.IsOption(app) {
		flag = app
		dockeroptions.CommandReportAll("", flag)
	} else {
		dockeroptions.CommandReportAll(app, flag)
	}
}
