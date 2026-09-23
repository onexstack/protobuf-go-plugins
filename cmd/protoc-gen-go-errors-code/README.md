# protoc-gen-go-errors-code

`protoc-gen-go-errors-code` turns an enum's error annotations into a Markdown
table, so the error-code list in a docs directory is derived from the contract
rather than maintained beside it.

Annotate an enum with `(errors.default_code)` and its values with `(errors.code)`:

```proto
enum DemoError {
  option (errors.default_code) = 500;
  DEMO_UNKNOWN   = 0 [(errors.code) = 501];
  DEMO_NOT_FOUND = 1 [(errors.code) = 404];
}
```

```sh
protoc --go-errors-code_out=paths=source_relative:./docs demo.proto
```

writes `demo_code.md`:

| Reason | HTTP Status Code | Description |
| :----: | :--------------: | :---------- |
| DemoUnknown | 501 |  |
| DemoNotFound | 404 |  |

An enum carrying no `errors.code` anywhere is skipped, and a file whose enums
are all skipped produces no output.

The directory name is load-bearing: this plugin is invoked as
`--go-errors-code_out`. It was previously named `protoc-gen-go-errordoc` — the
scaffold's name, not the plugin's — while its own `main.go` and template both
said `go-errors-code`. The directory was the odd one out, so it moved.

`errors/` is the extension definition (`errors.default_code`, `errors.code`) the
plugin reads, derived from the kratos errors proto. It is a local copy rather
than an import of `github.com/go-kratos/kratos/v2/errors` so that the extension
numbers this plugin looks for are pinned here; the upstream `go_package` still
named kratos and was corrected to this directory, so regenerating no longer
writes the descriptor back to a path this module does not own.

## Acknowledgements

Originally from onexstack/onex's `tools/` tree, under the MIT license. The
extension definitions in `errors/` derive from go-kratos/kratos.
