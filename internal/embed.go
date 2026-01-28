package internal

import "embed"

// Efs embeds the runtime files.
//go:embed runtime
var Efs embed.FS
