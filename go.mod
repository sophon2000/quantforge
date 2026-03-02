module github.com/sophon2000/quantforge

go 1.25.0

require (
	github.com/scmhub/ibapi v0.10.44
	github.com/scmhub/ibsync v0.10.44
	github.com/sdcoffey/big v0.7.0
	github.com/sdcoffey/techan v0.0.0-20211117160920-e192d24cb693
	github.com/spf13/cobra v1.10.2
)

// 使用 fork 的代码，因 fork 的 go.mod 仍声明为 sdcoffey/techan，需用 replace 指向 fork 的 dev 分支
replace github.com/sdcoffey/techan => github.com/sophon2000/techan v0.0.0-20211117160920-e192d24cb693

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robaho/fixed v0.0.0-20251201003256-beee5759f86a // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.41.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
