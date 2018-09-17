package main

import (
	"math/rand"
	"net/http"
	"time"
)

func verb(method, uri string, resp *http.Response, err error) {
	var status string
	var code int
	var cl int64
	if resp != nil {
		status, code, cl = resp.Status, resp.StatusCode, resp.ContentLength
	}
	info(cl, code, status, err, method, uri)
}
func randSleep() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
}
