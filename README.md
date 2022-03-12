# Gank

My Personal Golang Dev-Kit. This package is intended to wrap many common functions to be ready to use without much hassle. All services are wrapped inside an interface so you can replace it anytime. This package also using options pattern to maximize customization yet also maintain the simplicity. To handle application lifecycle (start & stop) and also dependencies injection, we are using [fx](https://github.com/uber-go/fx) & [dig](https://github.com/uber-go/dig) and wrap it under interface called lifecycler.

Included internal services are : logger, http-server, http-router, middlewares (auth, http-response, logger), env, config, secret, hash, token, etc. Included external services are : postgres-db.

## How to Install
`go get github.com/gndw/gank`

### :white_check_mark: How to Create Empty Application

```go
package main

import (
	"log"

	"github.com/gndw/gank"
)

func main() {

	// main function to create your application
	err := gank.CreateAndRunApp(
  
		// lifecycler object to handle start, stop, and dependencies injection (must use)
		gank.DefaultLifecycler(),
    
		// inject all default services such as : logger, http server, router, middlewares, token, etc
		// those services will not be called if not used
		gank.WithDefaultInternalProviders(),
	)
	if err != nil {
		log.Fatal(err)
	}

}
```
- create your `main.go` with code above
- run app using `go run main.go`
```
$ go run main.go 
INFO[2022-03-09T00:39:33+07:00] application startup in 410.592µs
```
- application will start but doing nothing since we are not initializing anything
- close application using SIGTERM or cmd+c

<br />

### :white_check_mark: How to Create Application with HTTP Health Check
```go
err := gank.CreateAndRunApp(
  // ... previous code
  
  // use this option to start http health service
  gank.WithHealthHandler(),
)
```
- this option `gank.WithHealthHandler()` will startup http server and adding health check system under endpoint `/health`
- run your app & do health check
```
$ go run main.go
INFO[2022-03-09T00:48:03+07:00] starting application with env: development   
INFO[2022-03-09T00:48:03+07:00] s> starting http server ...                  
INFO[2022-03-09T00:48:03+07:00] d> starting http server in 194.191µs         
INFO[2022-03-09T00:48:03+07:00] http server is listening at port 9000        
INFO[2022-03-09T00:48:03+07:00] application startup in 2.787555ms
```
```
$ curl -i http://localhost:9000/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 08 Mar 2022 17:52:35 GMT
Content-Length: 79

{"data":{"services":[{"service_name":"app","is_healthy":true,"status":"OK"}]}}
```


### :white_check_mark: How to Create your Custom Endpoint
```go
err := gank.CreateAndRunApp(
  // ... previous code
  
  // add new http endpoint
  gank.WithInvokers(
	func(server server.Service) (err error) {
	
		// adding OK endpoint
		err = server.AddHttpHandler(model.AddHTTPRequest{
			Method:   constant.HTTPMethodGet,
			Endpoint: "/my-custom-endpoint/ok",
			Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
				return "OK", nil
			},
		})
		if err != nil {
			return err
		}

		// adding Bad Request endpoint
		err = server.AddHttpHandler(model.AddHTTPRequest{
			Method:   constant.HTTPMethodGet,
			Endpoint: "/my-custom-endpoint/bad",
			Handler: func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {
				return nil, errors.New("bad request response")
			},
		})
		if err != nil {
			return err
		}

		return nil
	},
  ),
)
```
- this option `gank.WithInvokers()` will start any function passed to it. App will automatically pass any parameter needed in those functions as long as those services are injected before. In this case `server.Service` is used so that user can add custom http endpoint.
- run your app & check your custom endpoint
```
$ curl -i localhost:9000/my-custom-endpoint/ok
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Mar 2022 11:43:50 GMT
Content-Length: 14

{"data":"OK"}
```
```
$ curl -i localhost:9000/my-custom-endpoint/bad
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Mar 2022 11:43:56 GMT
Content-Length: 36

{"errors":["bad request response"]}
```

## Services
TODO:explain
- logger
- http-server
- http-router
- middlewares
- env
- config
- secret
- hash
- token
- flag
- marshal
- db

## Other
TODO:explain
- errorsg (custom error)
