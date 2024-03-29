package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/LucasAndFlores/user_api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	POSTGRES_DB   = "my_db_test"
	POSTGRES_USER = "postgres"
	DB_HOST       = "localhost"
)

var DB_PORT string

func TestMain(t *testing.M) {
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()

	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%v", POSTGRES_DB),
			"POSTGRES_HOST_AUTH_METHOD=trust",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	pg.Expire(10)

	postgresPort := pg.GetPort("5432/tcp")

	os.Setenv("POSTGRES_PORT", postgresPort)

	DB_PORT = postgresPort

	if err := pool.Retry(func() error {
		_, connErr := gorm.Open(postgres.Open(fmt.Sprintf("postgresql://postgres@localhost:%s/%s", postgresPort, POSTGRES_DB)), &gorm.Config{})
		if connErr != nil {
			return connErr
		}

		return nil
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	code := t.Run()

	os.Exit(code)
}

func runTestServer() *fiber.App {
	os.Setenv("DB_HOST", DB_HOST)
	os.Setenv("POSTGRES_USER", POSTGRES_USER)
	os.Setenv("POSTGRES_DB", POSTGRES_DB)
	os.Setenv("DB_PORT", DB_PORT)

	db, err := database.ConnectDatabase()

	if err != nil {
		log.Fatalf("An error occurred when tried to connect to database: %v", err)
	}

	return SetupApp(db)
}

func TestCreateUserSuccessfulScenario(t *testing.T) {
	tApp := runTestServer()

	body := map[string]interface{}{
		"name":          "test user",
		"email":         "test@example.com",
		"id":            "54022f9e-2301-428f-80de-ba73273341fb",
		"date_of_birth": "1990-01-01T00:00:00Z",
	}

	request, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Failed to marshal paylod to JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/save", bytes.NewReader(request))

	req.Header.Set("content-type", "application/json")

	resp, err := tApp.Test(req, -1)

	if err != nil {
		t.Fatalf("Failed when trying to execute fiber.Test: %v", err)
	}

	value, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Failed when trying to execute io.ReadAll: %v", err)
	}

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("The result is different from expected. Result: %v. Expected: %v", resp.StatusCode, fiber.StatusCreated)
	}

	expectedBody := "{\"message\":\"user successfully created\"}"

	if string(value) != expectedBody {
		t.Fatalf("The body result is different from expected. Result: %v. Expected: %v", string(value), expectedBody)
	}

}

type CreateUserErrorTest struct {
	request            map[string]interface{}
	expectedStatusCode int
	expectedBody       string
}

func TestCreateUserErrorScenario(t *testing.T) {
	ts := runTestServer()

	testCases := []CreateUserErrorTest{
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user1@example.com",
				"id":            "2dd002d0-dd56-4491-b77e-61b7dcce7123",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusCreated,
			expectedBody:       "{\"message\":\"user successfully created\"}",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user1@example.com",
				"id":            "3a47386e-56d4-4bd8-a015-c2b8bdf646f8",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusConflict,
			expectedBody:       "{\"message\":\"user already exists\"}",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2@example.com",
				"id":            "2dd002d0-dd56-4491-b77e-61b7dcce7123",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusConflict,
			expectedBody:       "{\"message\":\"user already exists\"}",
		},
	}

	for i, value := range testCases {

		request, err := json.Marshal(value.request)

		if err != nil {
			t.Fatalf("Failed to marshal paylod to JSON: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/save", bytes.NewReader(request))

		req.Header.Set("content-type", "application/json")

		resp, err := ts.Test(req, -1)

		if err != nil {
			t.Fatalf("Failed when trying to execute fiber.Test: %v", err)
		}

		rBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("Failed when trying to execute io.ReadAll: %v", err)
		}

		if resp.StatusCode != value.expectedStatusCode {
			t.Fatalf("Result is different from expected. Result: %v. Expected: %v. Test case index: %v", resp.StatusCode, value.expectedStatusCode, i)
		}

		if string(rBody) != value.expectedBody {
			t.Fatalf("The body result is different from expected. Result: %v. Expected: %v. Test case index: %v", string(rBody), value.expectedBody, i)
		}

	}

}

