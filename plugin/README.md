# escapediff plugin

This plugin adds the `comparing-go-escapes-and-inlining` skill to Claude Code. When you ask whether a Go change adds heap allocations or loses inlining, Claude runs the bundled `scripts/escapediff.go` with `go run`. It builds the base revision and your working tree with `-gcflags=-m=2` and lists only what changed: new heap escapes, functions no longer inlined, new parameter leaks and lost devirtualization. Claude then reads the code for each one and says whether it's on a hot path. It needs Go 1.22 or newer and git. See the repository README for examples and limits.

Install it with `/plugin install escapediff --marketplace GillMacedo/escapediff`, or with `claude plugin marketplace add GillMacedo/escapediff` and then `claude plugin install escapediff@escapediff-marketplace`.
