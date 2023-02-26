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

