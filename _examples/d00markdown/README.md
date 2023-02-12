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
| `GET /pets` | [main.FindPets](#mainfindpets-get-pets)  | | Returns all pets |
| `POST /pets` | [main.AddPet](#mainaddpet-post-pets)  | | Creates a new pet |
| `DELETE /pets/{id}` | [main.DeletePet](#maindeletepet-delete-petsid)  | | Deletes a pet by ID |
| `GET /pets/{id}` | [main.FindPetByID](#mainfindpetbyid-get-petsid)  | | Returns a pet by ID |


### main.FindPets `GET /pets`

Returns all pets

| name | value | 
| --- | --- |
| operationId | main.FindPets |
| endpoint | `GET /pets` |

#### input

```go
// GET /pets

type Input struct {
}
```

#### output (application/json)

```go
// GET /pets default
type OutputDefault struct { // Error
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

#### input

```go
// POST /pets

type Input struct {
}
```

#### output (application/json)

```go
// POST /pets default
type OutputDefault struct { // Error
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

#### input

```go
// DELETE /pets/{id}

type Input struct {
}
```

#### output (application/json)

```go
// DELETE /pets/{id} default
type OutputDefault struct { // Error
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

#### input

```go
// GET /pets/{id}

type Input struct {
}
```

#### output (application/json)

```go
// GET /pets/{id} default
type OutputDefault struct { // Error
}
```

#### description

Returns a pet based on a single ID



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [Pet](#pet) |  |


### Pet



```go
type Pet struct {
	// Unique id of the pet
	id integer `format:"int64"`
	// Name of the pet
	name string
	// Type of the pet
	tag? string
}

```
