# go-clean-arch-grpc

This is an example of implementation of Clean Architecture in Go (Golang) projects. With GRPC


The explanation about this project's structure  can read from this medium's post :  https://medium.com/@imantumorang/implementing-grpc-service-in-golang-afb9e05c0064


The client project can seen here : https://github.com/bxcodec/sample-client-grpc

### How To Run This Project

#### Dowload the project
```bash
# Download the project 
go get github.com/bxcodec/go-clean-arch-grpc

#move to directory
cd $GOPATH/src/github.com/bxcodec/go-clean-arch-grpc
 
# Install Dependencies
glide install -v

# Make File
make
```

#### Set  the config
Open `config.json`
Change to your own database config
```js
{
  "debug": true,
  "server": {
    "address": ":8080"
  },
  "database": {
      "host": "localhost",
      "port": "33061",
      "user": "root",
      "pass": "password",
      "name": "article"
  }

}

```


####  Run Project

```bash
go run main.go

```
 

> Make Sure you have run the article.sql in your mysql
