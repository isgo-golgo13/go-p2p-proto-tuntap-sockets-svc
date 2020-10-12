 # Go parameters
    GOCMD=go
    GOBUILD=go build
    GOCLEAN=go clean
    GOTEST=go test
    SERVER=p2ps_svc
  
    all: build
    build: 
		$ env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(SERVER) -v
    test: 
		$(GOTEST) -v ./...
    clean: 
		$(GOCLEAN)
			rm -f $(SERVER)
    run:
		$ env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(SERVER) -v ./...
		sudo ./$(SERVER)
