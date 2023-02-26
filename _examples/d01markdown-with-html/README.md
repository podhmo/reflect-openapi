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
| `GET /hello2/{name}` | [main.HelloHTML2](#mainhellohtml2-get-hello2name)  | `main text/html` |  |


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



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [Error](#error) | is custom error response |



### Error

```go
// Error is custom error response
type Error struct {
	message string
}
```

- [output of main.HelloHTML2 (default) as `Error`](#mainhellohtml2-get-hello2name)