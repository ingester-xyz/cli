module github.com/ingester-xyz/cli

go 1.20

replace github.com/ingester-xyz/cli/cmd => /Users/albererre/Github/ingester-xyz/cli/cmd

replace github.com/ingester-xyz/cli/pkg/s3 => /Users/albererre/Github/ingester-xyz/cli/pkg

require (
	github.com/aws/aws-sdk-go v1.55.7
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)
