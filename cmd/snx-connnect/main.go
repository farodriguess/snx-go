package main

import (
	"fmt"
	"snxgo/snx"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
	"golang.org/x/term"
)

var CLI = snx.SNXParams{}

func main() {
	ctx := kong.Parse(&CLI)
	fmt.Println(ctx.Command())
	getPasswordIfNecessary()
	if CLI.Debug {
		printCLIArgs()
	}

	snxConnect := snx.SNXConnect{Params: CLI}
	snxConnect.Connect()

}

func getPasswordIfNecessary() {

	if strings.TrimSpace(CLI.Password) == "" {
		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		CLI.Password = string(bytePassword)
		fmt.Println()
	}

}

func printCLIArgs() {
	fmt.Printf(`
CLI Args:
Host: %s
User: %s
Password: ***
Realm: %s
SkipSecurity: %v
`, CLI.Host, CLI.User, CLI.Realm, CLI.SkipSecurity)
}
