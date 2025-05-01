.PHONY: docker
docker:
	@rm fire || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o fire .
	@docker rmi -f event10342012/fire:v0.0.1
	@docker build -t event10342012/fire:v0.0.1 .

mock:
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
