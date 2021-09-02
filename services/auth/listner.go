package main

import (
	"flag"
	"github.com/bhrg3se/flahmingo-homework/services/auth/pb/proto"
	"github.com/bhrg3se/flahmingo-homework/services/auth/server"
	"github.com/bhrg3se/flahmingo-homework/services/auth/store"
	"github.com/bhrg3se/flahmingo-homework/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func startServer() {
	path := flag.String("c", "/etc/flahmingo", "config file location")
	writeToFile := flag.Bool("f", false, "write logs to file")
	flag.Parse()

	config := utils.ParseConfig(*path)

	// set log level and file
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	//logrus.SetReportCaller(true)

	if *writeToFile {
		f, err := os.OpenFile("/var/log/boiler.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	}

	// initialise database and other dependencies (store)
	s := store.NewStore(config)

	// register and start a gRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, server.NewServer(s))

	listener, err := net.Listen("tcp", config.Server.Listen)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("starting gRPC server in %s", config.Server.Listen)
	grpcServer.Serve(listener)
}
