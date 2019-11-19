# Authentication

List of Authentication APIs

## **[GET]** /v1/authenticate

### Request

header:

- sid: `my_session_id`
- uid: `my_user_id`

url value:

- type: `otp|pin`
- username: `your_username` | username is a `phone_number`
- state: `your_state`
- action: `register|pay`

below url value only valid if authentication `type` is `otp`

- country: `ID`
- code_length: `6`
- code_expiry_time_seconds: `150`

example: `/v1/authenticate?type=otp&username=628XXXXX&country=ID&state=thisismystate`

### Response

```json
{
    "status": "OK",
    "data": {
        "authentication_id": "xxxxxxxx"
    }
}
```

## **[POST]** /v1/authenticate

### Request

params:

- state: `your_state`
- authentication_id: `authentication_id`
- username: `your_username` [only required if authentication type is password]
- password: `your_password`

### Response

```json
{
    "status": "OK",
    "data": {
        "authentication_id": "xxxxxxxx",
        "authentication_status": "authenticated"
    }
}
```