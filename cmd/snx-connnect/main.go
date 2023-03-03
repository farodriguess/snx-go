package main

import (
	"fmt"
	"os"
	"runtime"
	"snxgo/snx"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
	"golang.org/x/term"
)

var (
	buildTime string
	version   string
)

var CLI struct {
	Host         string      `help:"VPN Hostname" name:"host" type:"string" required:""`
	User         string      `help:"VPN Username" name:"user" type:"string" required:""`
	Password     string      `help:"User's password" name:"password" type:"string"`
	Realm        string      `help:"VPN Realmd" name:"realm" type:"string" required:""`
	SkipSecurity bool        `help:"Skip TLS Verify in HTTPS Connection" name:"skip-security" type:"bool"`
	Debug        bool        `help:"Enable debug log" name:"debug" type:"bool"`
	Version      versionFlag `help:"Show build version" name:"version" type:"bool"`
}

func main() {

	ctx := kong.Parse(&CLI)
	fmt.Println(ctx.Command())

	getPasswordIfNecessary()

	if CLI.Debug {
		printCLIArgs()
	}

	snxConnect := snx.SNXConnect{Params: snx.SNXParams{
		Host:         CLI.Host,
		User:         CLI.User,
		Password:     CLI.Password,
		Realm:        CLI.Realm,
		SkipSecurity: CLI.SkipSecurity,
		Debug:        CLI.Debug,
	}}
	snxConnect.Connect()

}

type versionFlag bool

func (v versionFlag) BeforeApply() error {
	fmt.Printf("Version:\t%s\n", version)
	fmt.Printf("Build time:\t%s\n", buildTime)
	fmt.Printf("OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
	return nil
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
