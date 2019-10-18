


.PHONY: clean
clean:
	rm -rf generated

config-schema:
	skema generate configuration_schema.skema --jsonschema ./mongoke/config_schema.json

.PHONY: play
generate-spec: clean
	python -m mongoke confs/spec_conf.yaml

.PHONY: play
play: generate-spec
	python -m generated

.PHONY: play
tests: generate-pr
	pytest


