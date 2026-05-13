# Architecture: The `declare` Paradigm

`declare` is a declarative specification language and toolchain
designed for the agentic AI era. It sits at the boundary between
human intent and machine-generated implementation, holding the
*idea* of the system in a form that humans can review, version, and
argue about — and that AI agents can consume, validate, and produce
imperative code from.

## 1. Philosophy

A program is two things at once: an artifact (the source code, in some particular language) and an idea (what the system is supposed to *be*). For most of computing history the two have been fused. The idea exists only as it is encoded in the artifact; reading the idea means reading the code. When the artifact changes, the idea changes — and there is no separate place where the idea lives that can be checked, versioned, reviewed, or argued about.

The intellectual position behind `declare` is that this fusion is now optional. In a world where AI writes the imperative artifact, the *idea* of the system is the load-bearing thing humans should attend to, and the artifact is a derived witness — one of many possible implementations, all equally valid if they satisfy the idea.

A `.dx` file is meant to be that idea, written down. It does not specify *how* the system computes anything; it specifies *what is true* about the system's observable behavior. The implementation is a witness that those truths hold. Two witnesses, in different languages, both honest, are equivalent.

This is an old idea given new traction. **Formal-methods languages** (Z, VDM, TLA+, B, Alloy) have made this same separation since the 1970s — designed for human formal proof, never widely adopted because authoring them is expensive and verifying them requires expert reviewers. **Knowledge-representation languages** (KIF, OWL, Cyc) attempted a similar move for AI consumers — they failed for the inverse reason: humans couldn't agree on shared semantics, and pre-LLM machines couldn't reason fast enough to make the cost worthwhile. **Design-by-contract** (Eiffel, JML, Code Contracts) embedded the idea inside imperative languages as annotations — it gained adherents but never became the primary artifact. **Denotational semantics** (Strachey, Scott) provided the philosophical foundation — meaning before mechanism — but kept the meaning in the heads of theorists rather than version-control systems.

Each of these traditions had the right insight and the wrong era. They lacked a *consumer* that could draft, validate, and act on a declarative specification at scale. The arrival of LLM-mediated coding workflows changes both sides of that equation: drafting becomes a conversation with an agent, consuming becomes a tool-use loop, and the cost calculus inverts.

`declare` is the minimum-viable declarative intermediate language for this new operating context. It is not a successor to TLA+ or Z or OWL; it inherits from all of them. What it adds is a *delivery mechanism* — the agent that consumes the spec and produces (or modifies, or verifies) the implementation — that none of the prior traditions could assume. The spec is the source of truth; the agent is the engine; the imperative code is downstream of both.

A few honest claims this position commits us to:

- **The `.dx` file is primary.** When the spec and the code disagree, the spec wins. If the spec is impossible to satisfy, the spec changes — not the code.
- **Implementations are plural and replaceable.** A `.dx` should never name a language, library, or framework. If two implementations in two languages both satisfy every contract, both are correct, and the question of which is "real" is malformed.
- **Heuristic leaps are first-class artifacts.** When an agent must guess, the guess is recorded in `assumptions:` *before* the code is written. Silent invention is the failure mode the entire system is designed to prevent.
- **Verification is observational.** A contract describes what an outside observer would see, not what the code does internally. This is the discipline that lets implementations diverge stylistically while remaining equivalent semantically.

What we are *not* claiming:

- That `.dx` files are formal proofs. They are not — `declare` does not currently include a model checker or theorem prover. Verification today is the [`judge`](skills/judge/SKILL.md) skill walking each contract; future revisions may add mechanical verification for the subset of contracts that admit it.
- That the spec language is final. v0.1.0 is deliberately minimal; whole categories of expression (temporal logic, refinement relations, cross-system contracts) are absent and will land or be rejected as the workflow's needs become clearer.
- That AI-mediated specification will replace human design. The opposite: it makes human design *more* central by giving humans a place to do design that isn't immediately swamped by syntax decisions.

This positioning is the project's intellectual moat. The toolchain and skills below operationalize it; the philosophy is what makes the operational choices coherent.

## 2. Operating Principles

The philosophy in §1 is realized through three operating principles that govern every concrete choice in the toolchain, the skills, and the workflow:

- **Decoupling.** The `.dx` file defines the *what* and the *constraints*. The generated imperative code defines the *how*. Neither side is permitted to invade the other.
- **Orthogonality.** The specification must never leak imperative logic or dictate internal architecture (e.g., unit tests, library choice, threading model). It focuses purely on observable system state and the constraints that bind it.
- **Formalized hallucination.** Agents are required to explicitly declare heuristic leaps in an `assumptions:` block, turning silent hallucinations into auditable, promotable workflow state.

