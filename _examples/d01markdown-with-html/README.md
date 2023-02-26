---
title: hello
version: 1.0.0
---

# hello

This is the example has text/html output

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `POST /api/hello` | [main.Hello](#mainhello-post-apihello)  | `main` |  |
| `GET /hello/{name}` | [main.HelloHTML](#mainhellohtml-get-helloname)  | `main text/html` |  |
| `GET /hello2/{name}` | [main.HelloHTML2](#mainhellohtml2-get-hello2name)  | `main text/html` | with custom error response (responses['default']) |
| `GET /hello3/{name}` | [main.HelloHTML3](#mainhellohtml3-get-hello3name)  | `main text/html` | with response header |
| `POST /login` | [main.Login](#mainlogin-post-login)  | `main text/html` | Successfully authenticated. |


### main.Hello `POST /api/hello`



| name | value |
| --- | --- |
| operationId | main.Hello |
| endpoint | `POST /api/hello` |
| tags | `main` |


#### input (application/json)

```go
// POST /api/hello
type Input struct {
	JSONBody struct {	// 
		name string
	}
}
```

#### output (application/json)

```go
// POST /api/hello (200)
type Output200 struct {	// 
	message string
}
```


### main.HelloHTML `GET /hello/{name}`



| name | value |
| --- | --- |
| operationId | main.HelloHTML |
| endpoint | `GET /hello/{name}` |
| tags | `main` |


#### input

```go
// GET /hello/{name}
type Input struct {
	name string `in:"path"`
}
```

#### output (text/html)

html with greeting message
### main.HelloHTML2 `GET /hello2/{name}`

with custom error response (responses['default'])

| name | value |
| --- | --- |
| operationId | main.HelloHTML2 |
| endpoint | `GET /hello2/{name}` |
| tags | `main` |


#### input

```go
// GET /hello2/{name}
type Input struct {
	name string `in:"path"`
}
```

#### output (text/html)

html with greeting message

#### description

with custom error response (responses['default'])
### main.HelloHTML3 `GET /hello3/{name}`

with response header

| name | value |
| --- | --- |
| operationId | main.HelloHTML3 |
| endpoint | `GET /hello3/{name}` |
| tags | `main` |


#### input

```go
// GET /hello3/{name}
type Input struct {
	name string `in:"path"`
}
```

#### output (text/html)

html with greeting message

#### description

with response header
### main.Login `POST /login`

Successfully authenticated.

| name | value |
| --- | --- |
| operationId | main.Login |
| endpoint | `POST /login` |
| tags | `main` |


#### input

```go
// POST /login
type Input struct {
	JSONBody struct {	// LoginInput
		name string

		password string
	}
}
```

#### output (text/html)



#### description

Successfully authenticated.
The session ID is returned in a cookie named `JSESSIONID`. You need to include this cookie in subsequent request



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [Error](#error) | is custom error response |
| [LoginInput](#logininput) | https://swagger.io/docs/specification/authentication/cookie-authentication/ |



### Error

```go
// Error is custom error response
type Error struct {
	message string
}
```

- [output of main.HelloHTML2 (default) as `Error`](#mainhellohtml2-get-hello2name)
- [output of main.HelloHTML3 (default) as `Error`](#mainhellohtml3-get-hello3name)

### LoginInput

```go
// https://swagger.io/docs/specification/authentication/cookie-authentication/
type LoginInput struct {
	name string

	password string
}
```

- [input of main.Login as `LoginInput`](#mainlogin-post-login)