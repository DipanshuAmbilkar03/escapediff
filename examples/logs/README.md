# Example logs

All logs are real compiler output from Go 1.27.3 on linux/arm64. None is edited.

## Shop example

Made by `make logs` from `examples/shop`, a made-up order service:

| File | Command |
|---|---|
| `shop-before.log` | `go build -gcflags=-m=2 -o /dev/null ./... ` in `examples/shop/before` |
| `shop-after.log` | the same in `examples/shop/after` |
| `shop-comments.log` | the same in `examples/shop/comments`, which only adds comments and renames a local variable |
| `shop-before-m1.log`, `shop-after-m1.log` | `-gcflags=-m` instead of `-m=2` |

## Real projects

Each pair was made in a clone of the project with `escapediff git --head C --save DIR C^`, which builds `./...` at the commit's parent and at the commit.

| Files | Project and commit | License |
|---|---|---|
| `real/zerolog-14dcf38-*` | github.com/rs/zerolog `14dcf38`, "Make AppendKey in json inlinable", package `./internal/json/` | MIT |
| `real/uuid-85223fa-*` | github.com/google/uuid `85223fa`, "Reduce custom error allocation" | BSD-3-Clause |
| `real/chi-aba8618-*` | github.com/go-chi/chi `aba8618`, "fix: panic on Mount() when a router is mounted onto itself" | MIT |
| `real/chi-be60b2e-*` | github.com/go-chi/chi `be60b2e`, "Simplify chi.walk with slices.Concat" | MIT |
| `real/chi-cf537d4-*` | github.com/go-chi/chi `cf537d4`, "Optimize throttle middleware by avoiding unnecessary timer creation". No compiler decision changed | MIT |

The logs hold file names, function names and short code expressions from these projects, as the compiler prints them.
