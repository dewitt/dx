# Architecture: The `declare` Paradigm

`declare` is a Heuristic Intermediate Representation (HIR) and toolchain designed for the agentic AI era. It acts as the boundary layer between human intent and machine-generated implementation.

## 1. Core Philosophy

In a paradigm where LLMs write the majority of imperative code, the human role shifts from syntax authoring to system design and constraint definition. However, current "vibe coding" suffers from silent AI hallucinations and semantic drift. 

`declare` solves this by forcing a formal, version-controllable design phase via `.dx` files. 
- **Decoupling:** The `.dx` file defines the *what* and the *constraints*. The generated imperative code defines the *how*.
- **Orthogonality:** The specification must never leak imperative logic or dictate internal architecture (e.g., unit tests). It focuses purely on observable system state and Service Level Objectives (SLOs).
- **Formalized Hallucination:** Agents are required to explicitly declare heuristic leaps in an `assumptions` block, turning silent hallucinations into auditable, promotable workflow states.

## 2. The `.dx` Artifact

`.dx` files are written in a strictly constrained subset of YAML. This provides syntax highlighting out-of-the-box, ergonomic multi-line string support for humans, and a deterministic Abstract Syntax Tree (AST) for machine parsing.

### Core Blocks
*   **`intent`**: High-level semantic goals and business logic.
*   **`invariants`**: Non-negotiable physical, systemic, or performance constraints.
*   **`assumptions`**: Heuristics the agent filled in, waiting for human promotion or rejection.
*   **`contracts`**: Black-box verification rules (state transitions, standard I/O, exit codes) used to prove the implementation satisfies the invariants.
*   **`unconstrained`**: Explicitly declared degrees of freedom, preventing over-specification and guiding agent restraint.

## 3. The Multi-Agent Loop

The `declare` ecosystem is designed to be operated by a swarm of specialized agents, coordinated via the `.dx` file.

1.  **The Archaeologist (Extraction):** Distills legacy imperative code into semantic intent and observable invariants, outputting a base `.dx` file.
2.  **The Architect (Refinement):** Modernizes the `.dx` file, applies system-wide constraints, and flags semantic gaps as `assumptions`.
3.  **The Implementer (Coding):** Reads the `.dx` file (often via an AST-compiled JSON export) and generates the imperative code. It operates strictly within the defined invariants.
4.  **The Judge (Verification):** Executes the implementation against the `contracts` block using black-box testing. It cross-references failures with the `.dx` file and issues specific correction prompts to the Implementer.

## 4. The CLI Toolchain

The `declare` binary contains **no LLM**. It is a blindingly fast, deterministic toolchain built to enforce the `.dx` specification.

*   `declare lint`: Enforces the strict YAML subset and structural rules.
*   `declare fmt`: Formats the `.dx` file to a canonical representation.
*   `declare diff`: Parses two `.dx` files and outputs a semantic ledger of operations (e.g., `[PROMOTED] assumptions.x -> invariants.y`), rather than a noisy text diff.
*   `declare export`: Compiles the `.dx` file into a tightly packed, token-optimized format for ingestion by agent context windows, stripping human comments and standardizing structure.

## 5. Security & Safety

Agents modify `.dx` files directly, but rely on `declare lint` in their event loop to catch structural entropy. Semantic conflicts are resolved through the `declare diff` tool and explicit `reason` fields within the `.dx` invariants, creating an auditable ledger of architectural decisions.
