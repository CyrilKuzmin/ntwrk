package ntwrk

import (
	"fmt"
	"net"
	"time"
)

const protoFmt = "ntwrk%s :%s\r\n"

type Client struct {
	pingCount int
	timeout   time.Duration
	addr      string
	cli       bool
}

func NewClient(host string, port int) *Client {
	addr := fmt.Sprintf("%s:%d", host, port)
	return &Client{
		pingCount: 10,
		timeout:   time.Second * time.Duration(5),
		addr:      addr,
		cli:       true,
	}
}

// testContext holds a test function, action name, and address to connect to.
type testContext struct {
	Action string
	Fn     func(net.Conn, time.Duration) (int64, error)
	cl     *Client
}

// StartCLI starts the network test suite.
func (c *Client) StartCLI() {
	err := c.ping()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	ds, err := perform(testContext{"download", download, c})
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Printf(" %v\n", ds)
	us, err := perform(testContext{"upload", upload, c})
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Printf(" %v\n", us)
}

func (c *Client) Measure() (string, string, error) {
	c.cli = false
	ds, err := perform(testContext{"download", download, c})
	if err != nil {
		return "", "", err
	}
	us, err := perform(testContext{"upload", upload, c})
	if err != nil {
		return "", "", err
	}
	return ds, us, nil
}

// perform runs a network test and prints the recorded bandwidth.
func perform(ctx testContext) (res string, err error) {
	conn, err := openConn(ctx.cl.addr, ctx.Action)
	if err != nil {
		return
	}

	since := time.Now()
	if ctx.cl.cli {
		ticker := time.NewTicker(time.Millisecond * 150)
		go func() {
			for t := range ticker.C {
				elapsed := t.Sub(since)
				progress := formatProgress(elapsed, ctx.cl.timeout)
				fmt.Printf("\r %9s: %s", ctx.Action, progress)
			}
		}()
		defer ticker.Stop()
	}

	bytes, err := ctx.Fn(conn, ctx.cl.timeout)
	if err != nil {
		return
	}

	elapsed := time.Since(since).Seconds()
	return formatBytes(bytes, elapsed), nil
}

// ping performs a network latency test.
func (c *Client) ping() error {
	conn, err := openConn(c.addr, "echo")
	if err != nil {
		return err
	}
	resp := make([]byte, 6)
	since := time.Now()
	for i := 0; i < c.pingCount; i++ {
		conn.Write([]byte("echo\r\n"))
		conn.Read(resp)
		if string(resp) != "echo\r\n" {
			return fmt.Errorf("invalid echo reply")
		}
	}
	elapsed := time.Duration(int(time.Since(since).Milliseconds())/c.pingCount) * time.Millisecond
	if c.cli {
		fmt.Printf("\r %9s: %s\n", "latency", elapsed)
	}
	return nil
}

// Whoami requests the client's external IP address from `host` and prints it.
func (c *Client) Whoami() (string, error) {
	resp := make([]byte, 40)

	conn, err := openConn(c.addr, "whoami")
	if err != nil {
		return "", err
	}

	conn.Read(resp)
	if c.cli {
		fmt.Print(string(resp))
	}
	return string(resp), nil
}

// openConn opens a connection to `host` and writes a formatted message to it.
func openConn(host string, action string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", host)
	if err != nil {
		return
	}
	fmt.Fprintf(conn, protoFmt, proto, action)
	return
}
