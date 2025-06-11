module text-generation-example

go 1.23.1

toolchain go1.23.10

require (
	github.com/hypermodeinc/modus/sdk/go v0.17.0
	github.com/tidwall/gjson v1.18.0
)

require (
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)

replace github.com/hypermodeinc/modus/sdk/go => ../..
