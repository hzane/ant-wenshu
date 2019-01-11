package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
	"golang.org/x/net/publicsuffix"
	"jaytaylor.com/html2text"
)

// WenshuBrowser ...
type WenshuBrowser struct {
	http.Client
	UserAgent          func() string
	RewriteRequest     func(*http.Request) *http.Request
	js                 *otto.Otto
	guid               string
	number             string
	vl5x               string
	referer            string
	Locations          []string // higher-courts.csv
	intermediateCourts []string // intermediate-courts.csv
	primaryCourts      []string // primary-courts.csv
	causes             []string
	levels             map[string]string //
	types              map[string]string //
	docTypes           map[string]string //
	counts             doc
	tree               doc
}

// NewWenshu ...
func NewWenshu(repo string) *WenshuBrowser {
	ret := &WenshuBrowser{
		Client:         http.Client{Timeout: time.Second * 30},
		UserAgent:      macChrome,
		RewriteRequest: nop,
		guid:           GUID(),
		js:             otto.New(),
		levels: map[string]string{
			"全部":   "all",
			"最高法院": "1",
			"高级法院": "2",
			"中级法院": "3",
			"基层法院": "4",
		},
		types: map[string]string{
			"刑事案件": "1",
			"民事案件": "2",
			"行政案件": "3",
			"赔偿案件": "4",
			"执行案件": "5",
		},
		docTypes: map[string]string{
			"全部":  "all",
			"判决书": "1",
			"裁定书": "2",
			"调解书": "3",
			"决定书": "4",
			"通知书": "5",
			"批复":  "6",
			"答复":  "7",
			"函":   "8",
			"令":   "9",
			"其他":  "10",
		},
	}
	ret.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	_ = ret.RunJSFile(path.Join(repo, "base64.js"))
	_ = ret.RunJSFile(path.Join(repo, "md5.js"))
	_ = ret.RunJSFile(path.Join(repo, "vl5x.js"))
	_ = ret.RunJSFile(path.Join(repo, "unzip.js"))
	_ = ret.RunJSFile(path.Join(repo, "sha1.js"))
	_ = ret.RunJSFile(path.Join(repo, "rawdeflate.js"))
	_ = ret.RunJSFile(path.Join(repo, "rawinflate.js"))
	_ = ret.RunJSFile(path.Join(repo, "core-min.js"))
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
func (c *WenshuBrowser) QueryReferer(word, ay, fycj, spcx, ajlx, wslx, slfy, cprq string) string {
	k := "conditions"
	p := "searchWord %s %s  %s:%s"
	params := url.Values{}
	params.Set("sortType", "1")
	if ay != "" {
		if word == "" {
			word = "002006002001002"
		}
		params.Add(k, fmt.Sprintf(p, word, "AY", "案由", ay))
	}
	if fycj != "" {
		if word == "" {
			word, _ = c.levels[fycj]
		}
		params.Add(k, fmt.Sprintf(p, word, "FYCJ", "法院层级", fycj))
	}
	if spcx != "" {
		if word == "" {
			word = "all_" + spcx
		}
		params.Add(k, fmt.Sprintf(p, word, "SPCX", "审判程序", spcx))
	}
	// 2019-01-09  TO  2019-01-10
	if cprq != "" {
		if word == "" {
			word = " "
		}
		params.Add(k, fmt.Sprintf(p, word, "CPRQ", "裁判日期", cprq))
	}
	if ajlx != "" {
		if word == "" {
			word, _ = c.types[ajlx]
		}
		params.Add(k, fmt.Sprintf(p, word, "AJLX", "案件类型", ajlx))
	}
	if wslx != "" {
		if word == "" {
			word, _ = c.docTypes[wslx]
		}
		params.Add(k, fmt.Sprintf(p, word, "WSLX", "文书类型", wslx))
	}
	if slfy != "" {
		word = slfy
		params.Add(k, fmt.Sprintf(p, word, "SLFY", "法院名称", slfy))
	}
	if len(params[k]) == 0 {
		params.Add(k, fmt.Sprintf(p, "1", "AJLX", "案件类型", "刑事案件"))
	}
	params.Set("number", c.number)
	params.Set("guid", c.guid)
	uri, _ := url.Parse(WSPREFIX + "/List/List")
	uri.RawQuery = params.Encode()
	c.referer = uri.String()
	return c.referer
}
func (c *WenshuBrowser) ListList(word, ay, fycj, spcx, ajlx, wslx, slfy, cprq string) error {
	// sorttype=1&conditions=searchWord+1+AJLX++案件类型:刑事案件
	// uri := WSPREFIX + `/List/List?sorttype=1&conditions=searchWord+1+AJLX++%E6%A1%88%E4%BB%B6%E7%B1%BB%E5%9E%8B:%E5%88%91%E4%BA%8B%E6%A1%88%E4%BB%B6`
	uri := c.QueryReferer(word, ay, fycj, spcx, ajlx, wslx, slfy, cprq)
	return c.Touch(uri)
}

func (c *WenshuBrowser) TreeList() (doc []map[string]interface{}, err error) {
	uri := WSPREFIX + `/List/TreeList`
	req, _ := http.NewRequest("POST", uri, nil)
	return c.DoJSN(req)
}
func (c *WenshuBrowser) TreeContent(ay, fycj, spcx, ajlx, wslx, slfy, cprq string) (entries []map[string]interface{}, err error) {
	var params []string
	if ay != "" {
		params = append(params, "案由:"+ay)
	}
	if fycj != "" {
		params = append(params, "法院层级:"+fycj)
	}
	if spcx != "" {
		params = append(params, "审判程序:"+spcx)
	}
	// 2019-01-09  TO  2019-01-10
	if cprq != "" {
		params = append(params, "裁判日期:"+cprq)
	}
	if ajlx != "" {
		params = append(params, "案件类型:"+ajlx)
	}
	if wslx != "" {
		params = append(params, "文书类型:"+wslx)
	}
	if slfy != "" {
		params = append(params, "法院名称:"+slfy)
	}
	data, err := c.Submit(WSPREFIX+"/List/TreeContent", "",
		"Param", strings.Join(params, ","), // "案件类型:刑事案件",
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

// CourtTreeContent 案件类型:民事案件,中级法院:北京市第一中级人民法院
// parval: 北京市第一中级人民法院
func (c *WenshuBrowser) CourtTreeContent(casetyp, level, parent string) ([]map[string]interface{}, error) {
	if casetyp == "" {
		casetyp = "民事案件"
	}
	data, err := c.Submit(WSPREFIX+"/List/CourtTreeContent", "", "Param", sparam("案件类型", casetyp, level, parent), "parval", parent)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		err = fmt.Errorf("typ: %s, level: %s, parent:%s\n", casetyp, level, parent)
	}
	info(string(data))
	return c.JSJSON(data)
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

// GetDicValue AJLX 1,1,1,1,1,1,1
func (c *WenshuBrowser) GetDicValue(typ string, keys ...string) (string, error) {
	uri := WSPREFIX + "/List/GetDicValue"
	data, err := c.Submit(uri, "", "dicId", typ, "dicKey", strings.Join(keys, ","))
	return string(data), err
}
func (c *WenshuBrowser) CreateContentJS(docid string) (map[string]interface{}, error) {
	uri, _ := url.Parse(WSPREFIX + "/CreateContentJS/CreateContentJS.aspx")
	params := uri.Query()
	params.Set("DocID", docid)
	uri.RawQuery = params.Encode()

	ref, _ := url.Parse(WSPREFIX + "/content/content?KeyWord=")
	ref.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", uri.String(), nil)
	body, err := c.DoData(req)

	if err != nil {
		return nil, err
	}
	r1 := regexp.MustCompile(`JSON.stringify\(({.*?})\);`)
	result := r1.FindStringSubmatch(string(body))
	if len(result) == 0 {
		return nil, fmt.Errorf("invalid r1")
	}
	var doc map[string]interface{}
	tmp := result[0] // tmp is javascript json object, but we treat it as a json string
	fmt.Println(tmp)

	x, err := c.js.Run(tmp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(x.String()), &doc)

	r1 = regexp.MustCompile(`jsonHtmlData\s?=\s?("{.*}");`)
	tmp = r1.FindStringSubmatch(string(body))[1]
	x, err = c.js.Run("(" + tmp + ")")

	err = json.Unmarshal([]byte(x.String()), &doc)
	if htm, ok := doc["Html"].(string); ok {
		if txt, err := html2text.FromString(htm, html2text.Options{PrettyTables: true}); err == nil {
			doc["text"] = txt
		}
	}
	return doc, err
}
func (c *WenshuBrowser) RunJSFile(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	_, err = c.js.Run(data)
	return err
}

// ListContent "案件类型","刑事案件"
// word, ay, fycj, spcx, ajlx, wslx, cprq string
func (c *WenshuBrowser) ListContent(
	idx string, order, direction string, ay,
	fycj, spcx, ajlx, wslx, slfy, cprq string) (cases []doc, err error) {
	if order == "" {
		order = "裁判日期"
	}
	if direction == "" {
		direction = "desc"
	}
	if idx == "" {
		idx = "1"
	}
	var params []string
	if ay != "" {
		params = append(params, "案由:"+ay)
	}
	if fycj != "" {
		params = append(params, "法院层级:"+fycj)
	}
	if spcx != "" {
		params = append(params, "审判程序:"+spcx)
	}
	// 2019-01-09  TO  2019-01-10
	if cprq != "" {
		params = append(params, "裁判日期:"+cprq)
	}
	if ajlx != "" {
		params = append(params, "案件类型:"+ajlx)
	}
	if wslx != "" {
		params = append(params, "文书类型:"+wslx)
	}
	if slfy != "" {
		params = append(params, "法院名称:"+slfy)
	}
	if len(params) == 0 {
		params = append(params, "案件类型:刑事案件")
	}
	info("ListContent", idx, order, direction, "\n")
	info(strings.Join(params, ","))
	data, err := c.Submit(WSPREFIX+"/List/ListContent", "", "Param",
		strings.Join(params, ","),
		"Index", idx,
		"Page", "10",
		"Order", order,
		"Direction", direction,
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

func (c *WenshuBrowser) GetCode() (string, error) {
	uri := WSPREFIX + `/ValiCode/GetCode`
	body, err := c.Submit(uri, "", "guid", c.guid)
	c.number = string(body)
	return c.number, err
}

func (c *WenshuBrowser) Submit(uri, refer string, kv ...string) ([]byte, error) {
	req, _ := http.NewRequest("POST", uri, form(kv...))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if refer != "" {
		req.Header.Set("Referer", refer)
	}
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

// Do ...
func (c *WenshuBrowser) Do(req *http.Request) (*http.Response, error) {
	if c.RewriteRequest != nil {
		req = c.RewriteRequest(req)
	}
	req.Header.Set("accept", "*/*")
	// req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("accept-language", "en,en-US;q-0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6")
	req.Header.Set("DNT", "1")
	if req.Header.Get("origin") == "" {
		req.Header.Set("origin", WSPREFIX)
	}
	if req.Header.Get("user-agent") == "" {
		req.Header.Set("user-agent", c.UserAgent())
	}
	if req.Header.Get("referer") == "" {
		req.Header.Set("referer", c.referer)
	}
	info(req.Method, req.URL.String())
	printHeaders(req.Header)
	resp, err := c.Client.Do(req)
	if err == nil && resp.StatusCode < http.StatusInternalServerError {
		printHeaders(resp.Header)
		return resp, err
	}
	resp, err = c.Client.Do(req) // try again
	if resp != nil {
		printHeaders(resp.Header)
	}
	return resp, err
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