## Backend for the [flashcards.rocks](https://flashcards.rocks) site.

🌎 https://flashcards.rocks  ![](https://github.com/alvox/flcrd-web/workflows/Build%20image/badge.svg)

### Build app
`./build.sh`

### Run tests
`go test ./...`

### Build test database
`cd db && ./build.sh`

### Run
`cd env && docker-compose up -d`

### Launch parameters
- port Default: 5000
- dsn Postgresql datasource. Default: postgres://flcrd:flcrd@flcrd-test-db/flcrd?sslmode=disable
- appkey Secret used to sign JWT tokens
- mail_api_url Mailgun url
- mail_api_key Mailgun token