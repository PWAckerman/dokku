package dockeroptions

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dokku/dokku/plugins/common"

	. "github.com/onsi/gomega"
)

var (
	testAppName = "test-app-2"
	dokkuRoot   = common.MustGetEnv("DOKKU_ROOT")
	testAppDir  = strings.Join([]string{dokkuRoot, testAppName}, "/")
)

func setupTestApp() (err error) {
	Expect(os.MkdirAll(testAppDir, 0766)).To(Succeed())
	b := []byte("export testKey=TESTING\n")
	if err = ioutil.WriteFile(strings.Join([]string{testAppDir, "/ENV"}, ""), b, 0644); err != nil {
		return
	}

	return
}

func teardownTestApp() {
	os.RemoveAll(testAppDir)
}

func TestDockerOptionsGenPhaseFilePath(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	path1 := GenPhaseFilePath(testAppName, "build")
	path2 := GenPhaseFilePath(testAppName, "deploy")
	path3 := GenPhaseFilePath(testAppName, "run")
	expPath1 := "/home/dokku/test-app-2/DOCKER_OPTIONS_build"
	expPath2 := "/home/dokku/test-app-2/DOCKER_OPTIONS_deploy"
	expPath3 := "/home/dokku/test-app-2/DOCKER_OPTIONS_run"
	Expect(path1).To(Equal(expPath1), "should return expected build path")
	Expect(path2).To(Equal(expPath2), "should return expected deploy path")
	Expect(path3).To(Equal(expPath3), "should return expected run path")
}

func TestDockerOptionsGetPhaseFilePath(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	path, err := GetPhaseFilePath(testAppName, "build")
	expPath := "/home/dokku/test-app-2/DOCKER_OPTIONS_build"
	Expect(path).To(Equal(expPath), "Proper path should return if app and phase are valid.")
	Expect(err).To(BeNil(), "There should be no error if proper path is formed.")
	// Expect(GetPhaseFilePath("not-test-app", "build")).To(Panic(), "If the app name is not valid, should panic.")
}

func TestDockerOptionsGetPhases(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	phases := GetPhases([]string{"build", "run", "deploy"})
	eq := reflect.DeepEqual(phases, []string{"build", "run", "deploy"})
	Expect(eq).To(BeTrue(), "All valid phases should be returned if passed.")
	phases2 := GetPhases([]string{"build", "run", "deproy"})
	eq2 := reflect.DeepEqual(phases2, []string{"build", "run"})
	Expect(eq2).To(BeTrue(), "Invalid phases should be filtered out.")
}

func FileExists(path string) bool {
	var exists bool
	if _, err := os.Stat(path); err == nil {
		exists = true
	} else {
		exists = false
	}
	return exists
}

func TestCreatePhaseFileIfRequired(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	path, _ := GetPhaseFilePath(testAppName, "build")
	var exists bool
	exists = FileExists(path)
	Expect(exists).To(BeFalse(), "File shouldn't exist until created.")
	CreatePhaseFileIfRequired(path)
	exists = FileExists(path)
	Expect(exists).To(BeTrue(), "File should exist once created")
	fi1, _ := os.Stat(path)
	CreatePhaseFileIfRequired(path)
	fi2, _ := os.Stat(path)
	Expect(os.SameFile(fi1, fi2)).To(BeTrue(), "File should only be created once.")
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsDisplayPhaseOptions(t *testing.T) {
	RegisterTestingT(t)
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsDisplayAllPhaseOptions(t *testing.T) {
	RegisterTestingT(t)
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsDisplayPassedPhasesOptions(t *testing.T) {
	RegisterTestingT(t)
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsAddPassedDockerOption(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	AddPassedDockerOption("--add-host=docker:x.x.x.1", testAppName, "build")
	AddPassedDockerOption("--squash", testAppName, "build")
	AddPassedDockerOption("--not-a-real-option", testAppName, "build", "deploy")
	path1, _ := GetPhaseFilePath(testAppName, "build")
	path2, _ := GetPhaseFilePath(testAppName, "deploy")
	lines1, _ := common.FileToSlice(path1)
	lines2, _ := common.FileToSlice(path2)
	Expect(lines1[0]).To(Equal("--add-host=docker:x.x.x.1"), "Options should persist in file.")
	Expect(lines1[1]).To(Equal("--squash"), "Options should be in order.")
	Expect(lines1[2] == lines2[0]).To(BeTrue(), "Options should populate in multiple phase files if phases are passed.")
}

func TestDockerOptionsRemovePassedDockerOption(t *testing.T) {
	RegisterTestingT(t)
	Expect(setupTestApp()).To(Succeed())
	defer teardownTestApp()
	path1, _ := GetPhaseFilePath(testAppName, "build")
	AddPassedDockerOption("--add-host=docker:x.x.x.1", testAppName, "build")
	RemovePassedDockerOption("--add-host=docker:x.x.x.1", testAppName, "build")
	lines1, _ := common.FileToSlice(path1)
	Expect(len(lines1)).To(Equal(0), "Specified option should be removed from file.")
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsCommandReportAll(t *testing.T) {
	RegisterTestingT(t)
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsReportSingle(t *testing.T) {
	RegisterTestingT(t)
	// Expect(MustGetEnv("DOKKU_ROOT")).To(Equal("/home/dokku"))
}

func TestDockerOptionsDockerOptions(t *testing.T) {
	RegisterTestingT(t)
	// Expect(setupTestApp()).To(Succeed())
	// defer teardownTestApp()
	// path1, _ := GetPhaseFilePath(testAppName, "build")
	// AddPassedDockerOption("--add-host=docker:x.x.x.1", testAppName, "build")
	// options := DockerOptions(testAppName, "build")

}
