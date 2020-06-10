

play:
	FIRESTORE_EMULATOR_HOST=localhost:8080 WEB_UI_ASSETS=./web-ui/out go run ./cmd/main.go --localhost --path example_mongoke.yml

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

bench:
	go test ./benchmarks -bench=BenchmarkQuery/main -blockprofile=block.out -memprofile memprofile.out -cpuprofile cpu.out
	go tool pprof -http=":8081" benchmarks.test cpu.out
