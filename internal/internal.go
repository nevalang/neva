// Package internal only needed to embed runtime source code.
package internal

import "embed"

//go:embed runtime
var RuntimeFiles embed.FS
