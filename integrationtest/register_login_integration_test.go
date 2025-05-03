//go:build integrationtest

package integrationtest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"time"

	"github.com/fgeck/go-register/internal/service/config"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/internal/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

const (
	POSTGRES_IMAGE    = "postgres:latest"
	POSTGRES_USER     = "testuser"
	POSTGRES_PASSWORD = "testpassword"
	POSTGRES_DB       = "postgres"
)

var postgresCfg = PostgresConfig{
	Image:    POSTGRES_IMAGE,
	Username: POSTGRES_USER,
	Password: POSTGRES_PASSWORD,
	Database: POSTGRES_DB,
}

func TestIntegrationRegisterLogin(t *testing.T) {
	var err error
	postgresContainer, host, port, err := StartPostgres(postgresCfg)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Set environment variables for the app
	//nolint
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", POSTGRES_USER)
	os.Setenv("DB_PASSWORD", POSTGRES_PASSWORD)
	os.Setenv("DB_DATABASE", POSTGRES_DB)
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	cfgDirPath := filepath.Join("../", "cmd/", "web/")
	cfg, err := config.NewLoader().LoadConfig(cfgDirPath)
	require.NoError(t, err)

	go func(t *testing.T) {
		e := echo.New()
		web.InitServer(e, cfg)
		if err := e.Start(":" + cfg.App.Port); err != nil {
			log.Printf("failed to start server: %s", err)
			t.Fatal(err)
		}
	}(t)
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