## 3. The `.dx` Artifact

`.dx` files are written in a strictly constrained subset of YAML. This provides syntax highlighting out-of-the-box, ergonomic multi-line string support for humans, and a deterministic Abstract Syntax Tree (AST) for machine parsing.

### Core Blocks
*   **`intent`**: High-level semantic goals and business logic.
*   **`invariants`**: Non-negotiable physical, systemic, or performance constraints.
*   **`assumptions`**: Heuristics the agent filled in, waiting for human promotion or rejection.
*   **`contracts`**: Black-box verification rules (state transitions, standard I/O, exit codes) used to prove the implementation satisfies the invariants.
*   **`unconstrained`**: Explicitly declared degrees of freedom, preventing over-specification and guiding agent restraint.

## 4. The Multi-Agent Loop

The `declare` ecosystem is designed to be operated by a swarm of specialized agents, coordinated via the `.dx` file.

1.  **The Archaeologist (Extraction):** Distills legacy imperative code into semantic intent and observable invariants, outputting a base `.dx` file.
2.  **The Architect (Refinement):** Modernizes the `.dx` file, applies system-wide constraints, and flags semantic gaps as `assumptions`.
3.  **The Implementer (Coding):** Reads the `.dx` file (often via an AST-compiled JSON export) and generates the imperative code. It operates strictly within the defined invariants.
4.  **The Judge (Verification):** Executes the implementation against the `contracts` block using black-box testing. It cross-references failures with the `.dx` file and issues specific correction prompts to the Implementer.

## 5. The CLI Toolchain

The `declare` binary contains **no LLM**. It is a blindingly fast, deterministic toolchain built to enforce the `.dx` specification.

The v0.1.0 command set, with implementation status:

*   `declare lint` (implemented): Enforces SPEC §2 (no anchors/aliases, no folded scalars, no custom tags, scalar leaves under `invariants`/`assumptions`/`unconstrained`) and SPEC §3 (required-key presence). Walks the retained YAML node graph for the physical-rule checks; strict-decodes the AST for the structural-decode pass.
*   `declare fmt` (implemented): Reformats a `.dx` file to its canonical representation: top-level keys in SPEC §2 order, alphabetized maps, literal block scalars for multi-line strings, no trailing whitespace. Idempotent (`fmt(fmt(x)) == fmt(x)`) and AST-preserving. Prints to stdout by default; `--write` overwrites in place.
*   `declare diff` (implemented): Parses two `.dx` files and outputs a semantic ledger of operations (`[ADDED]`, `[REMOVED]`, `[MUTATED]`, `[PROMOTED]`, `[DEMOTED]`, `[RENAMED]`), rather than a noisy text diff.
*   `declare export` (implemented): Emits the AST as canonical YAML (default) or compact JSON, with comments stripped, for ingestion by agent context windows. Byte-stable for the same AST so two agents can agree on hashes.
*   `declare contracts list` (implemented): Enumerates contract identifiers in alphabetical order. Defaults to plain text (one ID per line, shell-friendly); `--verbose` adds an indented one-line preview of given/when/then; `--format=json` emits a structured object with full bodies. Used by the [`judge`](skills/judge/SKILL.md) skill to drive a deterministic walk over each contract.
*   `declare verify` (deferred to v0.2 per SPEC §4): Will run the `contracts:` block as a black-box test harness. Until it ships, contract execution is performed by an agent under the [`judge`](skills/judge/SKILL.md) skill, typically driven by `declare contracts list` to enumerate what to check.

See [`skills/declare-toolchain/SKILL.md`](skills/declare-toolchain/SKILL.md) for invocation details, exit codes, and the post-merge ritual.

## 6. Security & Safety

Agents modify `.dx` files directly, but rely on `declare lint` in their event loop to catch structural entropy. Semantic conflicts are resolved through the `declare diff` tool, which produces a deterministic ledger of how the spec evolved between two revisions.

The auditable design-decision ledger envisioned for the long term — where each invariant carries a `reason:`, an `author:`, and a `since:` field — is reserved as a future-compatibility shape in SPEC §6 but is **not** yet expressible in v0.1.0 (where leaves under `invariants:` / `assumptions:` / `unconstrained:` are scalar strings). Today, the audit trail lives in git: the commit that introduced an invariant is its provenance, and `declare diff` is the lens for reading the change.
