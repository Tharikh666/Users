package main

import (
	"User_Sample/common"
	User "User_Sample/domain/user"
	myUser "User_Sample/user"
	"log"
	"net"

	"google.golang.org/grpc"
)

func init() {

	common.LoadConfig()
}

func main() {

	log.Println("USER started ------>")

	lis, err := net.Listen("tcp", ":"+common.Config.ServerPort)
	log.Println("Server Connection: ", common.Config.ServerPort)

	if err != nil {

		log.Println("Error occured in opening listener", err)
	}

	mygRPCServer := grpc.NewServer()
	User.RegisterUsersServer(mygRPCServer, &myUser.UserService{})

	log.Println("Listening....")

	mygRPCServer.Serve(lis)
}
