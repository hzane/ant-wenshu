package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	cj "net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// WSPREFIX ...
const (
	WSPREFIX = "http://wenshu.court.gov.cn"
)

var (
	info = log.Println
	inf  = log.Printf
)

type (
	doc = map[string]interface{}
)

func main() {
	wb := NewWenshu()
	err := wb.Home()
	panice(err)

	err = wb.Criminal()
	panice(err)

	v := wb.GetCookie("vjkl5")
	info("vjkl5\t: ", v)
	_, _ = wb.js.Run(fmt.Sprintf(`function getCookie(cname) {
    return "%s"
  }`, v))
	vlx, err := wb.VL5X()
	panice(err)
	info("vl5x :\t", vlx)

	entries, err := wb.TreeList()
	panice(err)
	printd(entries)

	courts, err := wb.GetCourt("北京")
	panice(err)
	printd(courts)

	// courts, err = wb.GetChildAllCourt(courts)
	// panice(err)
	// printd(courts)

	// cases, err := wb.ListContent()
	// panice(err)
	// printd(cases)
	entries, err = wb.TreeContent()
	panice(err)
	printd(entries)
}
func init() {

}

// GUID ...
func GUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	return fmt.Sprintf("%x-%x-%x%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// WenshuBrowser ...
type WenshuBrowser struct {
	http.Client
	UserAgent      func() string
	RewriteRequest func(*http.Request) *http.Request
	js             *otto.Otto
	guid           string
	number         string
	vl5x           string
	counts         doc
	tree           doc
}

// NewWenshu ...
func NewWenshu() *WenshuBrowser {
	ret := &WenshuBrowser{
		Client:         http.Client{Timeout: time.Second * 30},
		UserAgent:      macChrome,
		RewriteRequest: nop,
		guid:           GUID(),
		js:             otto.New(),
	}
	ret.Jar, _ = cj.New(&cj.Options{PublicSuffixList: publicsuffix.List})
	_ = ret.RunJSFile("base64.js")
	_ = ret.RunJSFile("md5.js")
	_ = ret.RunJSFile("vl5x.js")
	_ = ret.RunJSFile("unzip.js")
	_ = ret.RunJSFile("sha1.js")
	_ = ret.RunJSFile("rawdeflate.js")
	_ = ret.RunJSFile("rawinflate.js")
	_ = ret.RunJSFile("core-min.js")
	return ret
}

// Get ...
func (c *WenshuBrowser) Touch(uri string) error {
	req, _ := http.NewRequest("GET", uri, nil)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer xclose(resp.Body)
	_, err = io.Copy(ioutil.Discard, resp.Body)
	info(resp.StatusCode, uri)
	printHeaders(resp.Header)
	return err
}

func (c *WenshuBrowser) Home() error {
	return c.Touch(WSPREFIX)
}

func (c *WenshuBrowser) Criminal() error {
	uri := WSPREFIX + `/List/List?sorttype=1&conditions=searchWord+1+AJLX++%E6%A1%88%E4%BB%B6%E7%B1%BB%E5%9E%8B:%E5%88%91%E4%BA%8B%E6%A1%88%E4%BB%B6`
	return c.Touch(uri)
}

func (c *WenshuBrowser) TreeList() (doc []map[string]interface{}, err error) {
	uri := WSPREFIX + `/List/TreeList`
	req, _ := http.NewRequest("POST", uri, nil)
	return c.DoJSN(req)
}
func (c *WenshuBrowser) TreeContent() (entries []map[string]interface{}, err error) {
	data, err := c.Submit(WSPREFIX+"/List/TreeContent", "",
		"Param", "案件类型:刑事案件",
		"vl5x", c.vl5x,
		"guid", c.guid,
		"number", c.number)
	if err != nil {
		return
	}
	entries, err = c.JSJSON(data)
	return
}
func (c *WenshuBrowser) GetCourt(province string) (docs []map[string]interface{}, err error) {
	data, err := c.Submit(WSPREFIX+"/Index/GetCourt", "", "province", province)
	if err != nil {
		return
	}

	s, err := c.JSString(data)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(s), &docs)
	return
}
func (c *WenshuBrowser) GetChildAllCourt(courts []map[string]interface{}) (rets []map[string]interface{}, err error) {
	var keys []string
	for _, court := range courts {
		if v, ok := court["key"]; ok {
			if v, ok := v.(string); ok {
				keys = append(keys, v)
			}
		}
	}
	data, err := c.Submit(WSPREFIX+"/Index/GetChildAllCourt", "", "keyCodeArrayStr", strings.Join(keys, ","))
	if err != nil {
		return
	}
	s, err := c.JSString(data)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(s), &rets)
	return
}
func (c *WenshuBrowser) VL5X() (ret string, err error) {
	x, err := c.js.Run(`getKey()`)
	if err != nil {
		return
	}
	ret, err = x.ToString()
	c.vl5x = ret
	return
}
func (c *WenshuBrowser) RunJSFile(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	_, err = c.js.Run(data)
	return err
}
func (c *WenshuBrowser) ListContent() (cases []doc, err error) {
	data, err := c.Submit(WSPREFIX+"/List/ListContent", "", "Param", "案件类型:刑事案件",
		"Index", "1",
		"Page", "10",
		"Order", "法院层级",
		"Direction", "asc",
		"vl5x", c.vl5x,
		"number", c.number,
		"guid", c.guid)
	if err != nil {
		return
	}
	cases, err = c.JSJSON(data)
	if len(cases) <= 1 {
		err = fmt.Errorf(string(data))
		return
	}

	runeval, _ := cases[0]["RunEval"].(string)
	_, err = c.js.Run(fmt.Sprintf(`eval(unzip("%s"))`, runeval))
	if err != nil {
		return
	}
	cases = cases[1:]
	for _, kase := range cases {
		did, _ := kase["文书ID"].(string)
		x, err := c.js.Run(fmt.Sprintf(`com.str.Decrypt(unzip("%s"))`, did))
		if err != nil {
			info(did, err)
		}

		did, _ = x.ToString()
		kase["文书ID"] = did

	}
	return
}
func (c *WenshuBrowser) JSJSON(data []byte) (docs []doc, err error) {
	s, err := c.JSString(data)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(s), &docs)
	if err != nil {
		s = strings.Replace(s, `\"`, `"`, -1)
		err = json.Unmarshal([]byte(s), &docs)
	}
	if s == "remind" { // hold a while
		err = fmt.Errorf("%s", s)
	}
	return
}
func (c *WenshuBrowser) JSString(data []byte) (ret string, err error) {
	val, err := c.js.Run(data)
	if err != nil {
		return
	}

	ret, err = val.ToString()
	val, err = c.js.Run(fmt.Sprintf(`JSON.stringify(%s)`, ret)) // -> json string
	ret, err = val.ToString()
	return
}

