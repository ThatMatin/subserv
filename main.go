package main

import "github.com/thatmatin/subserv/cmd"

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
