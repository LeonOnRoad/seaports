# Seaports - Client API and domain service

## Preconditions:
In order to be able lint, build, run and test the applications, some tools are required: `go, protoc, protoc-gen-go, golangci-lint, docker`. So, make sure to install them before trying to run the make commands.

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

## Dev corner
Make sure to check other make commands available useful for development:
- `proto-gen` - generates golang code from the proto files provided in `./proto` directory
- `build` - build services from the golang code and place binaries in `./bin` directory
- `images` - build docker images for both services
- `stop-containers` - perform `docker down`, that will stop all containers which are defined in `docker-compose.yaml`
- `start-containers` - perform `docker up -d`, that will start all containers which are defined in `docker-compose.yaml`
- `start-build-containers` - build docker images for both service, stop current running containers and start the newly built ones
- `e2e-test` - perform a very simple test
- `clean` - delete `.bin` directory
- `clean-proto` - delete files generated from proto files, which are now available under `proto/src`
- `lint` - perform a lint check on the `go` source code
  
> Note: To see logs from running containers run `docker-compose logs `

## TODO - improvements
- Add comments to all exported methods, types and type fields for generating documentation directly from code using `godoc`
- Provide a error count or limit the error array within the response for import ports insted of the error array. Using the array is a bad idea if the number of errors is big, because it will increase the needed memory in order to accumulate them. Also the response size can increase drastically and the end-user will have to wait longer to get the response.
- Using a production logger like `zap.Logger` which I use in my daily job. Reason for using it can be found on https://github.com/uber-go/zap. It is fast, consumes less memory than other, it's robust and maintained on github. 
- Adding unit-tests
- Probably using Mongo would be better
