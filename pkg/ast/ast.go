// Package ast defines the in-memory representation of a `.dx` declaration.
//
// The AST mirrors the schema described in SPEC.md (v0.1.0). It is intentionally
// shallow: the `.dx` file is the source of truth, and the AST is a transparent
// projection of it. We retain the original *yaml.Node graph alongside the
// decoded values so that downstream tooling (lint, fmt, diff) can inspect
// physical YAML features that strict-mode parsing would otherwise discard --
// most notably literal vs. folded block scalars (SPEC §4.2), node positions for
// diagnostics, and head/foot comments for round-trip formatting.
package ast

import "gopkg.in/yaml.v3"

// Declaration is the root of a parsed `.dx` file.
//
// Field ordering follows SPEC §4.2 ("Root Key Ordering") so that re-emitted
// canonical forms preserve the recommended agent ergonomics.
type Declaration struct {
	System        string              `yaml:"system"`
	Intent        Intent              `yaml:"intent"`
	Invariants    map[string]string   `yaml:"invariants"`
	Assumptions   map[string]string   `yaml:"assumptions"`
	Contracts     map[string]Contract `yaml:"contracts,omitempty"`
	Unconstrained map[string]string   `yaml:"unconstrained,omitempty"`

	// Node is the raw YAML document node retained after decoding. It is
	// populated by the loader (see pkg/lint) and is required for physical
	// checks the SPEC mandates -- e.g., refusing folded scalars (SPEC §4.2)
	// or anchors/aliases.
	//
	// ASSUMPTION (ast.node_retention): SPEC §4.2 says we must reject folded
	// scalars and anchors. The decoded Go values discard that information,
	// so we keep the node graph here to enable physical inspection. This
	// is a structural choice not mandated by SPEC; documented per
	// AGENTS.md §2.
	Node *yaml.Node `yaml:"-"`
}

// Intent expresses the high-level semantic purpose of the declaration
// (SPEC §4.3 / `intent`).
type Intent struct {
	// Primary is the core objective. Required.
	Primary string `yaml:"primary"`

	// Secondary is an optional set of supporting objectives or
	// non-functional goals.
	//
	// ASSUMPTION (ast.intent_secondary_shape): SPEC §4.3 describes
	// `secondary` as "Supporting objectives or non-functional goals"
	// without pinning the shape. We model it as a list of strings,
	// which is the most natural fit for a multi-item enumeration of
	// goals and matches typical agent emission patterns. If a single
	// string is provided, the loader can normalize it to a one-element
	// slice. Documented per AGENTS.md §2.
	Secondary []string `yaml:"secondary,omitempty"`
}

// Contract is a single black-box verification rule (SPEC §4.3 / `contracts`).
//
// ASSUMPTION (ast.contract_field_types): SPEC §4.3 specifies the three
// fields (`given`, `when`, `then`) as state/triggers/outcomes without
// constraining their YAML shape. We model them as free-form strings to
// preserve human-authored prose (typically literal block scalars). A
// future revision may introduce a structured form; the linter currently
// accepts only scalar values for these fields. Documented per AGENTS.md §2.
type Contract struct {
	Given string `yaml:"given"`
	When  string `yaml:"when"`
	Then  string `yaml:"then"`
}
