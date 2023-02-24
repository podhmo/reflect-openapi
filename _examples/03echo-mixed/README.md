---
title: Sample API
version: 0.0.0
---

# Sample API



- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `GET /users` | [main.ListUsers](#mainlistusers-get-users)  | `main` | ListUsers returns a list of users. |
| `POST /users` | [main.InsertUser](#maininsertuser-post-users)  | `main` | InsertUser inserts user. |
| `GET /users/{userId}` | [main.GetUser](#maingetuser-get-usersuserid)  | `main` | GetUser returns user |


### main.ListUsers `GET /users`

ListUsers returns a list of users.

| name | value | 
| --- | --- |
| operationId | main.ListUsers |
| endpoint | `GET /users` |
| tags | `main` |



#### output (application/json)

```go
// GET /users (200)
type Output200 []struct {	// User
	id integer

	// for go-playground/validator
	name string
}

// GET /users (default)
// default error
type OutputDefault struct {	// APIError
	message string

	details map[string]struct {	// FieldError
		path string

		message string
	}
}
```

#### description

ListUsers returns a list of users.
### main.InsertUser `POST /users`

InsertUser inserts user.

| name | value | 
| --- | --- |
| operationId | main.InsertUser |
| endpoint | `POST /users` |
| tags | `main` |


#### input (application/json)

```go
// POST /users
type Input struct {
	JSONBody struct {	// User
		id integer

		// for go-playground/validator
		name string
	}
}
```

#### output (application/json)

```go
// POST /users (200)
type Output200 struct {	// User
	id integer

	// for go-playground/validator
	name string
}

// POST /users (default)
// default error
type OutputDefault struct {	// APIError
	message string

	details map[string]struct {	// FieldError
		path string

		message string
	}
}
```

#### description

InsertUser inserts user.
### main.GetUser `GET /users/{userId}`

GetUser returns user

| name | value | 
| --- | --- |
| operationId | main.GetUser |
| endpoint | `GET /users/{userId}` |
| tags | `main` |


#### input (application/json)

```go
// GET /users/{userId}
type Input struct {
	userId integer `in:"path"`
}
```

#### output (application/json)

```go
// GET /users/{userId} (200)
type Output200 struct {	// User
	id integer

	// for go-playground/validator
	name string
}

// GET /users/{userId} (default)
// default error
type OutputDefault struct {	// APIError
	message string

	details map[string]struct {	// FieldError
		path string

		message string
	}
}
```

#### description

GetUser returns user



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [APIError](#apierror) |  |
| [FieldError](#fielderror) |  |
| [User](#user) |  |



### APIError

```go
type APIError struct {
	message string

	details map[string]struct {	// FieldError
		path string

		message string
	}
}
```

- [output of main.ListUsers (default) as `APIError`](#mainlistusers-get-users)
- [output of main.InsertUser (default) as `APIError`](#maininsertuser-post-users)
- [output of main.GetUser (default) as `APIError`](#maingetuser-get-usersuserid)

### FieldError

```go
type FieldError struct {
	path string

	message string
}
```


### User

```go
type User struct {
	id integer

	// for go-playground/validator
	name string
}
```

- [output of main.ListUsers (200) as `[]User`](#mainlistusers-get-users)
- [input of main.InsertUser as `User`](#maininsertuser-post-users)
- [output of main.InsertUser (200) as `User`](#maininsertuser-post-users)
- [output of main.GetUser (200) as `User`](#maingetuser-get-usersuserid)