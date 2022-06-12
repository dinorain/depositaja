start:
	docker-compose up -d

dep:
	go mod tidy
	go mod vendor

service:
	go run cmd/service/main.go

processor:
	go run cmd/processor/main.go -collector -detector -flagger