package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"math/rand"
	"gitlab.com/hearts.zhang/tools"
	"github.com/robertkrimen/otto"
	// _ "github.com/robertkrimen/otto/underscore"

	"io/ioutil"
	"path"
	"fmt"
	"time"
	"net/http"
	"bytes"
)

func main() {
	guid := GUID()

	client := tools.NewHTTPClient2(time.Second*15, 2, nil, nil)
	client.RewriteRequest = appendHeader
	number := GetCode(client, guid)
	fmt.Println(guid, number)
	vjkl5 := VJKL5(client, guid, number)
	vx, err := vl5x(vjkl5)
	fmt.Println(vjkl5, vx, err)
}
func appendHeader(req *http.Request) *http.Request {
	req.Header.Set("Origin", host)
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", host)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	return req
}

const codeURL = "http://wenshu.court.gov.cn/ValiCode/GetCode"
const host = "http://wenshu.court.gov.cn"

// GetCode ...
func GetCode(client *tools.Client, guid string) (number string) {
	data := url.Values{}
	data.Set("guid", guid)

	req, _ := http.NewRequest("POST", codeURL, bytes.NewBufferString(data.Encode()))
	//req.Header.Set("Origin", host)
	//req.Header.Set("Referer", host)
	// req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	number = string(b)
	return
}

// GUID ...
func GUID() (string) {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

const listURL = `http://wenshu.court.gov.cn/list/list/?sorttype=1&number=%v&guid=%v&conditions=searchWord+QWJS+++%v`

// VJKL5 ...
func VJKL5(client *tools.Client, guid, number string) (string) {
	uri := fmt.Sprintf(listURL, number, guid, url.QueryEscape("全文检索:三年"))

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := client.Do(req)
	_, _ = resp, err
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	for _, ck := range resp.Cookies() {
		if ck.Name == "vjkl5" {
			return ck.Value
		}
	}
	return ""
	//b, _ := ioutil.ReadAll(resp.Body)
	//return string(b)
}
func vl5x(vjkl5 string) (string, error) {

	vm := otto.New()
	compile(vm, path.Join(config.js, `md5.js`),
		path.Join(config.js, `sha1.js`),
		path.Join(config.js, `base64.js`),
		path.Join(config.js, `vl5x.js`))

	value, err := vm.Run(`vl5x("a9993e364706816aba3e25717850c26c9cd0d89d");`)
	if err != nil {
		return "", err
	}
	return value.ToString()
}
func compile(vm *otto.Otto, files ...string) {
	for _, file := range files {
		if bjs, err := ioutil.ReadFile(file); err == nil {
			vm.Run(string(bjs))
		}
	}
}
func main_() {

	s := tools.NewSpider(config.cuckoo,
		config.repo,
		config.domain,
		8<<20,
		config.workers,
	)
	s.Info = log.Println
	s.Accept = RejectAPK
	//s.Lookup = tools.Lookup(config.hostF, config.domain)

	cancel := s.BootInfinite(config.bootstrap)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	s.Close(cancel)
}

// RejectAPK ...
// http://www.wandoujia.com/apps/com.sohu.inputmethod.sogou/binding?source=web_inner_referral_binded
func RejectAPK(uri *url.URL) bool {
	ad := tools.AcceptDomain(config.domain)
	if !ad(uri) {
		return false
	}

	x := !strings.HasSuffix(uri.Path, "/binding")
	x = x && !strings.HasSuffix(uri.Path, "/download")
	x = x && !strings.HasSuffix(uri.Path, "/comment1")
	x = x && !strings.HasSuffix(uri.Path, "/history")
	x = x && !strings.HasSuffix(uri.Path, "award")
	x = x && !strings.HasSuffix(uri.Path, "/help")
	x = x && !strings.HasSuffix(uri.Path, "/wdjweb/recommend")
	x = x && !strings.HasPrefix(uri.Path, "/wdjweb/faq")

	x = x && uri.Fragment == ""

	return x
}

func init() {
	flag.StringVar(&config.js, "js-dir", ".", "javascript file folder")
	flag.StringVar(&config.repo, "repo", "/repo/spiding/wenshu/repo", "")
	flag.StringVar(&config.bootstrap, "bootstrap", "https://wenshu.court.gov.cn/", "")
	flag.StringVar(&config.hostF, "host-file", "", "")
	flag.StringVar(&config.cuckoo, "cuckoo", "/repo/spiding/wenshu/wenshu.cuckoo", "")
	flag.StringVar(&config.domain, "domain", "wenshu.court.gov.cn", "")
	flag.StringVar(&config.proxies, "proxies", "", "")
	flag.IntVar(&config.workers, "workers", 1, "")

	flag.Parse()

}

var config struct {
	js        string
	domain    string
	hostF     string
	bootstrap string
	repo      string
	cuckoo    string
	proxies   string
	workers   int
}

/*
javascript:Navi("DcKOwrcRw4AwDMOEVmISw7UswpnCtMO/SHYNw5wBFcKuw5vDtz7DtlLCsxkhw53DgxTCksKnWsOgdkBVw53CtCrCggXCosKDdjkpw7wROsKbwpEWw5oEZcKfwpzCjcKneMKmw5Zxw7zCshp2Zz3Dv0QVwr/DsMKqw7QgdWjDl3AIw4Mzw5XDtxrDi1vCtgk6bsK9P8Kbw5DDh2DDiMKJwqHDrnMZGm8JwqHCiMOyRsOiXsOJw7FTGC0pw7BvwpFbworDlT82w7oB","")*/
