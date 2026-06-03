package golang

// ReleaseBuildArgs returns standardized `go build` args for end-user artifacts.
// Flags are intentionally centralized to avoid drift across CLI/backends and
// to keep binaries compact and reproducible by default.
func ReleaseBuildArgs(outputPath, target string) []string {
	return []string{
		"build",
		"-trimpath",       // remove host-specific paths from debug metadata
		"-buildvcs=false", // do not embed VCS metadata into produced artifacts
		"-ldflags",
		"-s -w", // strip debug/symbol data to reduce binary size
		"-o",
		outputPath,
		target,
	}
}
