build:
	docker build -t golang-rest-api .

run:
	docker run -p 8080:8080 golang-rest-api

test:
	go test -count=1 -failfast -v -race ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html && go tool cover -func coverage.out

# Push the Docker image to the registry
push:
    docker tag golang-rest-api <your-registry>/golang-rest-api:<version>
    docker push <your-registry>/golang-rest-api:<version>

generate:
	swag init -g pkg/infrastructure/rest/routes/routes.go