package version

import (
	"fmt"
	"runtime"
)

// Build information. These variables are set at build time using ldflags.
var (
	// Version is the semantic version of the application
	Version = "dev"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// BuildDate is the date when the binary was built
	BuildDate = "unknown"

	// GoVersion is the version of Go used to compile the binary
	GoVersion = runtime.Version()

	// Platform is the operating system and architecture
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// Info returns a formatted version string with all build information
func Info() string {
	return fmt.Sprintf(`DocLoom version %s
  Build Date: %s
  Git Commit: %s
  Go Version: %s
  Platform:   %s`,
		Version,
		BuildDate,
		GitCommit,
		GoVersion,
		Platform,
	)
}

// Short returns just the version string
func Short() string {
	return Version
}
