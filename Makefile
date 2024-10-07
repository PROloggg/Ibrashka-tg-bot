build-docker:
	@docker compose up -d --build --remove-orphans

up:
	@docker compose up -d

sh:
	@docker compose run go-app sh

down:
	@docker compose down

down-v:
	@docker compose down -v

logs:
	@docker compose logs go-app

build-exe:	prepare-build
	@docker compose run go-app sh -c "GOOS=windows GOARCH=amd64 go build -o build/Ibrashka.exe main.go"

build-unix:	prepare-build
	@docker compose run go-app sh -c "go build -o build/Ibrashka main.go"

prepare-build:
	mkdir -p app/build
	cp app/.env app/build/
	cp -r app/src app/build/src