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
| `GET /users` | [main.ListUser](#mainlistuser-get-users)  |  |  |


### main.ListUser `GET /users`



| name | value |
| --- | --- |
| operationId | main.ListUser |
| endpoint | `GET /users` |
| tags |  |


#### input (application/json)

```go
// GET /users
type Input struct {
	cursor? string `in:"query"`
	pageSize? integer `in:"query"`
	sort? "asc" | "desc" `in:"query"`
	query? string `in:"query"`
}
```

#### output (application/json)

```go
// GET /users (200)
type Output200 struct {	// PaginatedOutput[[]main.User]
	hasMore boolean

	cursor string

	nextCursor string

	items []struct {	// User
		name string
	}
}
```





----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [PaginatedOutput__main.User](#paginatedoutput__mainuser) |  |
| [User](#user) |  |



### PaginatedOutput__main.User

```go
type PaginatedOutput[[]main.User] struct {
	hasMore boolean

	cursor string

	nextCursor string

	items []struct {	// User
		name string
	}
}
```

- [output of main.ListUser (200) as `PaginatedOutput[[]main.User]`](#mainlistuser-get-users)

### User

```go
type User struct {
	name string
}
```
