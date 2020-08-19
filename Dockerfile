FROM alpine:3.10
ADD app app
CMD ./app -port=$PORT -dsn=$DSN -appkey=$APPKEY -mail_api_url=$MAIL_API_URL -mail_api_key=$MAIL_API_KEY -aws_id=$AWS_ID -aws_secret=$AWS_SECRET