func TestCreateUserInvalidRequestScenario(t *testing.T) {
	ts := runTestServer()

	testCases := []CreateUserErrorTest{
		{
			request: map[string]interface{}{
				"name":          "",
				"email":         "user1@example.com",
				"id":            "497980c5-dac6-4af7-ac56-a7a0d2dad51a",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"Name\",\"Tag\":\"required\",\"Value\":\"\"}]",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2",
				"id":            "3a47386e-56d4-4bd8-a015-c2b8bdf646f8",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"Email\",\"Tag\":\"email\",\"Value\":\"\"}]",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2@example.com",
				"id":            "",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"ExternalId\",\"Tag\":\"uuid\",\"Value\":\"\"}]",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2@example.com",
				"id":            "testestesteststes",
				"date_of_birth": "1990-01-01T00:00:00Z",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"ExternalId\",\"Tag\":\"uuid\",\"Value\":\"\"}]",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2@example.com",
				"id":            "3a47386e-56d4-4bd8-a015-c2b8bdf646f8",
				"date_of_birth": "",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"DateOfBirth\",\"Tag\":\"required\",\"Value\":\"invalid date format\"},{\"Field\":\"DateOfBirth\",\"Tag\":\"required\",\"Value\":\"\"}]",
		},
		{
			request: map[string]interface{}{
				"name":          "test user",
				"email":         "user2@example.com",
				"id":            "3a47386e-56d4-4bd8-a015-c2b8bdf646f8",
				"date_of_birth": "1990-01-01",
			},
			expectedStatusCode: fiber.StatusUnprocessableEntity,
			expectedBody:       "[{\"Field\":\"DateOfBirth\",\"Tag\":\"required\",\"Value\":\"invalid date format\"}]",
		},
	}

	for i, value := range testCases {

		request, err := json.Marshal(value.request)

		if err != nil {
			t.Fatalf("Failed to marshal paylod to JSON: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/save", bytes.NewReader(request))

		req.Header.Set("content-type", "application/json")

		resp, err := ts.Test(req, -1)

		if err != nil {
			t.Fatalf("Failed when trying to execute fiber.Test: %v", err)
		}

		rBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("Failed when trying to execute io.ReadAll: %v", err)
		}

		if resp.StatusCode != value.expectedStatusCode {
			t.Fatalf("Result is different from expected. Result: %v. Expected: %v. Test case index: %v", resp.StatusCode, value.expectedStatusCode, i)
		}

		if string(rBody) != value.expectedBody {
			t.Fatalf("The body result is different from expected. Result: %v. Expected: %v. Test case index: %v", string(rBody), value.expectedBody, i)
		}

	}

}

func TestFindUserByIdSuccessfulScenario(t *testing.T) {
	tApp := runTestServer()

	createdUser := map[string]interface{}{
		"name":          "test user",
		"email":         "test_create@example.com",
		"id":            "d553a9de-eff6-4b3d-9c70-8c9266692782",
		"date_of_birth": "1990-01-01T00:00:00Z",
	}

	request, err := json.Marshal(createdUser)

	if err != nil {
		t.Fatalf("Failed to marshal paylod to JSON: %v", err)
	}

	reqCreate := httptest.NewRequest("POST", "/api/save", bytes.NewReader(request))

	reqCreate.Header.Set("content-type", "application/json")

	_, err = tApp.Test(reqCreate, -1)

	if err != nil {
		t.Fatalf("Failed when trying to execute fiber.Test to create user: %v", err)
	}

	getUrl := fmt.Sprintf("/api/%v", createdUser["id"])

	req := httptest.NewRequest("GET", getUrl, nil)

	res, err := tApp.Test(req, -1)

	if err != nil {
		t.Fatalf("Failed when trying to execute fiber.Test: %v", err)
	}

	value, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed when trying to execute io.ReadAll: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("The result is different from expected. Result: %v. Expected: %v", res.StatusCode, fiber.StatusOK)
	}

	expectedBody := "{\"user\":{\"name\":\"test user\",\"email\":\"test_create@example.com\",\"id\":\"d553a9de-eff6-4b3d-9c70-8c9266692782\",\"date_of_birth\":\"1990-01-01 01:00:00 +0100 CET\"}}"

	if string(value) != expectedBody {
		t.Fatalf("The body result is different from expected. Result: %v. Expected: %v", string(value), expectedBody)
	}
}

type FindUserByIdErrorTest struct {
	getUrl             string
	expectedStatusCode int
	expectedBody       string
}

func TestFindUserByIdErrorScenario(t *testing.T) {
	tApp := runTestServer()

	testCases := []FindUserByIdErrorTest{
		{
			getUrl:             "/api/111111111111",
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBody:       "{\"message\":\"unable to parse the id\"}",
		},
		{
			getUrl:             "/api/testestestest",
			expectedStatusCode: fiber.StatusBadRequest,
			expectedBody:       "{\"message\":\"unable to parse the id\"}",
		},
		{
			getUrl:             "/api/8269b23f-1417-4f9d-9662-83b609a4e6dd",
			expectedStatusCode: fiber.StatusNotFound,
			expectedBody:       "{\"message\":\"user not found\"}",
		},
	}

	for i, value := range testCases {
		req := httptest.NewRequest("GET", value.getUrl, nil)

		res, err := tApp.Test(req, -1)

		if err != nil {
			t.Fatalf("Failed when trying to execute fiber.Test: %v", err)
		}

		rBody, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Failed when trying to execute io.ReadAll: %v", err)
		}

		if res.StatusCode != value.expectedStatusCode {
			t.Fatalf("Result is different from expected. Result: %v. Expected: %v. Test case index: %v", res.StatusCode, value.expectedStatusCode, i)
		}

		if string(rBody) != value.expectedBody {
			t.Fatalf("The body result is different from expected. Result: %v. Expected: %v. Test case index: %v", string(rBody), value.expectedBody, i)
		}
	}

}
