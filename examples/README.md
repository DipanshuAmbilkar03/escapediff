# Examples

- `shop/` is a made-up order service in three versions: `before`, `after` with typical pull request changes, and `comments`, which only moves lines. Each version is its own module, so the paths in their logs line up.
- `logs/` holds real compiler output from the shop versions and from five commits in open-source projects. See `logs/README.md`.

From the repository root:

    go run ./plugin/scripts compare examples/logs/shop-before.log examples/logs/shop-after.log
