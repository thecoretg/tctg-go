module github.com/thecoretg/threatdown-site-list-lambda

go 1.26.2

require (
	github.com/aws/aws-lambda-go v1.54.0
	github.com/joho/godotenv v1.5.1
	github.com/thecoretg/tctg-go v0.0.0-20260609204303-c5c702e55b1f
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	resty.dev/v3 v3.0.0-rc.1 // indirect
)

replace github.com/thecoretg/tctg-go => ../..
