# Equity perpetual dividend funding reproduction

A retrospective, five-event accounting reproduction of Seo et al., [arXiv:2610.22839v1](https://arxiv.org/abs/2610.22839). Start with [REPORT.md](REPORT.md). The prior inspection, numerical result, provenance limits and post-seal documentary checks are disclosed there and in [FOLLOWUP.md](FOLLOWUP.md).

This repository contains the study code, frozen records, sealed run and handoff. It does **not** distribute the paper, exchange responses or other source bytes, which are held in ignored `data/`. The study code depends on a separate research-platform checkout. No permissions to redistribute third-party source documents are asserted here; public accessibility and noncommercial use alone do not establish redistribution rights.

## Checking the sealed run

Exact verification requires Go 1.27.2, a compatible research-platform checkout, and the six exact source files listed in [manifest.json](records/equity-perpetual-dividend-reproduction/manifest.json). The sealed run used platform repository commit `9fb4290b6c50558401275ea06b820f663f549724` and study implementation commit `7ddd773fc64e7ee3804656f6c60b37f897acaed6`; the sealed plan records their content identities. This public snapshot uses a relative platform path in `go.mod`. The sealed plan records the original implementation identity; this snapshot has a separate Git history. The research platform is not included here. Create a temporary Go workspace referencing both checkouts instead. For example, from the study root, with `STUDY_DIR` and `PLATFORM_DIR` set to absolute checkout paths:

```sh
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
workdir=$(mktemp -d)
(cd "$workdir" && GOWORK=off go work init "$STUDY_DIR" "$PLATFORM_DIR")
```

Run `go work init` **inside** the temporary directory, then set `GOWORK` to its `go.work` file before running the commands below. Alternatively adapt the example to your shell. Both repositories must have the same modules as the study and platform checkouts. The exact local paths in source manifests are relative to the study root.

```sh
export GOWORK="$workdir/go.work"
go test ./...
go run ./cmd/equity-perpetual-dividend-reproduction verify 20261111T125325Z-0a66a12c1a04
python3 scripts/result-table.py
```

The manifest gives the original source URLs, byte sizes and SHA-256 digests. The five API URLs are public queries, and the paper locator specifies v1; they can be retrieved into the listed paths **only in a separate copy** and checked against the manifest before verification. A later response, changed formatting, or revised web content may yield a different hash even if the observations look the same. The paper HTML can also change without a version change. A different hash cannot verify this sealed run: obtain the exact original bytes through an authorized channel, or perform a separately disclosed reproduction with new source identities. Do not overwrite the original source files or sealed artifacts.

The sealed study's semantic check reuses its evaluator. Its PASS is a replay result, distinct from platform integrity and from independent scientific validation. `FOLLOWUP.md` and `followup-sources.json` describe additional, post-seal documentary sources; they are not used by the sealed run.
