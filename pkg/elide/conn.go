package elide

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

func (c *Client) MakeRequest(buf []byte) ([]byte, error) {
	switch c.Protocol {
	case 1:
		return c.makeTCPRequest(buf)
	default:
		return c.makeHTTPRequest(buf)
	}
}

func (c *Client) makeHTTPRequest(buf []byte) ([]byte, error) {
	r, err := http.NewRequest(http.MethodPost, c.addr, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	buf, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return buf, r.Body.Close()
}

func (c *Client) makeTCPRequest(buf []byte) ([]byte, error) {
	_, err := c.conn.Write(append(buf, 0))
	if err != nil {
		return nil, err
	}
	buf, err = bufio.NewReader(c.conn).ReadBytes(0)
	if err != nil {
		return nil, err
	}
	return buf[:len(buf)-1], nil
}
