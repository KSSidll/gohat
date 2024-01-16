<h1 align="center">Gohat</h1>

</br>

# Tech Stack & Libraries

- [Go](https://go.dev/) based
- [HTMX](https://htmx.org/)
- [Templ](https://github.com/a-h/templ)
- [Tailwind](https://tailwindcss.com/)

# Running

## Requirements

[Go](https://go.dev/)

# How to run

```shell
go run .
```

# Development

## Requirements

[Go](https://go.dev/)
[Templ](https://github.com/a-h/templ)
[Tailwind](https://tailwindcss.com/)

# How to run

generate templates
```shell
templ generate
```

generate tailwind css
```shell
npx tailwindcss -i ./src/assets/css/main-tailwind.css -o ./src/assets/css/main.min.css --minify
```

```shell
go run ./src
```