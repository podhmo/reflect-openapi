---
title: Swagger Petstore
version: 1.0.0
---

# Swagger Petstore

A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `GET /users` | [main.ListUser](#mainlistuser-get-users)  | `main` |  |
| `GET /users/{id}` | [main.GetUser](#maingetuser-get-usersid)  | `main` | get user |


### main.ListUser `GET /users`



| name | value |
| --- | --- |
| operationId | main.ListUser[  <sub>(source)</sub>](https://github.com/podhmo/reflect-openapi/blob/main/_examples/d02markdown-with-wrapper/main.go#L62) |
| endpoint | `GET /users` |
| input | Input |
| output | [`Pagination[[]User]`](#user) ｜ [`Error`](#error) |
| tags | `main` |



#### output (application/json)

```go
// GET /users (200)
type Output200 struct {	// Pagination[main.User]
	hasMore boolean	// default: false

	items []struct {	// User
		// Name of the user
		name string

		// Age of the user
		age? integer
	}
}

// GET /users (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// Error message
	message string
}
```

examples

```json

// GET /users (default)

{
  "code": 444,
  "message": "unexpected error!"
}
```
### main.GetUser `GET /users/{id}`

get user

| name | value |
| --- | --- |
| operationId | main.GetUser[  <sub>(source)</sub>](https://github.com/podhmo/reflect-openapi/blob/main/_examples/d02markdown-with-wrapper/main.go#L53) |
| endpoint | `GET /users/{id}` |
| input | Input |
| output | [`GetUserOutput[User]`](#user) ｜ [`Error`](#error) |
| tags | `main` |


#### input (application/json)

```go
// GET /users/{id}
type Input struct {
	pretty? boolean `in:"query"`

	id string `in:"path"`
}
```

#### output (application/json)

```go
// GET /users/{id} (200)
type Output200 struct {	// GetUserOutput
	user struct {	// User
		// Name of the user
		name string

		// Age of the user
		age? integer
	}
}

// GET /users/{id} (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// Error message
	message string
}
```

examples

```json

// GET /users/{id} (default)

{
  "code": 444,
  "message": "unexpected error!"
}
```

#### description

get user





----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [Error](#error) |  |
| [User](#user) |  |



### Error



```go
type Error struct {
	// Error code
	code integer `format:"int32"`

	// Error message
	message string
}
```

exmaples

```js
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```

- [output of main.ListUser (default) as `Error`](#mainlistuser-get-users)
- [output of main.GetUser (default) as `Error`](#maingetuser-get-usersid)

### User



```go
type User struct {
	// Name of the user
	name string

	// Age of the user
	age? integer
}
```

- [output of main.ListUser (200) as `Pagination[[]User]`](#mainlistuser-get-users)
- [output of main.GetUser (200) as `GetUserOutput[User]`](#maingetuser-get-usersid)