package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	hashmap "github.com/silasbue/exam-2021/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ownPort, _ := strconv.Atoi(os.Args[1])

	file, _ := openLogFile("./logs/client.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	servers := make([]hashmap.HashmapClient, 3)

	for i := 0; i < 3; i++ {
		port := int32(3000) + int32(i)

		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Front end %v: Could not connect: %s", ownPort, err)
		}
		servers[i] = hashmap.NewHashmapClient(conn)
		defer conn.Close()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.Split(scanner.Text(), " ")
		command[0] = strings.ToLower(command[0])
		command[1] = strings.ToLower(command[1])

		key, _ := strconv.Atoi(command[1])

		if command[0] == "put" {
			command[2] = strings.ToLower(command[2])
			value, _ := strconv.Atoi(command[2])

			for id, server := range servers {
				res, err := server.Put(ctx, &hashmap.PutRequest{Key: int32(key), Value: int32(value)})
				if err != nil {
					log.Printf("server %v: ERROR - %v", id, err)
					continue
				}

				log.Printf("put succeeded: %v", res.GetSuccess())
			}
		}

		if command[0] == "get" {

			for id, server := range servers {
				res, err := server.Get(ctx, &hashmap.GetRequest{Key: int32(key)})
				if err != nil {
					log.Printf("server %v: ERROR - %v", id, err)
					continue
				}
				log.Printf("get call to key: %v returned value: %v", key, res.GetValue())
			}
		}
	}
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
