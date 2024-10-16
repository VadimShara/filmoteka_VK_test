lint:
	@golangci-lint --timeout=2m run

docker-up:
	docker compose --env-file deploy/.env -f deploy/docker-compose.yml -p filmoteka up -d --build

docker-down:
	docker compose -p filmoteka down

docker-logs:
	docker compose -p filmoteka logs

swagger-doc-generate:
	go-swagger3 --main-file-path cmd/main.go --handler-path api/rest/handlers --output api/doc/swagger.json --schema-without-pkg

wire:
	cd api/inject && wire