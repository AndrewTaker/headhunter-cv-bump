ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: all web scheduler cli

all: web scheduler cli

web:
	@echo "Building web..."
	cd apps/web && go build -o bin/web .

scheduler:
	@echo "Building scheduler..."
	cd apps/scheduler && go build -o bin/scheduler .

cli:
	@echo "Building cli..."
	cd apps/cli && go build -o bin/cli .
