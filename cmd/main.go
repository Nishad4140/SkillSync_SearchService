package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/Nishad4140/SkillSync_SearchService/initializer"
	"github.com/Nishad4140/SkillSync_SearchService/db"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	mongokey := os.Getenv("MONGO_KEY")
	mongoDB, err := db.InitMongoDB(mongokey)
	if err != nil {
		log.Fatal("error connecting to mongodb database")
	}

	listener, err := net.Listen("tcp", ":4003")
	if err != nil {
		log.Fatal("failed to listen on port 4003")
	}
	fmt.Println("search service listening on port 4003")

	services := initializer.Initializer(mongoDB)
	server := grpc.NewServer()

	pb.RegisterSearchServiceServer(server, services)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to listen on port 4003")
	}
}
