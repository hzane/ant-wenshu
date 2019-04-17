package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitlab.com/hearts.zhang/ants"
)

func main() {
	http.HandleFunc("/", wenshu)
	_ = http.ListenAndServe(":18071", nil)
}
func wenshu(w http.ResponseWriter, r *http.Request) {
	cookies := r.Header.Get("cookie")
	log.Println(r.URL)
	log.Println(cookies)
	if r.URL.Path == "/" || r.URL.Path == "/List/List" {
		wenshuWZWS(w, r)
		return
	}
	w.Header().Set("content-type", "text/html")
	_, _ = fmt.Fprintf(w, "<html><head></head><body>%s\n%s</body></html>", r.URL.String(), cookies)
}
func wenshuWZWS(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = "wenshu.court.gov.cn"
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}

	req, _ := http.NewRequest("GET", r.URL.String(), nil)
	for _, header := range []string{"accept-encoding", "accept", "accept-language", "cookie", "dnt", "user-agent"} {
		req.Header.Set(header, r.Header.Get(header))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("content-type", "text/html")
		_, _ = fmt.Fprintf(w, "<html><head></head><body>%s</body></html>", err.Error)
		return
	}

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
		log.Println(k, v)
	}
	_, _ = io.Copy(w, resp.Body)
	_ = resp.Body.Close()
}

func init() {
	flag.Parse()
	config.c = ants.NewHTTPClient(ants.HClientSettings{})
}

var config struct {
	c *ants.Client
}
