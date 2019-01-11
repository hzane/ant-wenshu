package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

// GUID ...
func GUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	return fmt.Sprintf("%x-%x-%x%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// WSPREFIX ...
const (
	WSPREFIX = "http://wenshu.court.gov.cn"
	macc     = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_%d_%d) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.%d.%d.%d Safari/%d.%d`
)

func macChrome() string {
	m1, m2 := rand.Intn(4)+7, rand.Intn(8)
	v1, v2, v3, v4 := rand.Intn(40)+27, rand.Intn(10), rand.Intn(4000), rand.Intn(1000)
	s1, s2 := rand.Intn(30)+510, rand.Intn(50)
	return fmt.Sprintf(macc, m1, m2, v1, v2, v3, v4, s1, s2)
}

func form(kv ...string) io.Reader {
	data := params(kv...)
	printv(data)
	return strings.NewReader(data.Encode())
}
func printv(vals url.Values) {
	for k, items := range vals {
		for _, item := range items {
			info(k, item)
		}
	}
}

func params(kv ...string) url.Values {
	data := url.Values{}
	for i := 0; i < len(kv); i += 2 {
		data.Add(kv[i], kv[i+1])
	}
	return data
}
func sparam(kv ...string) string {
	var params []string
	for i := 0; i < len(kv); i += 2 {
		params = append(params, kv[0]+":"+kv[1])
	}
	return strings.Join(params, ",")
}
func nop(r *http.Request) *http.Request {
	return r
}
func panice(err error) {
	if err != nil {
		panic(err)
	}
}

func xclose(rc io.ReadCloser) {
	_ = rc.Close()
}

func printHeaders(headers http.Header) {
	for k, vals := range headers {
		inf("%s:\t%s\n", k, vals[0])
	}
}
func printd(d []doc) {
	data, _ := json.MarshalIndent(d, "", "  ")
	fmt.Println(string(data))
}
