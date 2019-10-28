# User

List of user APIs

## [GET] /v1/user/login

### Request

url value:

- username: `my_username`
- type: `otp|password`
- metadata: `redirect_url=https://redirect_url,foo=bar`

example: `/v1/login?username=my_username&type=otp&metadata=redirect_url=https://redirect_url,foo=bar`

### Response

```json
{
    "status": "OK",
    "data": {
        "authentication_url": "/v1/login",
        "authentication_method": "post",
        "authentication_state": "thisismystate"
    }
}
```

## [POST] /v1/user/login

### Request

params:

- password: `my_password`
- state: `thisismystate`

### Response

```json
{
    "status": "OK",
    "data": {
        "session": "my_session",
        "user_id": "my_user_id",
        "metadata": {
            "redirect_url": "https://redirect_url"
        }
    }
}
```

## [POST] /v1/user/register

### Request

header:

- sid : `session_id`

params:

- username: `my_phone_number`
- full_name: `my_full_name`
- email: `my_email`

## Response

```json
{
    "status": "OK",
    "data": {
        "authentication_state": "your_state",
        "authentication_id": "xxxx"
    }
}
```

## [POST] /v1/user/register/confirm

### Request

header:

- sid: `session_id`

params:

- password: `my_password`
- authentication_id: `authentication_id`
- state: `my_state`

### Response

```json
{
    "status": "OK",
    "data": {
        "user_id": "my_user_id",
        "session": "my_session"
    }
}
```