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
| `GET /pets` | [main.FindPets](#mainfindpets-get-pets)  | `pet read` | Returns all pets |
| `POST /pets` | [main.AddPet](#mainaddpet-post-pets)  | `pet write` | Creates a new pet |
| `DELETE /pets/{id}` | [main.DeletePet](#maindeletepet-delete-petsid)  | `pet write` | Deletes a pet by ID |
| `GET /pets/{id}` | [main.FindPetByID](#mainfindpetbyid-get-petsid)  | `pet read` | Returns a pet by ID |


### main.FindPets `GET /pets`

Returns all pets

| name | value | 
| --- | --- |
| operationId | main.FindPets |
| endpoint | `GET /pets` |


#### input (application/json)

```go
// GET /pets
type Input struct {
	// tags to filter by
	tags? []string `in:"query"`
	// maximum number of results to return
	limit? integer `in:"query"`
}
```

#### output (application/json)

```go
// GET /pets (200)
// pet response
type Output200 []struct {	// Pet
	// Unique id of the pet
	id integer `format:"int64"`
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}
```

exmaples

```js
// GET /pets
// sample output
[
  {
    "id": 1,
    "name": "foo",
    "tag": "A"
  },
  {
    "id": 2,
    "name": "bar",
    "tag": "A"
  },
  {
    "id": 3,
    "name": "boo",
    "tag": "B"
  }
]
```

#### output (application/json)

```go
// GET /pets (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`
	// Error message
	message string
}
```

exmaples

```js
// GET /pets
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```


#### description

Returns all pets from the system that the user has access to
Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.

Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
### main.AddPet `POST /pets`

Creates a new pet

| name | value | 
| --- | --- |
| operationId | main.AddPet |
| endpoint | `POST /pets` |


#### input (application/json)

```go
// POST /pets
type Input struct {
	JSONBody struct {	// AddPetInput
		// Name of the pet
		name string
		// Type of the pet
		tag? string
	}
}
```

#### output (application/json)

```go
// POST /pets (200)
// pet response TODO:
type Output200 struct {	// Pet
	// Unique id of the pet
	id integer `format:"int64"`
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}
```

#### output (application/json)

```go
// POST /pets (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`
	// Error message
	message string
}
```

exmaples

```js
// POST /pets
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```


#### description

Creates a new pet in the store. Duplicates are allowed
### main.DeletePet `DELETE /pets/{id}`

Deletes a pet by ID

| name | value | 
| --- | --- |
| operationId | main.DeletePet |
| endpoint | `DELETE /pets/{id}` |


#### input (application/json)

```go
// DELETE /pets/{id}
type Input struct {
	// ID of pet to delete
	id integer `in:"path"`
}
```

#### output (application/json)

```go
// DELETE /pets/{id} (204)
type Output204 struct {	// 
}
```

#### output (application/json)

```go
// DELETE /pets/{id} (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`
	// Error message
	message string
}
```

exmaples

```js
// DELETE /pets/{id}
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```


#### description

deletes a single pet based on the ID supplied
### main.FindPetByID `GET /pets/{id}`

Returns a pet by ID

| name | value | 
| --- | --- |
| operationId | main.FindPetByID |
| endpoint | `GET /pets/{id}` |


#### input (application/json)

```go
// GET /pets/{id}
type Input struct {
	// ID of pet to fetch
	id integer `in:"path"`
}
```

#### output (application/json)

```go
// GET /pets/{id} (200)
type Output200 struct {	// Pet
	// Unique id of the pet
	id integer `format:"int64"`
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}
```

#### output (application/json)

```go
// GET /pets/{id} (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`
	// Error message
	message string
}
```

exmaples

```js
// GET /pets/{id}
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```


#### description

Returns a pet based on a single ID



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [AddPetInput](#addpetinput) |  |
| [Error](#error) |  |
| [Pet](#pet) | pet object. |


### AddPetInput

```go
type AddPetInput struct {
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}

```
- [input of main.AddPet](#mainaddpet-post-pets)
### Error

```go
type Error struct {
	// Error code
	code integer `format:"int32"`
	// Error message
	message string
}

```
- [output of main.FindPets (default)](#mainfindpets-get-pets)
- [output of main.AddPet (default)](#mainaddpet-post-pets)
- [output of main.DeletePet (default)](#maindeletepet-delete-petsid)
- [output of main.FindPetByID (default)](#mainfindpetbyid-get-petsid)

exmaples

```js
// 
{
  "code": 444,
  "message": "unexpected error!"
}
```
### Pet

```go
// Pet : pet object.
type Pet struct {
	// Unique id of the pet
	id integer `format:"int64"`
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}

```
- [output of main.AddPet (200)](#mainaddpet-post-pets)
- [output of main.FindPetByID (200)](#mainfindpetbyid-get-petsid)