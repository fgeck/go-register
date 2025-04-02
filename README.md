# go-register
Template repo for user registration and signin based on goth stack

## Prerequisites

```bash
brew install sqlc tailwindcss
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
```

## How to run for Development

```bash
tailwindcss -i templates/css/app.css -o public/styles.css --watch
templ generate --watch
air
```

...or run each command in a dedicated terminal window

## Tools used

* [echo](https://echo.labstack.com/)
* [sqlc](https://sqlc.dev/)
* [templ](https://github.com/a-h/templ)
* [air](https://github.com/air-verse/air)
* [tailwindcss](https://tailwindcss.com/)
* [daisyui](https://daisyui.com/)

## ToDO

* [x] Echo Server renders templ templates
* [x] Add Tailwind to templates
* [ ] Add htmx to templates
* [ ] Add DaisyUI
* [ ] 
