build:
	go build -i -v

fmt:
	go fmt ./...

image:
	docker build -t cirocosta/perfer .
