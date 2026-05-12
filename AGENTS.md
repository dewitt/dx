# Agent Instruction Protocol: Working with `declare`

This document defines the behavioral constraints for all AI agents (Architects, Implementers, and Judges) contributing to this repository.

## 1. The Primacy of the Declaration
The `.dx` file is the source of truth. You must never generate imperative code that violates a defined invariant in the `.dx` file. If an invariant is technically impossible to satisfy, you must propose a mutation to the `.dx` file rather than "fixing it in code."

## 2. Explicit Assumption Logging
When implementation requires a choice not specified in the `intent` or `invariants`, you **must not** choose silently.
1. Add a new entry to the `assumptions` block in the `.dx` file.
2. Document the heuristic leap and why it was made.
3. Only proceed with implementation once the assumption is recorded.

## 3. Verification Loop
Before declaring a task "complete":
1. Execute `declare lint` on all modified `.dx` files.
2. Generate/Run the implementation.
3. Compare the implementation behavior against the `contracts` block.
4. If a contract fails, treat the failure as a semantic bug.

## 4. Pruning and Parsimony
As an Architect, your goal is the minimum viable constraint set. Avoid over-specifying. If the user intent can be achieved without a specific invariant, move that constraint to the `unconstrained` block.

## 5. Communication with Humans
When discussing changes, use `declare diff` to explain semantic shifts. Do not summarize code changes; summarize changes to the **intent** and **invariants**.
