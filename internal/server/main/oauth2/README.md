# Oauth2

List of Oauth2 http apis

## [GET] **/v1/login/oauth/authorize**

### Request

Query:

- client_id: `user client_id, used for client_credentials grant`
- redirect_uri: `user authorization callback URL`
- state: `random string to protect againts forgery attacks`
- login: `specific account to use for signing in`

### Example

Header:

- grant: `client_credentials`

`/v1/oauth2/auth?client_id=xxx&client_secret=xxxx`

Header:

- grant: `password`

`/v1/oauth2/auth?username=xxxx&password=xxx`

## [POST] **/v1/login/oauth/authenticate

Header:

- grant_type: `password|client_credentials`
- otp: `required`