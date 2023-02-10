## Documentation

Postman collection: [https://documenter.getpostman.com/view/14566466/2s935snMUV](https://documenter.getpostman.com/view/14566466/2s935snMUV)

## Usage

1. clone repository
2. change `.env` values - optional
3. `$ docker-compose up -d` - Run docker-compose
4. Generate test token with JWT endpoint
5. Test other endpoints using postman collection above

## Tests

- `$ go test httpserver.go httpserver_test.go -v` - run tests
