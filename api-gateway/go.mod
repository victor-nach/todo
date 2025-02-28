module github.com/victor-nach/todo/api-gateway

go 1.23.1

replace github.com/victor-nach/todo/proto => ../proto

require github.com/victor-nach/todo/proto v0.0.0-00010101000000-000000000000

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/stretchr/testify v1.10.0
	github.com/uptrace/bunrouter v1.0.22
	go.uber.org/mock v0.5.0
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)
