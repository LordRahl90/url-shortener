.PHONY: bi
all: bi
	docker-compose up

start:
	@go run ./cmd/

test:
	@go test ./... -v --cover

start-image:
	@docker-compose up

build-image:
	@docker build -t gcr.io/neurons-be-test/shortener:latest .

push-image:
	@docker push gcr.io/neurons-be-test/shortener:latest

bi: build-image
pi: push-image
si: start-image