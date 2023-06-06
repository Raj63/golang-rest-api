build:
	docker build -t golang-rest-api .

run:
	docker-compose down
	docker-compose up

reload:
	docker-compose down
	docker image rm golang-rest-api
	docker image rm golang-rest-api_app
	docker build -t golang-rest-api .
	docker-compose up

test:
	go test -count=1 -failfast -v -race ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html && go tool cover -func coverage.out

# Push the Docker image to the registry
push:
    docker tag golang-rest-api <your-registry>/golang-rest-api:<version>
    docker push <your-registry>/golang-rest-api:<version>

generate:
	swag init -g pkg/infrastructure/rest/routes/routes.go

db-update:
	docker run golang-rest-api update-db
	
