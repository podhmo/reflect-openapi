---
title: hello
version: 1.0.0
---

# hello



- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `POST /api/hello` | [main.Hello](#mainhello-post-apihello)  | `main` |  |

htmls

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `GET /hello/{name}` | [main.HelloHTML](#mainhellohtml-get-helloname)  | `main` |  |


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