func (c *WenshuBrowser) DoJSN(req *http.Request) (docs []map[string]interface{}, err error) {
	data, err := c.DoData(req)
	if err != nil {
		return
	}
	s, err := c.JSString(data)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(s), &docs)
	return
}
func (c *WenshuBrowser) GetVJKL5() string {
	return c.GetCookie("vjkl5")
}

func (c *WenshuBrowser) GetCookie(name string) string {
	root, _ := url.Parse(WSPREFIX)
	for _, cookie := range c.Jar.Cookies(root) {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func (c *WenshuBrowser) GetAllCountRefresh() error {
	uri := WSPREFIX + `/Index/GetAllCountRefresh?refresh=Refresh`
	req, _ := http.NewRequest("POST", uri, nil)
	req.Header.Set("Origin", WSPREFIX)
	req.Header.Set("Referer", WSPREFIX)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	doc, err := c.DoJSON(req)
	// fmt.Println(doc)
	if err == nil {
		c.counts = doc
	}
	return err
}

func (c *WenshuBrowser) Submit(uri, refer string, kv ...string) ([]byte, error) {
	req, _ := http.NewRequest("POST", uri, form(kv...))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", refer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer xclose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *WenshuBrowser) GetCode() (string, error) {
	uri := WSPREFIX + `/ValiCode/GetCode`
	body, err := c.Submit(uri, "guid", c.guid)
	c.number = string(body)
	return c.number, err
}

func form(kv ...string) io.Reader {
	data := url.Values{}
	for i := 0; i < len(kv); i += 2 {
		data.Add(kv[i], kv[i+1])
	}
	return strings.NewReader(data.Encode())
}

// Do ...
func (c *WenshuBrowser) Do(req *http.Request) (*http.Response, error) {
	if c.RewriteRequest != nil {
		req = c.RewriteRequest(req)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("accept-language", "en,en-US;q-0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6")
	req.Header.Set("DNT", "1")
	if req.Header.Get("origin") == "" {
		req.Header.Set("origin", WSPREFIX)
	}
	if req.Header.Get("user-agent") == "" {
		req.Header.Set("user-agent", c.UserAgent())
	}
	if req.Header.Get("referer") == "" {
		req.Header.Set("referer", WSPREFIX)
	}
	// printHeaders(req.Header)
	resp, err := c.Client.Do(req)
	if err == nil && resp.StatusCode < http.StatusInternalServerError {
		return resp, err
	}
	return c.Client.Do(req) // try again
}
func (c *WenshuBrowser) DoData(req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer xclose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	data, err1 := ioutil.ReadAll(resp.Body)
	if err == nil && err1 != nil {
		err = err1
	}
	return data, err
}
func (c *WenshuBrowser) DoJSON(req *http.Request) (map[string]interface{}, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer xclose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	var doc map[string]interface{}
	if err1 := json.NewDecoder(resp.Body).Decode(&doc); err == nil && err1 != nil {
		err = err1
	}
	return doc, err
}

var _ = NewWenshu

// chrome
const macc = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_%d_%d) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.%d.%d.%d Safari/%d.%d`

func macChrome() string {
	m1, m2 := rand.Intn(4)+7, rand.Intn(8)
	v1, v2, v3, v4 := rand.Intn(40)+27, rand.Intn(10), rand.Intn(4000), rand.Intn(1000)
	s1, s2 := rand.Intn(30)+510, rand.Intn(50)
	return fmt.Sprintf(macc, m1, m2, v1, v2, v3, v4, s1, s2)
}

// 指定baidu referer是不是可以更像一个浏览器
const ref = `https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=wenshu&rsv_pq=%x&rsv_t=&rqlang=cn&rsv_enter=1&rsv_sug3=6&rsv_sug1=5&rsv_sug7=100&rsv_sug2=0&inputT=753&rsv_sug4=%d`

func referer() string {
	r1, r2 := rand.Int63(), rand.Intn(10000)
	return fmt.Sprintf(ref, r1, r2)
}

type Task struct {
	ID            string      `json:"id"`
	URI           string      `json:"uri,omitempty"`
	REF           string      `json:"ref,omitempty"`
	Reason        string      `json:"reason,omitempty"`
	EXT           string      `json:"ext,omitempty"`
	Title         string      `json:"title,omitempty"`
	ContentType   string      `json:"content-type,omitempty"`
	FilePath      string      `json:"file-path,omitempty"`
	ContentLength int64       `json:"content-length,omitempty"`
	StatusCode    int         `json:"status-number,omitempty"`
	TTL           int         `json:"ttl,omitempty"`
	Headers       http.Header `json:"headers,omitempty"`
	error         error
	html          *html.Node
	doc           *gq.Document
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
