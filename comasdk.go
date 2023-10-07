package comasdkgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type Coma struct {
	RetryWaitTime time.Duration
	Retry         int
	origin        string
	host          string
	port          string
	key           string
	conn          *websocket.Conn
}

type ComaOption func(c *Coma)

func SetRetry(retry int) ComaOption {
	return func(c *Coma) {
		c.Retry = retry
	}
}

func SetRetryWaitTime(retryWaitTime time.Duration) ComaOption {
	return func(c *Coma) {
		c.RetryWaitTime = retryWaitTime
	}
}

// Coma constructor will establish connection to coma server
// if the connection establishment is fail, it return error
func New(origin, host, port, key string, opts ...ComaOption) (*Coma, error) {
	coma := &Coma{
		origin: origin,
		host:   host,
		port:   port,
		key:    key,
	}

	for _, opt := range opts {
		opt(coma)
	}

	if coma.Retry == 0 {
		coma.Retry = 1
	}

	err := coma.connect()
	if err != nil {
		return nil, err
	}
	return coma, nil
}

func (c *Coma) connect() error {
	var (
		conn *websocket.Conn
		err  error
	)

	for i := 1; i <= c.Retry; i++ {
		conn, err = websocket.Dial(
			fmt.Sprintf("ws://%s:%s/websocket?authorization=%s", c.host, c.port, c.key),
			"tcp", c.origin,
		)
		if err == nil {
			break
		}
		if err != nil {
			log.Printf("connecting retry: %d\n", i)
		}

		time.Sleep(c.RetryWaitTime)
	}

	if err != nil {
		return err
	}

	c.conn = conn
	log.Println("connected")
	return nil
}

// Observe uses for observing data from coma server
func (c *Coma) Observe(observer interface{}) error {
	var (
		message Message
		err     error
	)

	if c.conn == nil {
		return errors.New("err: connection does not found")
	}

	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() error {
		for {
			err = websocket.JSON.Receive(c.conn, &message)
			if err == io.EOF {
				if err := c.connect(); err != nil {
					return errors.New("err: reconnecting error")
				}
				err = nil
			}
			if err != nil {
				return err
			}

			err = json.Unmarshal(message.Data, observer)
			if err != nil {
				return err
			}
		}
	}()

	return nil
}

func (c *Coma) Shutdown(ctx context.Context) error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
