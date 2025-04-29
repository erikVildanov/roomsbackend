.PHONY: proto clean migrate-up migrate-down

PROTO_DIR=proto
OUT_DIR=proto/generated
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

proto:
	@mkdir -p $(OUT_DIR)
	@for file in $(PROTO_FILES); do \
		name=$$(basename $$file .proto); \
		mkdir -p $(OUT_DIR)/$$name; \
		protoc -I=$(PROTO_DIR) -I=third_party \
			--go_out=paths=source_relative:$(OUT_DIR)/$$name \
			--go-grpc_out=paths=source_relative:$(OUT_DIR)/$$name \
			--grpc-gateway_out=paths=source_relative:$(OUT_DIR)/$$name \
			--openapiv2_out=$(OUT_DIR)/$$name \
			$$file; \
	done

clean:
	rm -rf $(OUT_DIR)

migrate-up:
	migrate -path migrations/auth -database "postgres://postgres:postgres@localhost:5432/rooms_auth?sslmode=disable" up
	migrate -path migrations/user -database "postgres://postgres:postgres@localhost:5432/rooms_user?sslmode=disable" up
	migrate -path migrations/chats -database "postgres://postgres:postgres@localhost:5432/rooms_chats?sslmode=disable" up
	migrate -path migrations/participants -database "postgres://postgres:postgres@localhost:5432/rooms_participants?sslmode=disable" up
	migrate -path migrations/messages -database "postgres://postgres:postgres@localhost:5432/rooms_messages?sslmode=disable" up
	migrate -path migrations/files -database "postgres://postgres:postgres@localhost:5432/rooms_files?sslmode=disable" up
	migrate -path migrations/privacy -database "postgres://postgres:postgres@localhost:5432/rooms_privacy?sslmode=disable" up

migrate-remove:
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_auth?sslmode=disable" -path migrations/auth force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_chats?sslmode=disable" -path migrations/chats force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_files?sslmode=disable" -path migrations/files force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_messages?sslmode=disable" -path migrations/messages force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_participants?sslmode=disable" -path migrations/participants force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_privacy?sslmode=disable" -path migrations/privacy force 0
	migrate -database "postgres://postgres:postgres@localhost:5432/rooms_user?sslmode=disable" -path migrations/user force 0