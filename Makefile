

play:
	go run ./cmd/main.go --path example_mongoke.yml

.PHONY: build
build:
	cd web-ui && yarn build 