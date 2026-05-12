// Package export emits a parsed `.dx` declaration in agent-optimized formats.
//
// This is a stub. The first real implementation will provide a tightly packed
// JSON projection of the AST (see ARCHITECTURE.md §4) with comments stripped
// and key ordering canonicalized for stable LLM context-window ingestion.
package export

import (
	"errors"
	"io"

	"github.com/dewitt/declare/pkg/ast"
)

// Format names a target serialization for `declare export`.
type Format string

const (
	// FormatJSON emits a compact JSON projection of the AST.
	FormatJSON Format = "json"
)

// ErrNotImplemented signals that the export pipeline has not yet been wired
// up. The CLI surfaces this verbatim so users see a clear stub indicator
// rather than a silent no-op.
var ErrNotImplemented = errors.New("export: not yet implemented")

// Write serializes d in the requested format. It currently returns
// ErrNotImplemented for every format; callers should treat any nil error
// as a future contract.
func Write(_ io.Writer, _ *ast.Declaration, _ Format) error {
	return ErrNotImplemented
}
