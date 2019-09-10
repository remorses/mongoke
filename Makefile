


.PHONY: clean
clean:
	rm -rf generated

.PHONY: play
generate-pr: clean
	python -m src confs/pr_conf.yaml

.PHONY: play
play: generate-pr
	python -m generated

.PHONY: play
tests: generate-pr
	pytest


