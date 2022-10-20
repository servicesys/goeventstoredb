# goeventstoredb

EVENT STORE GOLANG


## migration
     - tern migrate 

## Test 
 #### test 
      - go test ./...
      - go clean -testcache
 #### coverage 
     - go test ./... -coverprofile=coverage.out
     - go tool cover -func=coverage.out 
     - go tool cover -html=coverage.out


## Docker 
- Stop all container :      docker kill $(docker ps -q)
- Remove all  containers:   docker rm $(docker ps -a -q)
- Remove allimages: docker rmi $(docker images -q)
