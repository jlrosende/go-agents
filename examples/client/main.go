/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	pb "github.com/jlrosende/go-agents/proto/a2a/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.NewClient("unix:///tmp/agent_one.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewA2AServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	// defer cancel()
	r, err := c.SendMessage(context.Background(), &pb.SendMessageRequest{
		Request: &pb.Message{
			MessageId: uuid.NewString(),
			Role:      pb.Role_ROLE_USER,
			Content: []*pb.Part{
				{
					Part: &pb.Part_Text{
						Text: "READ tyhe file .gitignore",
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("could send message: %v", err)
	}

	log.Printf("MSG Response: %s", r.GetMsg())

}
