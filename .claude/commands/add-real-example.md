Add a real-project example for: $ARGUMENTS

1. Clone the open-source project, check its license allows including compiler output, and pick a commit that changes escape analysis or inlining.
2. From the clone, run `go run <this repo>/plugin/scripts/escapediff.go git --head C --save DIR C^`, and copy `DIR/base.log` and `DIR/head.log` to `examples/logs/real/<project>-<short commit>-base.log` and `-head.log`.
3. Read the commit's diff and check that every finding matches a real change. Note the commit, Go version and license in `examples/logs/README.md`.
4. Add a test in `plugin/scripts/escapediff_test.go` that pins the result, then run `make test`.
