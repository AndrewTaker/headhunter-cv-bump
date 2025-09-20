ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: all web scheduler

all: web scheduler

web:
	@echo "Building web..."
	cd apps/web && go build -o bin/hhcv .

scheduler:
	@echo "Building scheduler..."
	cd apps/scheduler && go build -o bin/hhcv-scheduler .

run-web:
	@echo "Building web..."
	cd apps/web && go build -o bin/hhcv .
	@echo "Starting web..."
	./web/bin/hhcv http localhost 44444 true
