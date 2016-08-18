/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/07/18 10:00
 */

package proxy

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
)

type client struct {
	laddr string
	raddr string
}

func NewClient(laddr string, raddr string) *client {
	return &client{laddr: laddr, raddr: raddr}
}

func (c *client) Run() {
	log.Printf("proxy_client listen on '%s' with remote '%s'\n", c.laddr, c.raddr)

	go func() {
		log.Fatal(http.ListenAndServe(c.laddr, c))
	}()
}

func (c *client) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("tcp", c.raddr)
	if err != nil {
		return
	}
	defer conn.Close()

	conns := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
	r.Header.Add("realto", r.RequestURI)
	r.Write(conns)

	reader := bufio.NewReader(conns)
	resp, err := http.ReadResponse(reader, r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	for k, value := range resp.Header {
		for _, vv := range value {
			rw.Header().Add(k, vv)
		}
	}

	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}
