


.PHONY: clean
clean:
	rm -rf generated

.PHONY: play
generate-spec: clean
	python -m mongoke confs/spec_conf.yaml

.PHONY: play
play: generate-spec
	python -m generated

.PHONY: play
tests: generate-pr
	pytest


