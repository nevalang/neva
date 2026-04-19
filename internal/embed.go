//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package internal

import "embed"

// Efs embeds the runtime files.
//go:embed runtime
var Efs embed.FS
