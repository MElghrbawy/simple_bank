package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	db "github.com/MElghrbawy/simple_bank/db/sqlc"
	_ "github.com/MElghrbawy/simple_bank/doc/statik"
	"github.com/MElghrbawy/simple_bank/gapi"
	"github.com/MElghrbawy/simple_bank/pb"
	"github.com/MElghrbawy/simple_bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	defer cancel()

	if err != nil {
		log.Fatal("cannot register gateway server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("doc/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik file system:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create http listener :", err)
	}

	log.Printf("start HTTP server on %s", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}

}
