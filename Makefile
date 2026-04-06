include .env

PROJECT_NAME=inventory
MIGRATION_STEP=
MIGRATION_PATH=script/migration
DB_CONN=postgres://$(DB.POSTGRESQL.WRITE.USER:"%"=%):$(DB.POSTGRESQL.WRITE.PASSWORD:"%"=%)@$(DB.POSTGRESQL.WRITE.HOST:"%"=%):$(DB.POSTGRESQL.WRITE.PORT:"%"=%)/$(DB.POSTGRESQL.WRITE.NAME:"%"=%)?sslmode=$(DB.POSTGRESQL.WRITE.SSLMODE:"%"=%)

# development
dev: generate
	go run github.com/cosmtrek/air

# db migration
migrate_create:
	@read -p "migration name (do not use space): " NAME \
  	&& migrate create -ext sql -dir $(MIGRATION_PATH) $${NAME}

# windows migration
# @migrate create -ext sql -dir $(MIGRATION_PATH) $(NAME)

migrate_up:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_CONN)" up $(MIGRATION_STEP)

migrate_down:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_CONN)" down $(MIGRATION_STEP)

migrate_force:
	@read -p "please enter the migration version (the migration filename prefix): " VERSION \
  	&& migrate -path $(MIGRATION_PATH) -database "$(DB_CONN)" force $${VERSION}

migrate_version:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_CONN)" version 

migrate_drop:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_CONN)" drop

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/vektra/mockery/v2@latest
	swag init

local: generate
	go run github.com/cosmtrek/air -c config/.air.toml

test:
	go test -v -cover -coverprofile=cover.out ./...

# GENERATE PROTO
# Repository Module
# example: $ make gen-{service_name}
gen-%:
	# proto generate for request
	mkdir -p entity
	protoc -I/usr/include --proto_path=proto/$* --go_out=entity --go_opt=paths=source_relative proto/$*/$*.proto;
	# proto generate for adapter
	mkdir -p module/repository/grpc/$*
	protoc -I/usr/include --proto_path=proto/$* --go-grpc_out=module/repository/grpc/$* \
		--go-grpc_opt=paths=source_relative proto/$*/$*_service.proto;


lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

dev-up:
	docker-compose -f docker-compose-dev.yml up -d

dev-down:
	docker-compose -f docker-compose-dev.yml down

prod-up:
	docker-compose up -d

prod-down:
	docker-compose down

scan:
	sonar-scanner

mockery-usecase:
	cd module/auth/usecase && mockery --name=Usecase --output=../mocks 
	
mockery-repo:
	cd module/auth/store && mockery --name=Repository --output=../mocks

rebase-dev:
	git checkout development && git pull origin development && git checkout @{-1} && git rebase development

rebase-release:
	git checkout release && git pull origin release && git checkout @{-1} && git rebase release

rebase-master:
	git checkout master && git pull origin master && git checkout @{-1} && git rebase master

push:
	go fmt ./... && git push origin HEAD

push-force:
	go fmt ./... && git push -f origin HEAD

generate:
	go generate ./...
