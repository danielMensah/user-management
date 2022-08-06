# User Management

This is a small microservice to manage access to Users.

## Requirements

These dependencies can also be installed with [Homebrew](https://brew.sh/).

* Requires Go 1.18 or greater. This can be installed with brew `brew install go` or
  downloaded [here](https://golang.org/doc/install).
* Requires Golangci Lint. This can be installed with brew `brew install golangci-lint` or
  downloaded [here](https://golangci-lint.run/usage/install/#local-installation).
* Requires Docker and Docker Compose. This can be installed with brew `brew install docker` or
  downloaded [here](https://docs.docker.com/engine/install/).

### Install Dependencies

Install dependencies, issue the following command(s):

```bash
make install
```

### Testing and Formatting

To run the tests, issue the following command(s):

```bash
make test
```

#### Lint only

Run linting only:

```bash
make lint
```

## How to Run

To run the application, simply issue the following example command(s):

```bash
make local
```

## Environment Variables

Environment variables needed to start the application are:

| Variable                          | Description                                                    | Required | Example/Default       |
|-----------------------------------|----------------------------------------------------------------|----------|-----------------------|
| API_HOST                          | Host that the exposed api endpoints should be run on           | :x:      | 0.0.0.0               |
| API_PORT                          | Port that the exposed api endpoints will listen on for request | :x:      | 8000                  |
| API_MONGO_URI                     | Mongo instance URI                                             | &check;  | mongodb://mongo:27017 |                                                                                                                    | &check;  | E.G: us-east-1                         |
| API_MONGO_DB_NAME                 | Mongo Database Name to initialize                              | &check;  | usermanagement        |                                                                                                                    | &check;  | E.G: us-east-1                         |

## API Endpoints

Base URLs:

* http://localhost:8000/api/v1

### Health check

`GET /_healthz`

Returns 200 if the service is up and running

> Example responses

> 200 Response

```
"OK"
```

<h3 id="get___healthz-responses">Responses</h3>

| Status | Meaning                                                 | Description               | Schema |
|--------|---------------------------------------------------------|---------------------------|--------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | Service is up and running | string |


### getUsers

`GET /users`

<h3 id="getusers-parameters">Parameters</h3>

| Name    | In    | Type           | Required | Description              |
|---------|-------|----------------|----------|--------------------------|
| country | query | string         | false    | User country             |
| email   | query | string         | false    | User email               |
| page    | query | integer(int64) | true     | Page number              |
| limit   | query | integer(int64) | true     | Number of users per page |

> Example responses

> 200 Response

```json
{
  "users": [
    {
      "_id": "string",
      "first_name": "John",
      "last_name": "Doe",
      "nickname": "jd",
      "email": "js@example.com",
      "country": "UK",
      "created_at": "2019-08-24T14:15:22Z",
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ]
}
```

<h3 id="getusers-responses">Responses</h3>

| Status | Meaning                                                                    | Description           | Schema                                      |
|--------|----------------------------------------------------------------------------|-----------------------|---------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)                    | A list of users       | [GetUsersResponse](#schemagetusersresponse) |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)           | Invalid request       | [Error](#schemaerror)                       |
| 500    | [Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1) | Internal server error | [Error](#schemaerror)                       |

<aside class="success">
This operation does not require authentication
</aside>

## createUser

`POST /users`

> Body parameter

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "password": "worm",
  "country": "UK"
}
```

<h3 id="createuser-parameters">Parameters</h3>

| Name | In   | Type                                    | Required | Description |
|------|------|-----------------------------------------|----------|-------------|
| body | body | [UserCreateData](#schemausercreatedata) | false    | none        |

> Example responses

> 201 Response

```json
{
  "_id": "string"
}
```

<h3 id="createuser-responses">Responses</h3>

| Status | Meaning                                                                    | Description           | Schema                                          |
|--------|----------------------------------------------------------------------------|-----------------------|-------------------------------------------------|
| 201    | [Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)               | Created user          | [CreateUserResponse](#schemacreateuserresponse) |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)           | Invalid request       | [Error](#schemaerror)                           |
| 500    | [Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1) | Internal server error | [Error](#schemaerror)                           |

## deleteUser

<a id="opIddeleteUser"></a>

`DELETE /users/{id}`

<h3 id="deleteuser-parameters">Parameters</h3>

| Name | In   | Type   | Required | Description |
|------|------|--------|----------|-------------|
| id   | path | string | true     | User ID     |

> Example responses

> 400 Response

```json
{
  "message": "string"
}
```

<h3 id="deleteuser-responses">Responses</h3>

| Status | Meaning                                                                    | Description           | Schema                |
|--------|----------------------------------------------------------------------------|-----------------------|-----------------------|
| 204    | [No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)            | Deleted user          | None                  |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)           | Invalid request       | [Error](#schemaerror) |
| 500    | [Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1) | Internal server error | [Error](#schemaerror) |

## updateUser

<a id="opIdupdateUser"></a>

`PUT /users/{id}`

> Body parameter

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "password": "worm",
  "country": "UK"
}
```

<h3 id="updateuser-parameters">Parameters</h3>

| Name | In   | Type                                    | Required | Description |
|------|------|-----------------------------------------|----------|-------------|
| id   | path | string                                  | true     | User ID     |
| body | body | [UserUpdateData](#schemauserupdatedata) | false    | none        |

> Example responses

> 200 Response

```json
{
  "_id": "string",
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "country": "UK",
  "created_at": "2019-08-24T14:15:22Z",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

<h3 id="updateuser-responses">Responses</h3>

| Status | Meaning                                                                    | Description           | Schema                |
|--------|----------------------------------------------------------------------------|-----------------------|-----------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)                    | Updated user          | [User](#schemauser)   |
| 400    | [Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)           | Invalid request       | [Error](#schemaerror) |
| 500    | [Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1) | Internal server error | [Error](#schemaerror) |

# Schemas

<h2 id="tocS_GetUsersResponse">GetUsersResponse</h2>
<!-- backwards compatibility -->
<a id="schemagetusersresponse"></a>
<a id="schema_GetUsersResponse"></a>
<a id="tocSgetusersresponse"></a>
<a id="tocsgetusersresponse"></a>

```json
{
  "users": [
    {
      "_id": "string",
      "first_name": "John",
      "last_name": "Doe",
      "nickname": "jd",
      "email": "js@example.com",
      "country": "UK",
      "created_at": "2019-08-24T14:15:22Z",
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ]
}

```

### Properties

| Name  | Type                  | Required | Restrictions | Description |
|-------|-----------------------|----------|--------------|-------------|
| users | [[User](#schemauser)] | false    | none         | none        |

<h2 id="tocS_CreateUserResponse">CreateUserResponse</h2>
<!-- backwards compatibility -->
<a id="schemacreateuserresponse"></a>
<a id="schema_CreateUserResponse"></a>
<a id="tocScreateuserresponse"></a>
<a id="tocscreateuserresponse"></a>

```json
{
  "_id": "string"
}

```

### Properties

| Name | Type            | Required | Restrictions | Description |
|------|-----------------|----------|--------------|-------------|
| _id  | [Id](#schemaid) | true     | none         | none        |

<h2 id="tocS_User">User</h2>
<!-- backwards compatibility -->
<a id="schemauser"></a>
<a id="schema_User"></a>
<a id="tocSuser"></a>
<a id="tocsuser"></a>

```json
{
  "_id": "string",
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "country": "UK",
  "created_at": "2019-08-24T14:15:22Z",
  "updated_at": "2019-08-24T14:15:22Z"
}

```

### Properties

| Name       | Type                          | Required | Restrictions | Description |
|------------|-------------------------------|----------|--------------|-------------|
| _id        | [Id](#schemaid)               | true     | none         | none        |
| first_name | [FirstName](#schemafirstname) | true     | none         | none        |
| last_name  | [LastName](#schemalastname)   | true     | none         | none        |
| nickname   | [Nickname](#schemanickname)   | true     | none         | none        |
| email      | [Email](#schemaemail)         | true     | none         | none        |
| country    | [Country](#schemacountry)     | true     | none         | none        |
| created_at | [CreatedAt](#schemacreatedat) | true     | none         | none        |
| updated_at | [UpdatedAt](#schemaupdatedat) | true     | none         | none        |

<h2 id="tocS_UserUpdateData">UserUpdateData</h2>
<!-- backwards compatibility -->
<a id="schemauserupdatedata"></a>
<a id="schema_UserUpdateData"></a>
<a id="tocSuserupdatedata"></a>
<a id="tocsuserupdatedata"></a>

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "password": "worm",
  "country": "UK"
}

```

### Properties

| Name       | Type                          | Required | Restrictions | Description |
|------------|-------------------------------|----------|--------------|-------------|
| first_name | [FirstName](#schemafirstname) | false    | none         | none        |
| last_name  | [LastName](#schemalastname)   | false    | none         | none        |
| nickname   | [Nickname](#schemanickname)   | false    | none         | none        |
| email      | [Email](#schemaemail)         | false    | none         | none        |
| password   | [Password](#schemapassword)   | false    | none         | none        |
| country    | [Country](#schemacountry)     | false    | none         | none        |

<h2 id="tocS_UserCreateData">UserCreateData</h2>
<!-- backwards compatibility -->
<a id="schemausercreatedata"></a>
<a id="schema_UserCreateData"></a>
<a id="tocSusercreatedata"></a>
<a id="tocsusercreatedata"></a>

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "jd",
  "email": "js@example.com",
  "password": "worm",
  "country": "UK"
}

```

### Properties

| Name       | Type                          | Required | Restrictions | Description |
|------------|-------------------------------|----------|--------------|-------------|
| first_name | [FirstName](#schemafirstname) | true     | none         | none        |
| last_name  | [LastName](#schemalastname)   | true     | none         | none        |
| nickname   | [Nickname](#schemanickname)   | true     | none         | none        |
| email      | [Email](#schemaemail)         | true     | none         | none        |
| password   | [Password](#schemapassword)   | true     | none         | none        |
| country    | [Country](#schemacountry)     | true     | none         | none        |

<h2 id="tocS_Error">Error</h2>
<!-- backwards compatibility -->
<a id="schemaerror"></a>
<a id="schema_Error"></a>
<a id="tocSerror"></a>
<a id="tocserror"></a>

```json
{
  "message": "string"
}

```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| message | string | true     | none         | none        |

<h2 id="tocS_Id">Id</h2>
<!-- backwards compatibility -->
<a id="schemaid"></a>
<a id="schema_Id"></a>
<a id="tocSid"></a>
<a id="tocsid"></a>

```json
"string"

```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| *anonymous* | string | false    | none         | none        |

<h2 id="tocS_FirstName">FirstName</h2>
<!-- backwards compatibility -->
<a id="schemafirstname"></a>
<a id="schema_FirstName"></a>
<a id="tocSfirstname"></a>
<a id="tocsfirstname"></a>

```json
"John"

```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| *anonymous* | string | false    | none         | none        |

<h2 id="tocS_LastName">LastName</h2>
<!-- backwards compatibility -->
<a id="schemalastname"></a>
<a id="schema_LastName"></a>
<a id="tocSlastname"></a>
<a id="tocslastname"></a>

```json
"Doe"

```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| *anonymous* | string | false    | none         | none        |

<h2 id="tocS_Nickname">Nickname</h2>
<!-- backwards compatibility -->
<a id="schemanickname"></a>
<a id="schema_Nickname"></a>
<a id="tocSnickname"></a>
<a id="tocsnickname"></a>

```json
"jd"

```

### Properties

| Name         | Type    | Required  | Restrictions  | Description  |
|--------------|---------|-----------|---------------|--------------|
| *anonymous*  | string  | false     | none          | none         |

<h2 id="tocS_Email">Email</h2>
<!-- backwards compatibility -->
<a id="schemaemail"></a>
<a id="schema_Email"></a>
<a id="tocSemail"></a>
<a id="tocsemail"></a>

```json
"js@example.com"

```

### Properties

| Name         | Type    | Required  | Restrictions  | Description  |
|--------------|---------|-----------|---------------|--------------|
| *anonymous*  | string  | false     | none          | none         |

<h2 id="tocS_Password">Password</h2>
<!-- backwards compatibility -->
<a id="schemapassword"></a>
<a id="schema_Password"></a>
<a id="tocSpassword"></a>
<a id="tocspassword"></a>

```json
"worm"

```

### Properties

| Name         | Type    | Required  | Restrictions  | Description |
|--------------|---------|-----------|---------------|-------------|
| *anonymous*  | string  | false     | none          | none        |

<h2 id="tocS_Country">Country</h2>
<!-- backwards compatibility -->
<a id="schemacountry"></a>
<a id="schema_Country"></a>
<a id="tocScountry"></a>
<a id="tocscountry"></a>

```json
"UK"

```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| *anonymous* | string | false    | none         | none        |

<h2 id="tocS_CreatedAt">CreatedAt</h2>
<!-- backwards compatibility -->
<a id="schemacreatedat"></a>
<a id="schema_CreatedAt"></a>
<a id="tocScreatedat"></a>
<a id="tocscreatedat"></a>

```json
"2019-08-24T14:15:22Z"

```

### Properties

| Name        | Type              | Required | Restrictions | Description |
|-------------|-------------------|----------|--------------|-------------|
| *anonymous* | string(date-time) | false    | none         | none        |

<h2 id="tocS_UpdatedAt">UpdatedAt</h2>
<!-- backwards compatibility -->
<a id="schemaupdatedat"></a>
<a id="schema_UpdatedAt"></a>
<a id="tocSupdatedat"></a>
<a id="tocsupdatedat"></a>

```json
"2019-08-24T14:15:22Z"

```

### Properties

| Name        | Type              | Required | Restrictions | Description |
|-------------|-------------------|----------|--------------|-------------|
| *anonymous* | string(date-time) | false    | none         | none        |
