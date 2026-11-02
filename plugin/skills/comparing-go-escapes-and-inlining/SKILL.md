---
name: comparing-go-escapes-and-inlining
description: Compares the Go compiler's escape analysis and inlining decisions between a base branch and the current change, then explains each new heap escape, lost inlining, new parameter leak or lost devirtualization in the code. Use when the user asks whether a Go change adds allocations or heap escapes, why a benchmark got slower after a change, whether a function is still inlined, wants a performance review of a Go pull request, or shares `go build -gcflags=-m` output from two builds.
---

# Comparing Go escapes and inlining

`go build -gcflags=-m=2` prints every escape and inlining decision, which is thousands of lines for a real package. What matters in a review is what changed. The bundled tool builds two revisions and lists only the differences. It needs Go and git, and it compiles code but never runs it.

## Run the tool

From inside the Go module, compare the working tree with the base the change started from:

    go run "${CLAUDE_PLUGIN_ROOT}/scripts/escapediff.go" git --save /tmp/escapediff "$(git merge-base origin/main HEAD)" ./...

- The working tree includes uncommitted edits. Add `--head REF` to build a commit instead.
- Narrow `./...` to the packages the change touches when the module is large.
- `--save DIR` keeps both raw logs. Use them for the explanations below.
- If the user already has two saved `-gcflags=-m=2` outputs, run `compare BASE.log HEAD.log` instead. Both must come from the same Go version.
- If a build fails, report the error. Don't guess at the result.

## Read each finding in the code

The report lists new heap escapes, functions no longer inlined (with the compiler's reason and cost), new parameter leaks, calls no longer devirtualized, functions close to the inline budget, and what improved. For each regression:

1. Open the head position and say which value escapes and why. To see the compiler's reasoning, find the same `file:line:col` in `/tmp/escapediff/head.log`: the indented `flow:` lines under it show the path to the heap.
2. A new parameter leak or lost inlining in one function often explains heap escapes in its callers. Connect them, and start with the cause.
3. Judge whether the code is hot. A new escape in a panic message, an error path or startup code usually doesn't matter. One in a request handler loop or an encoder probably does. Say which, and why.
4. For a function no longer inlined, give its cost against the budget and what in the diff added the cost.

## Report

1. List regressions on hot paths first, with `file:line` and the cause in one sentence each.
2. Propose the smallest fix in the project's style, such as passing a value instead of a pointer, keeping a fast path small enough to inline, or moving a cold branch into its own function.
3. Say that these are compile-time decisions. Suggest `go test -bench . -benchmem -count 10` on both revisions with `benchstat` to measure the real effect before anyone acts on it.

## Don't

- Don't edit code unless the user asks.
- Don't claim a slowdown or speedup from this report alone. It shows what the compiler decided, not how often the code runs.
- Don't suggest `unsafe` or compiler flags only to make the report quiet.
- Don't treat "No changes" as proof that performance is unchanged. Runtime behavior, such as a new loop or a lock, isn't in this output.
