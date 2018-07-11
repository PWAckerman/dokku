package main

import (
	"os"
	"strings"

	dockeroptions "../../"
	"github.com/dokku/dokku/plugins/common"
)

// docker_options_main_cmd() {
// 	declare desc="Display applications docker options"
// 	local cmd="docker-options"
// 	dokku_log_warn "Deprecated: Please use docker-options:report"

// 	verify_app_name "$2" && local APP="$2"
// 	read -ra passed_phases <<< "$(get_phases "$3")"
// 	if [[ ! "${passed_phases[@]}" ]]; then
// 	  display_all_phases_options
// 	else
// 	  display_passed_phases_options passed_phases[@]
// 	fi
//   }

func main() {
	cmd := os.Args[1]
	app := os.Args[2]
	err := common.VerifyAppName(app)
	if err != nil {
		common.LogFail(err.Error())
	}
	phase := os.Args[3]
	phases := dockeroptions.GetPhases(strings.Split(phase, ","))
	if len(phases) == 0 {
		dockeroptions.DisplayAllPhaseOptions(app)
	} else {
		dockeroptions.DisplayPassedPhasesOptions(phases...)
	}
}
