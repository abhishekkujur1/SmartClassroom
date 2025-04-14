package main

import (
	"github.com/abhishekkujur1/SmartClassroom/server"
	"github.com/abhishekkujur1/SmartClassroom/cmd"
)

func main() {
	server.StartServer()
	cmd.StartCapture()
}