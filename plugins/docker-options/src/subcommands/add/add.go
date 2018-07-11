package main

import (
	"os"
	"strings"

	dockeroptions "../../"
	"github.com/dokku/dokku/plugins/common"
)

// docker_options_add_cmd() {
// 	declare desc="Add a docker option to application"
// 	local cmd="docker-options:add"

// 	verify_app_name "$2" && local APP="$2"
// 	read -ra passed_phases <<< "$(get_phases "$3")"
// 	shift 3 # everything else passed is the docker option
// 	local passed_docker_option="$*"
// 	[[ -z "$passed_docker_option" ]] && dokku_log_fail "Please specify docker options to add to the phase"
// 	add_passed_docker_option passed_phases[@] "${passed_docker_option[@]}"
//   }

//   docker_options_add_cmd "$@"

func main() {
	cmd := os.Args[1]
	app := os.Args[2]
	err := common.VerifyAppName(app)
	if err != nil {
		common.LogFail(err.Error())
	}
	phase := os.Args[3]
	phases := dockeroptions.GetPhases(strings.Split(phase, ","))
	options := os.Args[3:]
	if len(options) == 0 {
		common.LogFail("Please specify docker options to add to the phase.")
	}
	dockeroptions.AddPassedDockerOption(strings.Join(options, ""), app, phases...)
}
