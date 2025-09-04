module anthropic-functions-example

go 1.24.0

toolchain go1.25.0

require github.com/hypermodeinc/modus/sdk/go v0.18.0

require (
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)

replace github.com/hypermodeinc/modus/sdk/go => ../..
