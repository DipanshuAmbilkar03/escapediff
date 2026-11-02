# Changelog

All notable changes are documented here, in [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format. This project uses [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.0] - 2026-11-02
### Added
- `compare` command: compares two saved `go build -gcflags=-m=2` outputs
- `git` command: builds a base revision and the working tree, or `--head REF`, and compares them without changing the checkout
- Reports new heap escapes, functions no longer inlined with the compiler's reason and cost, inlined callees whose cost grew, new parameter leaks, calls no longer devirtualized, and functions close to the inline budget, plus what improved
- Matches by file, enclosing function and message text, so moved lines don't count as changes
- Hides changes in files outside the built packages, such as generic standard library code, unless `--all` is set
- `--json`, `--fail` and `--save` options
- Claude Code plugin with the `comparing-go-escapes-and-inlining` skill
