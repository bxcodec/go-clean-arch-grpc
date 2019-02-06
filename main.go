package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"

	"google.golang.org/grpc"

	deliveryGrpc "github.com/bxcodec/go-clean-arch-grpc/article/delivery/grpc"
	articleRepo "github.com/bxcodec/go-clean-arch-grpc/article/repository"
	articleUcase "github.com/bxcodec/go-clean-arch-grpc/article/usecase"
	cfg "github.com/bxcodec/go-clean-arch-grpc/config/env"
	_ "github.com/go-sql-driver/mysql"
)

var config cfg.Config

func init() {
	config = cfg.NewViperConfig()

	if config.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	dbHost := config.GetString(`database.host`)
	dbPort := config.GetString(`database.port`)
	dbUser := config.GetString(`database.user`)
	dbPass := config.GetString(`database.pass`)
	dbName := config.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && config.GetBool("debug") {
		fmt.Println(err)
	}
	defer dbConn.Close()

	ar := articleRepo.NewMysqlArticleRepository(dbConn)
	au := articleUcase.NewArticleUsecase(ar)
	list, err := net.Listen("tcp", config.GetString("server.address"))
	if err != nil {
		fmt.Println("SOMETHING HAPPEN")
	}

	server := grpc.NewServer()
	deliveryGrpc.NewArticleServerGrpc(server, au)
	fmt.Println("Server Run at ", config.GetString("server.address"))

	err = server.Serve(list)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}

}
