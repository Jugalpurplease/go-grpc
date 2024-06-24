# Define the proto source files
PROTO_FILES := *.proto

# Directory for generated Go files
PROTOC_OUT_DIR := pb

OPEN_API := openapi

# Commands
PROTOC := protoc
RM := rm -rf
GO := go

# Targets
generate:
	$(PROTOC) --go_out=$(PROTOC_OUT_DIR) --go_opt=paths=source_relative \
	          --go-grpc_out=$(PROTOC_OUT_DIR) --go-grpc_opt=paths=source_relative \
	          $(PROTO_FILES)

clean:
	$(RM) $(PROTOC_OUT_DIR)/*.go

run:
	$(GO) run notes_server/main.go

proto:
	$(PROTOC) --go_out=$(PROTOC_OUT_DIR) --go_opt=paths=source_relative \
	          --go-grpc_out=$(PROTOC_OUT_DIR) --go-grpc_opt=paths=source_relative \
	          notes.proto

gw:
	$(PROTOC)  --go_out=$(PROTOC_OUT_DIR) --go_opt=paths=source_relative \
	          --grpc-gateway_out=$(PROTOC_OUT_DIR) --grpc-gateway_opt=paths=source_relative \
			   --go-grpc_out=$(PROTOC_OUT_DIR) --go-grpc_opt=paths=source_relative \
	          --openapiv2_out=${OPEN_API} --openapiv2_opt=logtostderr=true \
	          $(PROTO_FILES)
