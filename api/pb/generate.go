package pb

//go:generate protoc --proto_path=../ --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ../chat.proto
//go:generate protoc --proto_path=../ --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ../auth.proto
