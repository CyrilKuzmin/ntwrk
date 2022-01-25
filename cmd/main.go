// Command ntwrk is a tool for testing network performance.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	ntwrk "github.com/CyrilKuzmin/ntwrk/pkg"
	"github.com/waits/update"
)

const updateUrl = "https://api.github.com/repos/CyrilKuzmin/ntwrk/releases/latest"

var tag string
var version = update.ParseVersion(tag)

func main() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	} else {
		cmd = "help"
	}

	clientFlags := flag.NewFlagSet("client", flag.ExitOnError)
	host := clientFlags.String("host", "ntwrk.waits.io", "server to test against")
	port := 1600

	switch cmd {
	case "help":
		help()
	case "ip":
		clientFlags.Parse(os.Args[2:])
		client := ntwrk.NewClient(*host, port)
		client.Whoami()
	case "server":
		srv := ntwrk.NewServer(port, log.Default())
		srv.Start()
	case "run":
		clientFlags.Parse(os.Args[2:])
		client := ntwrk.NewClient(*host, port)
		client.StartCLI()
	case "update":
		update.Auto(version, updateUrl, update.CheckGithub)
	case "version":
		fmt.Printf("ntwrk %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	case "lib":
		clientFlags.Parse(os.Args[2:])
		client := ntwrk.NewClient(*host, port)
		client.SetMeasureDuration(5 * time.Second)
		ds, us, err := client.Measure()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("{\n    \"download_speed\": \"%v\",\n    \"upload_speed\": \"%v\"\n}\n", ds, us)
	default:
		fmt.Printf("Unknown command '%v'.\n", cmd)
	}
}

func help() {
	cmds := []string{"help", "ip\t", "run\t", "server", "update", "version"}
	descriptions := []string{
		"Show this help message",
		"Print external IP address",
		"Run performance tests",
		"Start a test server",
		"Checks for and downloads an updated binary",
		"Print version number"}

	fmt.Print("usage: ntwrk <command> [arguments]\n\n")
	fmt.Print("commands:\n")
	for i, cmd := range cmds {
		fmt.Printf("    %v\t%v\n", cmd, descriptions[i])
	}
}
