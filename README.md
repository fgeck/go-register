# go-register

Template repo for user registration and signin based on GoTTH (Go, Tailwind, Templ, Htmx) stack

## Prerequisites

```bash
brew install sqlc tailwindcss mockery
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
go mod download
```

## How to run for Development

```bash
tailwindcss -i templates/css/app.css -o public/styles.css --watch &
templ generate --watch &
air &
```

...or run each command in a dedicated terminal window

## Tools used

* [echo](https://echo.labstack.com/)
* [sqlc](https://sqlc.dev/)
* [templ](https://github.com/a-h/templ)
* [air](https://github.com/air-verse/air)
* [tailwindcss](https://tailwindcss.com/)
* [daisyui](https://daisyui.com/) ?

## ToDO

* [x] Echo Server renders templ templates
* [x] Add Tailwind to templates
* [x] Add htmx to templates
* [ ] Add DaisyUI?
* [ ] On startup register admin user based on env
* [ ] build register handler for storing in db
* [ ] build login handler emitting JWT with proper claims
* [ ] Loggingframework like logrus
* [ ] Dependency Injection Framework?
* [ ] 
