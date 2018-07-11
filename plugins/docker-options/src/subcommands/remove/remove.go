package main

import (
	"os"
	"strings"

	// "github.com/dokku/dokku/plugins/common"
	dockeroptions "../../"
	"github.com/dokku/dokku/plugins/common"
)

// docker_options_remove_cmd() {
// 	declare desc="Remove a docker option from application"
// 	local cmd="docker-options:remove"

// 	verify_app_name "$2" && local APP="$2"
// 	read -ra passed_phases <<< "$(get_phases "$3")"
// 	shift 3 # everything else passed is the docker option
// 	# shellcheck disable=SC2154
// 	[[ -z ${passed_docker_option="$@"} ]] && dokku_log_fail "Please specify docker options to remove from the phase"
// 	remove_passed_docker_option passed_phases[@] "${passed_docker_option[@]}"
//   }

//   docker_options_remove_cmd "$@"

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
	dockeroptions.RemovePassedDockerOption(strings.Join(options, ""), app, phases...)
}
