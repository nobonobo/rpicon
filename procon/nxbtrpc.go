package procon

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
	"sync"
	"time"
)

type State struct {
	State  string  `json:"state"`
	Errors *string `json:"errors,omitempty"`
}

type codec struct {
	rpc.ClientCodec
	in  *io.PipeReader
	out *io.PipeWriter
}

func (c *codec) WriteRequest(res *rpc.Request, v interface{}) error {
	//log.Println("req:", res, v)
	if err := c.ClientCodec.WriteRequest(res, v); err != nil {
		return err
	}
	return nil
	//_, err := c.out.Write([]byte("\n"))
	//return err
}

func (c *codec) Close() error {
	c.in.Close()
	c.out.Close()
	return c.ClientCodec.Close()
}

type Client struct {
	*rpc.Client
	mu  sync.Mutex
	cmd *exec.Cmd
}

func New() *Client {
	return &Client{}
}

func (c *Client) Start(script string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stop()
	c.cmd = exec.Command("python3", script)
	rin, win := io.Pipe()
	rout, wout := io.Pipe()
	c.cmd.Stdin = rin
	c.cmd.Stdout = wout
	//c.cmd.Stderr = os.Stderr
	cc := &codec{
		ClientCodec: jsonrpc.NewClientCodec(
			struct {
				io.Reader
				io.Writer
				io.Closer
			}{
				Reader: rout, //io.TeeReader(rout, os.Stderr),
				Writer: win,  //io.MultiWriter(win, os.Stderr),
				Closer: win,
			},
		),
		in:  rout,
		out: win,
	}
	c.Client = rpc.NewClientWithCodec(cc)
	return c.cmd.Start()
}

func (c *Client) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stop()
}

func (c *Client) stop() error {
	if c.cmd == nil {
		return nil
	}
	defer func() {
		c.Client = nil
		c.cmd = nil
	}()
	c.Call("close", nil, nil)
	c.Client.Close()
	cmd := c.cmd
	c.cmd = nil
	tm := time.AfterFunc(5*time.Second, func() {
		if p := cmd.Process; p != nil {
			p.Kill()
		}
	})
	if err := cmd.Wait(); err != nil {
		log.Print("stop:", err)
	}
	tm.Stop()
	return nil
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Call("disconnect", nil, nil)
	for r := 0; r < 3; r++ {
		if err := c.Call("connect", nil, nil); err != nil {
			return err
		}
		res := State{}
		score := 0
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			if err := c.Call("state", nil, &res); err != nil {
				return err
			}
			if res.State != "connected" {
				break
			}
			score++
			if score == 5 {
				return nil
			}
		}
		if res.State == "crashed" {
			continue
		}
	}
	return fmt.Errorf("connect failed: retry exceeded")
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Call("disconnect", nil, nil)
}

func (c *Client) Input(in Input) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Call("input", in, nil)
}

func (c *Client) State() (*State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := &State{}
	if err := c.Call("state", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func ExampleTest() {
	c := New()
	log.Print("start")
	defer log.Print("stop")
	if err := c.Start("procon.py"); err != nil {
		log.Fatal(err)
	}
	defer c.Stop()
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()
	for i := 0; i < 100; i++ {
		if err := c.Input(Input{}); err != nil {
			log.Println(err)
			if err := c.Connect(); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(16666 * time.Microsecond)
	}
}
