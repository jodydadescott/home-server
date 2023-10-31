module github.com/jodydadescott/home-server

go 1.20

require (
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/jinzhu/copier v0.4.0
	github.com/jodydadescott/jody-go-logger v0.0.0-20231029171416-d58235a59670
	github.com/jodydadescott/unifi-go-sdk v0.0.0-20231026203353-cc40e5471ffa
	github.com/miekg/dns v1.1.56
	github.com/spf13/cobra v1.7.0
	go.uber.org/zap v1.26.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/jodydadescott/jody-go-logger => ../jody-go-logger

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
)
