


.PHONY: clean
clean:
	rm -rf generated

.PHONY: clean
image:
	docker build . -t mongoke/mongoke

config-schema:
	cat configuration_schema.skema | skema gen jsonschema > ./mongoke/config_schema.json

.PHONY: play
generate-spec: clean
	python -m mongoke tests/confs/spec_conf.yaml --generated-path example_generated_code

.PHONY: play
play: generate-spec
	DB_URL=mongodb://localhost/db python -m example_generated_code

.PHONY: play
tests: generate-pr
	pytest

json-example:
	cat tests/confs/skema | skema gen jsonschema > tests/confs/schema.json
	python -m mongoke tests/confs/json_conf.yaml --generated-path play_json
	DB_URL=mongodb://localhost/db python -m example_generated_code