# 0002: Match decisions by function and text, not by line

**Decision:** each diagnostic is keyed by kind, file, enclosing function and message text. Line numbers only break ties between identical entries, using the offset from the function's first line.

**Why:** almost every pull request moves lines, so matching by `file:line` reports most of a file as changed. Function and text are stable across unrelated edits. The comments-only shop example and chi `cf537d4` check that moved lines stay quiet.

**Cost:** a renamed function shows its escapes as gone and new. Two identical escapes in one function that swap places aren't reported, which is the right answer.

**Where the function comes from:** the `can inline` and `cannot inline` lines that `-m=2` prints on each function's first line, rather than parsing source. This keeps saved logs self-contained, so `compare` works without either revision's source.
