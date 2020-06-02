

play:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go run ./cmd/main.go --www ./web-ui/out --localhost --path example_mongoke.yml

.PHONY: build
build:
	cd web-ui && yarn build 

test:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./... -cover -failfast

cov:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go test ./...  -coverpkg=./... -failfast -coverprofile=coverage.out
	go tool cover -html=coverage.out

test2:
	FIRESTORE_EMULATOR_HOST=localhost:8080 gotestsum -f dots-v2
