# Seaports - Client API and domain service

## Preconditions:
In order to be able lint, build, run and test the applications, some tools are required:
go, protoc, protoc-gen-go, golangci-lint, docker

## 1. Build the client API and service as docker images
```
make images
```
## 2. Run client API and service as docker containers
```
make start-containers
```
> Note 1: Building images and running them can be done using command `make start-build-containers`
> Note 2: To stop the containers run: `make stop-containers`

## 3. Run simple end-to-end test 
```
make e2e-test
```
This test will:
 - access the root path of the server for the welcome message (test http server is alive)
 - import ports from the ./assignement command using curl: 
 - get a port by it's ID
 - import ports again to check if any changed ports get updated
 - check updated port

> Run `make clean` to clean up binaries


## TODO - improvements
- Add comments to all exported methods, types and type fields for generating documentation directly from code using `godoc`
- Provide a error count or limit the error array within the response for import ports insted of the error array. Using the array is a bad idea if the number of errors is big, because it will increase the needed memory in order to accumulate them. Also the response size can increase drastically and the end-user will have to wait longer to get the response.
- Using a production logger like `zap.Logger` which I use in my daily job. Reason for using it can be found on https://github.com/uber-go/zap. It is fast, consumes less memory than other, it's robust and maintained on github. 
- Adding unit-tests
- Probably using Mongo would be better
