package proxy

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

type server struct {
	laddr  string
	client *http.Client
}

func NewServer(laddr string) *server {
	return &server{laddr: laddr, client: &http.Client{}}
}

func (s *server) Run() {
	config := &tls.Config{Certificates: make([]tls.Certificate, 1)}
	config.Certificates[0], _ = tls.X509KeyPair([]byte(crt), []byte(key))
	ln, err := tls.Listen("tcp", s.laddr, config)
	if err != nil {
		return
	}

	log.Printf("proxy_server listen on '%s'\n", s.laddr)

	go func() {
		log.Fatal(http.Serve(ln, s))
	}()
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	realto := r.Header.Get("realto")

	req, err := http.NewRequest(r.Method, realto, r.Body)
	if err != nil {
		return
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	rw.WriteHeader(resp.StatusCode)
	rw.Write(data)

}
