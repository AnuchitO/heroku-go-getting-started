# Ariskill Backend

## Getting Started

1. start mongodb `make db-up`
1. seed data to database `make seed`
1. start backend api server `make run`
1. verify backend api server is up and running `make health`

  if you see following response, it means backend api server is up and running

```json
{"status":"success","message":"","data":"ariskill is ready and connected to database"}
```


# Set environment variables

- check the value of.env
- set following enviroment variable

  - OCI_CLIENT
  - MONGODB_USERNAME
  - MONGODB_PASSWORD
  - GOOGLE_OIDC_CLIENT_ID
    > for DEV mode, you could get google OIDC client id from frontend next.config
  - GOOGLE_OIDC_CLIENT_SECRET ?
  - GOOGLE_OIDC_REDIRECT_URI ?
    > these two variables required to define in but not used in DEV mode


# Makefile commands

- **seed**: insert data from local to database, support -_env_ prefix
- **run**: run main service only, support -_env_ prefix
- **test**: run all test
  - **test-integration** run test with deployed container
- **test-cover**: run unit test coverage

- **swagger-install** : install swag
- **set gopath** : set gopath incase cannot use command swag after install
- **swagger** : generate swagger file and format comment of godoc
