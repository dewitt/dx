# Specification: The `.dx` Language (v0.1.0)

## 1. Physical Format
Files must be valid YAML 1.2. The canonical file extension is `.dx`.

## 2. Structural Constraints
To maintain a deterministic Abstract Syntax Tree (AST) and prevent semantic drift during AI processing, the following restrictions apply:
*   **No Anchors/Aliases:** The use of `&` (anchors) and `*` (aliases) is strictly forbidden.
*   **No Complex Tags:** Custom YAML tags (e.g., `!!binary`, `!!set`) are not supported.
*   **Literal Scalars Only:** All multiline strings must use the literal block scalar (`|`). The folded scalar (`>`) is prohibited due to ambiguous whitespace handling in diverse LLM tokenizers.
*   **Root Key Ordering:** While YAML is unordered, agents should prefer the order: `system`, `intent`, `invariants`, `assumptions`, `contracts`, `unconstrained`.

## 3. Schema Definitions

### `system` (Required)
A unique identifier for the declaration.
- Type: String (Slug format)

### `intent` (Required)
The high-level semantic purpose of the implementation.
- `primary`: The core objective.
- `secondary`: (Optional) Supporting objectives or non-functional goals.

### `invariants` (Required)
Non-negotiable constraints that the implementation must satisfy. 
- Map of `id: string`.
- Keys should be prefixed by category (e.g., `sec_`, `perf_`, `iface_`).

### `assumptions` (Required)
Heuristics or design choices made by the agent that require human validation.
- Map of `id: string`.
- Empty maps are allowed but the key must exist to signal a "zero-assumption" state.

### `contracts` (Optional)
Verifiable state-transition rules for black-box testing.
- Map of named contract blocks.
- Fields: `given` (initial state), `when` (execution triggers), `then` (expected outcome/side-effect).

### `unconstrained` (Optional)
Explicitly declared degrees of freedom.
- Map of `category: description`.
