# Patched copy of docker-credential-acr-env

This directory is a copy of
[github.com/chrismellard/docker-credential-acr-env](https://github.com/chrismellard/docker-credential-acr-env)
at upstream commit `82a0ddb2758901b711d9d1614755b77e401598a1` (2023-03-04),
carried in-tree because upstream is unmaintained. kaniko uses it through a
`replace` directive in the top-level `go.mod`, which covers both the library
used by `pkg/creds` and the `docker-credential-acr-env` binary that
`deploy/Dockerfile` builds and ships in the images.

Only `LICENSE`, `README.md`, `go.mod`, `go.sum`, `main.go`, `cmd/` and `pkg/`
were copied. The licence (Apache-2.0) and upstream copyright headers are kept.

## Local changes

1. GO-2026-6225: `pkg/credhelper/helper.go` replaces the unanchored host regex
   with an anchored, label-aware check (`isACRRegistry`): exact `host[:port]`
   only, valid port range, no trailing dot, no malformed labels, still
   case-sensitive. Non-ACR hosts now get a "credentials not found" error from
   `Get` rather than a generic error. Tests are in
   `pkg/credhelper/helper_test.go`.

2. `go.mod` / `go.sum`: requirements raised to the versions kaniko's own
   build list already selects (for example `golang.org/x/crypto` from a 2022
   pseudo-version to the version kaniko uses), and the `go` directive to match. This does not
   change what kaniko builds (Go's minimum version selection already picked
   these versions); it stops scanners and Dependabot reporting the ancient
   upstream pins in this `go.mod` as findings. Keep them in step with the
   top-level `go.mod` when bumping dependencies.

Every local code change is marked with a `Chainguard patch` comment. To review
them, diff this directory against the upstream commit above.

Run the tests from this directory (it is its own Go module and is not part of
the top-level `go test ./...`):

    go test ./...
