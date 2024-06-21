modules := $(notdir $(basename $(wildcard schema/*.yml)))
generator := $(shell find internal/generate)

pkg/opnsense/client.go: $(generator) $(wildcard schema/*.yml)
	@echo "Generating opnsense client"
	go generate ./pkg/opnsense

pkg/%/controller.go: schema/%.yml $(generator) pkg/%/generate.go
	@echo "Generating $* controller"
	go generate ./$(@D)

$(modules): %: pkg/%/controller.go

all: $(modules) pkg/opnsense/client.go
