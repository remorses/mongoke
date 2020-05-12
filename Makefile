

play:
	go run ./cmd/main.go --www ./web-ui/out --localhost --path example_mongoke.yml

.PHONY: build
build:
	cd web-ui && yarn build 

test:
	go test ./... -cover

test2:
	gotestsum -f dots-v2
