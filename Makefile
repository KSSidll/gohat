PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@echo "build        build a binary"
	@echo "run          build and run local project"
	@echo "templ        generate templ templates"
	@echo "css          build tailwindcss"
	@echo "css-watch    watch build tailwindcss"
	@echo

.PHONY: build
build:
	go build -v -o ./build/main.exe ./src/main.go

.PHONY: run
start: build
	./build/main.exe

.PHONY: templ
templ:
	templ generate

.PHONY: css
css:
	npx tailwindcss -i ./src/assets/css/main-tailwind.css -o ./src/assets/css/main.min.css --minify

.PHONY: css-watch
css-watch:
	npx tailwindcss -i ./src/assets/css/main-tailwind.css -o ./src/assets/css/main.min.css --minify --watch