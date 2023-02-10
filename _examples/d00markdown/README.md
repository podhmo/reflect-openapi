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
| --- | --- | --- | --- | --- |
| `GET /pets` | [main.FindPets](#findpets--get-pets)  | | Returns all pets |
| `POST /pets` | [main.AddPet](#findpets--get-pets)  | | Creates a new pet |
| `DELETE /pets/{id}` | [main.DeletePet](#findpets--get-pets)  | | Deletes a pet by ID |
| `GET /pets/{id}` | [main.FindPetByID](#findpets--get-pets)  | | Returns a pet by ID |

## schemas