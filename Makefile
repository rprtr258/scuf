.PHONY: test
test:
	@gotestsum --format dots-v2 ./...

.PHONY: test-watch
test-watch:
	@reflex -r '\.go$$' -- make test

.PHONY: run-color-chart
run-color-chart:
	@go run cmd/color-chart/main.go

.PHONY: hello-world
run-hello-world:
	@go run cmd/hello-world/main.go

.PHONY: run-ssh
run-ssh:
	@go run cmd/ssh/main.go