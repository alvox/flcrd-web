FROM alpine:3.10
ADD app app
CMD ./app -port=$PORT -dsn=$DSN -auth_domain=$AUTH_DOMAIN -mail_api_url=$MAIL_API_URL -mail_api_key=$MAIL_API_KEY