.PHONY: install test test-e2e test-all

install:
	pip install -r requirements.txt

test: install
	pytest tests/ -v

test-e2e:
	bash scripts/test-e2e.sh

test-all: test test-e2e
