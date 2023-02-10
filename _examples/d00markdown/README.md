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
| `GET /pets` | [main.FindPets](#findpets--get-pets)  | | Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien. |
| `POST /pets` | [main.AddPet](#findpets--get-pets)  | | Creates a new pet |
| `DELETE /pets/{id}` | [main.DeletePet](#findpets--get-pets)  | | Deletes a pet by ID |
| `GET /pets/{id}` | [main.FindPetByID](#findpets--get-pets)  | | Returns a pet by ID |


### main.FindPets `GET /pets`

Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.

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

Creates a new pet

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

Deletes a pet by ID

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

Returns a pet by ID

Returns a pet based on a single ID

## schemas