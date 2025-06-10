package main

import "github.com/thatmatin/subserv/cmd"

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
