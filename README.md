# declare

A declarative language and toolchain for the agentic AI era.

## Overview
`declare` provides a formal boundary between high-level design and imperative implementation. In a world where AI writes the code, `declare` ensures that humans maintain control over the **intent** and **constraints** of a system without being mired in the syntax of implementation.

## Key Features
*   **Intent Compression:** Distill complex systems into readable `.dx` files.
*   **Verification:** Prove implementation compliance via black-box contracts.
*   **Auditability:** Track the evolution of system design through semantic diffs.
*   **Language Agnostic:** Decouple problem description from implementation choice (Python, Rust, Go, etc.).

## Quick Start
1. **Define:** Create a `system.dx` file specifying your invariants and intent.
2. **Lint:** Run `declare lint system.dx` to ensure structural integrity.
3. **Implement:** Pass the declaration to a coding agent to generate the artifact.
4. **Verify:** Use the integrated contracts to prove the implementation satisfies the declaration.

## Project Structure
*   `/cmd`: Source for the `declare` CLI toolchain.
*   `/pkg`: Core logic for YAML parsing, AST management, and semantic diffing.
*   `/docs`: Philosophical and technical documentation.
*   `/examples`: Sample `.dx` files and their multi-language implementations.
