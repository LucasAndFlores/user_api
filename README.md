# User API 
This is an API with two available endpoints, where you can save and return the data stored in a Postgres database. The API is built using Go, Fiber, Go ORM, and Postgresql.

## Usage and setup 
Before you start, be sure that you have Docker, Docker-compose, Golang, and Makefile installed.

As an initial step, copy all the variables from the `.env.example` and create a `.env` file. Define your `POSTGRES_USER` and `POSTGRES_PASSWORD` variables.

### Local usage - Docker
To run the API with docker, run: 

```bash
docker-compose up --build
```

And you will be able to access all the endpoints.

### Testing
To run the integration tests, you need to install all the packages included in the API, run:

```bash
go get -d -v ./...
```

After this, you can execute the integration tests with this command:
```bash
make run-integration-test
```

This command automatically generates a coverage report, if you want to check this report, run:
```bash
make coverage-report
```

## API Endpoints
`POST /api/save`

This route is responsible for storing the user in the database. 
A body object is required. The email and id are unique.

Example: 
```json
{
	"name":          "John", 
	"email":         "john@test.com",
	"id":            "4e6e0b08-f1a7-4ff3-85d3-f93fabc8ad5d",
	"date_of_birth": "1990-01-01T00:00:00Z",
}
```

Expected responses:

Status code: `201` <br>
Body: <br>
```json
{
   "message":"user successfully created"
}
```

Status code: `409` <br>
Error reason: The email or ID is already registered for a user and it can't be registered again. <br>
Body:
```json
{
   "message":"user already exists"
}
```

Status code: `422` <br>
Error reason: A field is missing or it is invalid. <br>
Body:
```json
[
   {
      "Field":"Name",
      "Tag":"Required",
      "Value":""
   }
]
```

Status code: `500` <br>
Error reason: An internal server error happened <br>
Body:
```json
{
	"message":"internal server error",
}
```


`GET /api/:id`

This endpoint will return user data or an error if the user doesn't exist.

Expected responses:

Status code: `200` <br>
Body:
```json
{
	"name":          "John", 
	"email":         "john@test.com",
	"id":            "4e6e0b08-f1a7-4ff3-85d3-f93fabc8ad5d",
	"date_of_birth": "1990-01-01T00:00:00Z",
}
```

Status code: `400` <br>
Error reason: Invalid ID type. <br>
Body:
```json
{
	"message":"unable to parse the id",
}
```

Status code: `404` <br>
Error reason: The user does not exist in the database <br>
Body:
```json
{
	"message":"user not found",
}
```

Status code: `500` <br>
Error reason: An internal server error happened <br>
Body:
```json
{
	"message":"internal server error",
}
```

