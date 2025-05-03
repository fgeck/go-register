//go:build integrationtest

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"time"

	"github.com/fgeck/go-register/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	POSTGRES_IMAGE    = "postgres:latest"
	POSTGRES_USER     = "testuser"
	POSTGRES_PASSWORD = "testpassword"
	POSTGRES_DB       = "postgres"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	initScripts, err := getInitScripts()
	require.NoError(t, err)

	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	defer os.Unsetenv("TESTCONTAINERS_RYUK_DISABLED")
	postgresContainer, err := postgres.Run(ctx,
		POSTGRES_IMAGE,
		postgres.WithInitScripts(initScripts...),
		postgres.WithDatabase(POSTGRES_DB),
		postgres.WithUsername(POSTGRES_USER),
		postgres.WithPassword(POSTGRES_PASSWORD),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	// Set environment variables for the app
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", POSTGRES_USER)
	os.Setenv("DB_PASSWORD", POSTGRES_PASSWORD)
	os.Setenv("DB_DATABASE", POSTGRES_DB)
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	go func() {
		main()
	}()
	time.Sleep(1 * time.Second)

	t.Run("A new user can register", func(t *testing.T) {
		testUser := "testuser"
		testEmail := "testuser@test.io"
		testPassword := "testuserPassword123!"

		formData := url.Values{
			"username": {testUser},
			"email":    {testEmail},
			"password": {testPassword},
		}

		resp, err := http.PostForm("http://localhost:8081/api/register", formData)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		var createdUser user.UserCreatedDto
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		require.NoError(t, err)
		assert.Equal(t, createdUser.Username, testUser)
		assert.Equal(t, createdUser.Email, testEmail)
	})

	t.Run("A user cannot register with an existing email", func(t *testing.T) {
		testUser := "othertestuser"
		testEmail := "othertestuser@test.io"
		testPassword := "othertestuserPassword123!"

		formData := url.Values{
			"username": {testUser},
			"email":    {testEmail},
			"password": {testPassword},
		}

		resp, err := http.PostForm("http://localhost:8081/api/register", formData)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		var createdUser user.UserCreatedDto
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		require.NoError(t, err)
		assert.Equal(t, createdUser.Username, testUser)
		assert.Equal(t, createdUser.Email, testEmail)

		resp, err = http.PostForm("http://localhost:8081/api/register", formData)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("A registered user can login", func(t *testing.T) {
		testUser := "anothertestuser"
		testEmail := "anothertestuser@test.io"
		testPassword := "anothertestuserPassword123!"

		formData := url.Values{
			"username": {testUser},
			"email":    {testEmail},
			"password": {testPassword},
		}

		resp, err := http.PostForm("http://localhost:8081/api/register", formData)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		var createdUser user.UserCreatedDto
		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		require.NoError(t, err)
		assert.Equal(t, createdUser.Username, testUser)
		assert.Equal(t, createdUser.Email, testEmail)

		formData = url.Values{
			"email":    {testEmail},
			"password": {testPassword},
		}
		resp, err = http.PostForm("http://localhost:8081/api/login", formData)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cookies := resp.Cookies()
		require.NotEmpty(t, cookies, "No cookies found in the response")

		var tokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "token" {
				tokenCookie = cookie
				break
			}
		}
		require.NotNil(t, tokenCookie, "Token cookie not found in the response")
		assert.NotEmpty(t, tokenCookie.Value, "Token cookie value is empty")
		assert.True(t, tokenCookie.HttpOnly, "Token cookie is not HttpOnly")
		assert.True(t, tokenCookie.Secure, "Token cookie is not Secure")
		assert.Equal(t, "/", tokenCookie.Path, "Token cookie path is incorrect")
		assert.Equal(t, http.SameSiteLaxMode, tokenCookie.SameSite, "Token cookie SameSite attribute is incorrect")
	})
}

func getInitScripts() ([]string, error) {
	var initScripts []string
	err := filepath.Walk("../../migrations", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		initScripts = append(initScripts, path)
		return nil
	})

	return initScripts, err
}
