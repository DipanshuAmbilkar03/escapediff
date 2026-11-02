# Roadmap

## Next
- A Markdown output for pull request comments
- Read GitHub Actions logs with timestamps in `compare`, after testing on a real Actions run
- Cover test files with `go test -gcflags=-m=2 -run '^$'` as an option

## Later
- Bounds checks that are no longer eliminated, from `-d=ssa/check_bce`
- Group a new escape with the parameter leak or lost inlining that caused it
- Match renamed functions by their body position, so a rename isn't shown as gone and new
