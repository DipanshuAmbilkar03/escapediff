# escapediff

Claude Code skill and Go CLI that compares the Go compiler's escape analysis and inlining decisions between two builds, such as a base branch and a pull request.

## Architecture
`plugin/scripts/escapediff.go` holds all the logic in one standard-library file, so `go run` needs no setup. `Parse` reads `-gcflags=-m=2` (or `-m`) output into diagnostics and per-function inlining decisions, `attribute` gives each diagnostic its enclosing function from the `can inline` and `cannot inline` lines, `Compare` matches the two builds, and `gitBuilds` unpacks a revision with `git archive` and builds it. `plugin/skills/comparing-go-escapes-and-inlining/SKILL.md` tells Claude how to use it. See `docs/architecture.md`.

## Commands
- `make test`: gofmt, vet, unit tests and integration tests in `tests/`
- `make lint`: plugin and marketplace validation
- `make demo`: compare the shop example logs
- `make logs`: rebuild the shop logs from `examples/shop`

## Conventions
- Standard library only. Never use the network, and never run the code being built. Building with `go build` is the only thing the tool does with a project.
- Never change the user's checkout. Revisions go into a temporary directory through `git archive`.
- Every new report section needs real compiler output in `examples/logs/` that it fires on, and must stay quiet on `shop-comments.log` and `real/chi-cf537d4-*`.
- Logs in `examples/logs/` are real compiler output. `examples/logs/README.md` says which Go version and command made each one.

## Do and don't
- Do call findings compile-time decisions. They aren't measured allocations.
- Don't say a change has no performance effect because the report is empty.

## Definition of done
`make test` and `make lint` pass, and `ROADMAP.md`, the README and docs are updated when relevant.
