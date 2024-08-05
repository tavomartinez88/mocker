## Mocker
Mocker is a tool to create mocks for emulate third-party api rest's response  

### Setup
#### Install Go
> install go >= 1.20

#### Run tidy
> go mod tidy

### Run locally
> go run main.go

### Add endpoint
> curl -X POST http://{host}:{port}/add-endpoint -H "Content-Type: application/json" -d '{
>"method": "POST",
>"pattern": "/my-new-endpoint",
>"body": "{\"message\": \"Dynamic endpoint reached\"}", "status": 201
>}'

### Remove endpoint
> curl -X POST http://{host}:{port}/remove-endpoint -H "Content-Type: application/json" -d '{
>"pattern": "/my-new-endpoint",
>}'

> You can import these curls on postman
> base path is: [mocker](https://disciplinary-corinna-gusmartinez-dev-aba94555.koyeb.app/)
