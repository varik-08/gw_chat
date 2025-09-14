PROTO_DIR := proto
GEN_DIR := internal/grpc/gen

PROTOC := protoc

.PHONY: proto
proto:
	@mkdir -p $(GEN_DIR)
	rm -rf $(GEN_DIR)/*
	$(PROTOC) \
		-I $(PROTO_DIR) \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO_DIR)/*.proto

.PHONY: tidy
tidy:
	go mod tidy

