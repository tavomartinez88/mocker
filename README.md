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
#### Example
Create on postman a request using:
[add-endpoint](https://disciplinary-corinna-gusmartinez-dev-aba94555.koyeb.app/add-endpoint)

And you can send a body (if you want it) as:
```sh
{
  "method": "POST",
  "pattern": "/my-new-endpoint",
  "body": {
      "property1": "https://example.com",
      "property2": {
          "property2.1": 1,
          "property2.2": "12345",
          "property2.3": "67890"
      }
   },
"status": 201
}
```

### Remove endpoint
#### Example
Create on postman a request using:
[remove-endpoint](https://disciplinary-corinna-gusmartinez-dev-aba94555.koyeb.app/remove-endpoint)

And you should send a body with endpoint do you want delete:
```sh
{
  "pattern": "/my-new-endpoint"
}
```

