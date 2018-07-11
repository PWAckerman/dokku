package dockeroptions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

type dockerPhase int

const (
	dBuild dockerPhase = iota
	dDeploy
	dRun
)

//PhaseMap is a mapping of phase name to docker phase const
var PhaseMap = func() map[string]dockerPhase {
	phaseMap := make(map[string]dockerPhase)
	phaseMap["build"] = dBuild
	phaseMap["deploy"] = dDeploy
	phaseMap["run"] = dRun
	return phaseMap
}()

//GenPhaseFilePath  Generates partial file path for a specific app and phase
func GenPhaseFilePath(app string, phase string) string {
	dokkuRoot := common.MustGetEnv("DOKKU_ROOT")
	values := []interface{}{dokkuRoot, app, phase}
	phaseFilePath := fmt.Sprintf(`%s/%s/DOCKER_OPTIONS_%s`, values...)
	return phaseFilePath
}

// GetPhaseFilePath gets the Phase File path for a particular phase and application
func GetPhaseFilePath(app string, phase string) (str string, err error) {
	err = common.VerifyAppName(app)
	if err != nil {
		common.LogFail("Error: phase_file_path is incomplete.")
		return "", err
	}
	return GenPhaseFilePath(app, phase), nil
}

// GetPhases gets a list of all possible phases as strings
func GetPhases(phases []string) []string {
	keys := []string{}
	outPhases := []string{}
	for key := range PhaseMap {
		keys = append(keys, key)
	}
	for _, el := range phases {
		if _, ok := PhaseMap[el]; !ok {
			msg := fmt.Sprintf("Phase(s) must be one of %s", strings.Join(keys[:], ","))
			fmt.Println(msg)
		} else {
			outPhases = append(outPhases, el)
		}
	}
	return outPhases
}

//CreatePhaseFileIfRequired creates a phase file if it doesn't already exists
func CreatePhaseFileIfRequired(path string) {
	exists := common.FileExists(path)
	if !exists {
		f, err := os.Create(path)
		if err != nil {
			common.LogFail(err.Error())
		}
		f.Close()
	}
}

//DisplayPhaseOptions prints the options for a particular phase to the terminal
func DisplayPhaseOptions(phase string, path string) {
	fmt.Sprintf("%s options:", phase)
	lines, _ := common.FileToSlice(path)
	for ln := range lines {
		fmt.Print(ln)
	}
}

//DisplayAllPhaseOptions prints the options for each phase
func DisplayAllPhaseOptions(appName string) {
	keys := []string{}
	for key := range PhaseMap {
		keys = append(keys, key)
	}
	for _, phase := range keys {
		path, _ := GetPhaseFilePath(appName, phase)
		if common.FileExists(path) {
			DisplayPhaseOptions(phase, path)
		}
	}
}

//DisplayPassedPhasesOptions displays the options only for those phases passed
func DisplayPassedPhasesOptions(appName string, phases ...string) {
	for _, phase := range phases {
		path, _ := GetPhaseFilePath(appName, phase)
		if common.FileExists(path) {
			fmt.Sprintf("%s options: none", path)
		} else {
			DisplayPhaseOptions(phase, path)
		}
	}
}

//AddPassedDockerOption adds an option to the phase files specified
func AddPassedDockerOption(option string, appName string, phases ...string) {
	for _, phase := range phases {
		path, err := GetPhaseFilePath(appName, phase)
		if err != nil {
			common.LogFail(err.Error())
		}
		CreatePhaseFileIfRequired(path)
		slices, err := common.FileToSlice(path)
		var strBuilder bytes.Buffer
		strBuilder.WriteString(option)
		strBuilder.WriteString("\n")
		fmt.Println(strBuilder.String())
		slices = append(slices, strBuilder.String())
		dat := common.RemoveBlankLines(strings.Join(slices, "\n"))
		err = ioutil.WriteFile(path, []byte(dat), 0777)
		if err != nil {
			common.LogFail(err.Error())
		}
	}
}

//RemovePassedDockerOption removes the specified option from the phases passed to the function
func RemovePassedDockerOption(option string, appName string, phases ...string) {
	for _, phase := range phases {
		path, _ := GetPhaseFilePath(appName, phase)
		if common.FileExists(path) {
			options, err := common.FileToSlice(path)
			out := make([]string, 0)
			for _, opt := range options {
				if opt != option {
					out = append(out, option+"\n")
				}
			}
			err = ioutil.WriteFile(path, []byte(strings.Join(out, "")), 0777)
			if err != nil {
				common.LogFail(err.Error())
			}
		}
	}
}

//CommandReportAll prints all options for the specified app (or for all apps if not specified)
func CommandReportAll(appName string, infoFlag string) {
	if appName == "" && infoFlag == "" {
		infoFlag = "true"
	}
	if appName == "" {
		apps, err := common.DokkuApps()
		if err != nil {
			common.LogFail(err.Error())
		}
		for _, app := range apps {
			ReportSingle(app, infoFlag)
		}
	} else {
		ReportSingle(appName, infoFlag)
	}
}

//ReportSingle prints the options for all phases (or specified phase) for a provided app
func ReportSingle(appName string, infoFlag string) (str string, err error) {
	if infoFlag == "true" {
		infoFlag = ""
	}
	err = common.VerifyAppName(appName)
	if err != nil {
		common.LogFail("Error: phase_file_path is incomplete.")
		return "", err
	}
	var flagMap = make(map[string]string)
	flagMap["--docker-options-build"] = DockerOptions(appName, "build")
	flagMap["--docker-options-deploy"] = DockerOptions(appName, "deploy")
	flagMap["--docker-options-run"] = DockerOptions(appName, "run")
	if len(infoFlag) == 0 {
		common.LogInfo2Quiet(fmt.Sprintf("%s docker options information", appName))
		for key, options := range flagMap {
			str := common.FormatFlagString(key)
			res := str + options
			common.LogVerbose(res)
			return res, nil
		}
	} else {
		match := false
		valueExists := false
		validFlags := ""
		for flag := range flagMap {
			validFlags = validFlags + " " + flag
			if flag == infoFlag {
				if len(flag) != 0 {
					valueExists = true
					match = true
				} else {
					match = true
				}
			}
		}
		if match != true {
			common.LogFail(fmt.Sprintf("Invalid flag passed, valid flags: %s", validFlags))
		}
		if valueExists != true {
			common.LogFail("not deployed")
		}
	}
	return "", nil
}

//DockerOptions retrieves all of the options for a specific app and phases
func DockerOptions(appName string, phase string) string {
	path, _ := GetPhaseFilePath(appName, phase)
	exists := common.FileExists(path)
	if exists {
		options, _ := common.FileToSlice(path)
		return common.RemoveNewLines(strings.Join(options, ""))
	} else {
		return ""
	}
}
