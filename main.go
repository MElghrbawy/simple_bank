package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/MElghrbawy/simple_bank/api"
	db "github.com/MElghrbawy/simple_bank/db/sqlc"
	"github.com/MElghrbawy/simple_bank/gapi"
	"github.com/MElghrbawy/simple_bank/pb"
	"github.com/MElghrbawy/simple_bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)
	// runGinServer(config, store)

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener :", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}

}
