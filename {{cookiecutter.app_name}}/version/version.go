package version

import (
	"fmt"
	"runtime"
)

// GitCommit returns the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version returns the main version number that is being run at the moment.
const Version = "0.1.0"

// BuildDate returns the date the binary was built
var BuildDate = ""

// GoVersion returns the version of the go runtime used to compile the binary
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

// AppName returns the name of the app
var AppName = "{{cookiecutter.app_name}}"

// Branch returns the branch name
var Branch = ""

// V is a struct that contains all the version information
type V struct {
	GitCommit string `json:"git_commit"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OSArch    string `json:"os_arch"`
	AppName   string `json:"app"`
	Branch    string `json:"branch"`
}

// Get returns a V struct with all the version information
// This is used to populate the version endpoint of the API
func Get() V {
	return V{
		GitCommit: GitCommit,
		Version:   Version,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		OSArch:    OsArch,
		AppName:   AppName,
		Branch:    Branch,
	}
}
