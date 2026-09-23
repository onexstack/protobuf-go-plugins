module github.com/onexstack/protobuf-go-plugins

// 1.24, not the toolchain-specific 1.26.3 this file was seeded with: a patch
// version in the go directive pins the build to exactly that toolchain, so the
// module stops building for anyone on 1.25 or an older 1.26. Nothing here uses
// a post-1.24 language feature.
go 1.25.0

// Pinned to the version onex compiles its generated code against, so the plugin
// and its consumers cannot disagree about protogen.
require google.golang.org/protobuf v1.36.12

require (
	github.com/iancoleman/strcase v0.3.0
	github.com/prometheus/common v0.71.0
	golang.org/x/text v0.41.0
)

require (
	github.com/spf13/afero v1.3.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
)

require (
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/stretchr/testify v1.12.1
)
