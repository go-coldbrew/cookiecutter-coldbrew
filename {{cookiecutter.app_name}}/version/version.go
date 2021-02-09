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

type V struct {
	GitCommit string `json:"git_commit"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OSArch    string `json:"os_arch"`
}

func Get() V {
	return V{
		GitCommit: GitCommit,
		Version:   Version,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		OSArch:    OsArch,
	}
}
