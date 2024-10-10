build:
	protoc --go_out=./user --go-grpc_out=./user ./user/user.proto
	protoc -I=./guild -I=/home/kamisama/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf/ --go_out=./guild --go-grpc_out=./guild --go_opt=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types --go-grpc_opt=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types ./guild/guild.proto