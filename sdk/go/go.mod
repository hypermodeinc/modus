module github.com/hypermodeinc/modus/sdk/go

go 1.23.0

require (
	github.com/hypermodeinc/modus/lib/manifest v0.0.0-20241015202547-171c88427e0a
	github.com/hypermodeinc/modus/lib/wasmextractor v0.0.0-20241015202547-171c88427e0a
)

// NOTE: during developement, you can use replace directives if needed
// to pick up changes to the above dependencies.  However, they MUST not be committed
// to this file, because that will break `go install` of the build tools.

// replace github.com/hypermodeinc/modus/lib/manifest => ../../lib/manifest
// replace github.com/hypermodeinc/modus/lib/wasmextractor => ../../lib/wasmextractor

require (
	github.com/fatih/color v1.17.0
	github.com/hashicorp/go-version v1.7.0
	github.com/mattn/go-isatty v0.0.20
	github.com/rs/xid v1.6.0
	github.com/tidwall/sjson v1.2.5
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	golang.org/x/mod v0.21.0
	golang.org/x/tools v0.26.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/jsonc v0.3.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
