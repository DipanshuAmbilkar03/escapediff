# 0003: Build the base revision from git archive

**Decision:** `git` mode unpacks the base revision with `git archive` into a temporary directory and builds it there. The head is the working tree as it is, or another `git archive` with `--head`.

**Why:** a checkout, stash or worktree changes the user's repository state, and a failure halfway could leave it changed. `git archive` only reads. Building the working tree as it is covers uncommitted edits, which is what people want to check before they push.

**Cost:** submodules, `export-ignore` files and symlinks are missing from the base tree. Builds that rely on them fail with a clear error rather than giving a wrong answer.
