package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	hashmap "github.com/silasbue/exam-2021/gRPC"
	"google.golang.org/grpc"
)

func main() {
	file, _ := openLogFile("./logs/serverlog.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	if len(os.Args) != 2 {
		log.Printf("Please input a number to run the server on. Fx. inputting 3 would run the server on port 3003")
		return
	}

	ownId := os.Args[1]

	listen, _ := net.Listen("tcp", "localhost:300"+ownId)

	convOwnId, _ := strconv.ParseInt(ownId, 10, 32)

	grpcServer := grpc.NewServer()
	hashmap.RegisterHashmapServer(grpcServer, &Server{
		id:     int32(convOwnId),
		values: make(map[int32]int32),
	})

	log.Printf("server listening at %v", listen.Addr())

	grpcServer.Serve(listen)
}

func (s *Server) Put(ctx context.Context, req *hashmap.PutRequest) (*hashmap.PutReply, error) {
	s.values[req.Key] = req.Value
	log.Printf("server %v: recieved a bid for key %v. value: %v", s.id, req.GetKey(), req.GetValue())
	return &hashmap.PutReply{Success: true}, nil
}

func (s *Server) Get(ctx context.Context, req *hashmap.GetRequest) (*hashmap.GetReply, error) {
	value := s.values[req.Key]
	return &hashmap.GetReply{Value: value}, nil
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

type Server struct {
	hashmap.UnimplementedHashmapServer
	id     int32
	values map[int32]int32
}
