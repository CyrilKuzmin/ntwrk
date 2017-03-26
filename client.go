package main

import (
	"fmt"
	"net"
	"time"
)

const protoFmt = "ntwrk%s :%s\r\n"
const timeout = time.Duration(15) * time.Second

// testContext holds a test function, action name, and address to connect to.
type testContext struct {
	Action string
	Fn     func(net.Conn, time.Duration) (int64, error)
	Addr   string
}

// startClient starts the network test suite.
func startClient(host string) {
	perform(testContext{"download", download, host})
	perform(testContext{"upload", upload, host})
}

// perform runs a network test and prints the recorded bandwidth.
func perform(ctx testContext) {
	conn := openConn(ctx.Addr, ctx.Action)

	since := time.Now()
	ticker := time.NewTicker(time.Millisecond * 150)
	go func() {
		for t := range ticker.C {
			elapsed := t.Sub(since)
			progress := formatProgress(elapsed, timeout)
			fmt.Printf("\r %9s: %s", ctx.Action, progress)
		}
	}()
	defer ticker.Stop()

	bytes, err := ctx.Fn(conn, timeout)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	elapsed := time.Since(since).Seconds()
	fmt.Printf(" %s\n", formatBytes(bytes, elapsed))
}

// openConn opens a connection to `host` and writes a formatted message to it.
func openConn(host string, action string) (conn net.Conn) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, protoFmt, proto, action)
	return
}

// whoami requests the client's external IP address from `host` and prints it.
func whoami(host string) {
	resp := make([]byte, 40)
	conn := openConn(host, "whoami")
	conn.Read(resp)
	fmt.Print(string(resp))
}
