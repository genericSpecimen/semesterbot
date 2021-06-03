VENV := env
BIN := $(VENV)/bin
PYTHON := $(BIN)/python3
SHELL := /bin/bash

.PHONY: venv
venv:
	python3 -m venv $(VENV) && source $(BIN)/activate

freeze:
	$(BIN)/pip freeze > requirements.txt

.PHONY: install
install: venv
	$(BIN)/pip install --upgrade -r requirements.txt

start: install
