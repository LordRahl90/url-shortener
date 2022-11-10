.PHONY: bi
all: bi
	docker-compose up

start:
	@go run ./cmd/

test:
	@go test ./... --cover

start-image:
	@docker-compose up

build-image:
	@docker build -t lordrahl/shortener:latest .

push-image:
	@docker push lordrahl/shortener:latest

bi: build-image
pi: push-image
si: start-image