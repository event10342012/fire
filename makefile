.PHONY: docker
docker:
	@rm fire || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o fire .
	@podman rmi -f event10342012/fire:v0.0.1
	@podman build -t event10342012/fire:v0.0.1 .
