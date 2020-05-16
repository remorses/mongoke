

play:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go run ./cmd/main.go --www ./web-ui/out --localhost --path example_mongoke.yml

.PHONY: build
build:
	cd web-ui && yarn build 

test:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./... -cover -failfast

test2:
	FIRESTORE_EMULATOR_HOST=localhost:8080 gotestsum -f dots-v2
