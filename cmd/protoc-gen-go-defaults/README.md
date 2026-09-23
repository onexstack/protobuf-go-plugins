# protoc-gen-go-defaults

`protoc-gen-go-defaults` compiles a field's `(defaults.value)` annotation into a
`Default()` method on its message. `internal/pkg/reqval` in onex calls that
method once per request, through an `interface{ Default() }` assertion.

It is a fork of [linka-cloud/protoc-gen-defaults](https://github.com/linka-cloud/protoc-gen-defaults),
kept in this repository so the generator can be corrected here rather than
pinned to an upstream release. The change made so far is under
[Difference from upstream](#difference-from-upstream).

The directory name is load-bearing: protoc invokes a plugin named
`protoc-gen-go-defaults` as `--go-defaults_out`.

## Two halves

This plugin is unusual in this repository in that it ships a **runtime** package
alongside the binary:

| Path | Role |
|---|---|
| `cmd/protoc-gen-go-defaults/` | the generator, run by protoc |
| `defaults/` (repository root) | the runtime half: the `(defaults.value)` extension definitions |

The generated `.pb.go` files carry a blank import of `defaults` — it comes from
`go_package` in `defaults/defaults.proto` — which is what registers the
extension with the protobuf runtime. That is why the package does not live under
`cmd/` with the binary: an import path containing `cmd/` would misdescribe what
it is, and it is application code that imports it, not just the generator.

It also means this module is a **runtime dependency of its consumers**, not only
a build-time tool. onex requires it, and `defaults` appears in the import graph
of every generated package.

## Difference from upstream

Upstream kept one `imports` set on the module and had the template read it:

```go
imports: func() string { for v := range m.imports { ... } }
```

pgs registers a template when a file is walked but renders it only after every
target file has been walked, so that function reads the set at **render** time.
Each generated file therefore received the union of every file's imports, and a
file that does not reference a package still imported it — which does not
compile.

The trigger is easy to miss: a proto that mentions a well-known type another
file in the same protoc invocation uses — a `google.protobuf.DescriptorProto`
field, say — exports that import to its siblings. It surfaced here only because
`tests/pb` contains exactly that pair, and because pgs walks them in name order
(`test.proto` before `types.proto`).

This fork keys the set by file path, so a generated file imports what it
references and nothing else. Keying by `FullyQualifiedName` is not enough: two
files in the same proto package share one.

## Generated

```proto
message Order {
    int64 page = 3 [(defaults.value).int64 = 1];
}
```

```go
func (x *Order) Default() {
	if x.Page == 0 {
		x.Page = 1
	}
}
```

## Tests

`tests/` exercises the generated methods against the fixtures in `tests/pb`; it
is the plugin's only behavioural coverage, and it is what caught the import
leak above. Regenerate a fixture with:

```sh
protoc \
  --proto_path=. \
  --proto_path=$(dirname $(dirname $(command -v protoc)))/include \
  --go-defaults_out=paths=source_relative:. \
  tests/pb/types.proto tests/pb/test.proto
```

run from this directory, with `defaults/` reachable on a second `--proto_path`
(the repository root).

## Acknowledgements

protoc-gen-defaults was written by Linka Cloud under the Apache License 2.0.
See [LICENCE](./LICENCE).
