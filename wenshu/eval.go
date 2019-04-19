package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"gitlab.com/hearts.zhang/ants"
)

func main() {
	wb, err := NewWenshuBrowser("js")
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = wb.LoadCookies()
	// err = wb.Home()
	// fmt.Println("home", err)

	_, ret, err := wb.List("刑事案件")
	fmt.Println(err)
	data, _ := json.MarshalIndent(ret, "", "  ")
	fmt.Println(string(data))
}

// WenshuBrowser ...
type WenshuBrowser struct {
	GUID string
	c    *ants.Client
	jse  *otto.Otto
	home *url.URL
	Log  func(...interface{})
}

// NewWenshuBrowser ...
func NewWenshuBrowser(repo string) (wb *WenshuBrowser, err error) {
	u, _ := url.Parse(`http://wenshu.court.gov.cn/`)
	// User-Agent在后面的xhr请求中被严格验证，所以在所有请求中应该保持ua一致
	ua := `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3759.4 Safari/537.36`
	wb = &WenshuBrowser{
		home: u,
		c: ants.NewHTTPClient(ants.HClientSettings{
			UserAgent: func() string { return ua },
			Rewrite: func(r *http.Request) *http.Request {
				// 自定义的反向代理看到这个头后用自己记录下来的cookie更新请求中的cookie
				r.Header.Add("use-wenshu-cookie", "4bd54ad7a845")
				// r.Header.Add("use-wenshu-trailer", "4bd54ad7a845")
				return r
			},
		}),
		GUID: "7df4cd4f-0d91-23411994-4bd54ad7a845", // 好像没有地方验证 GUID的一致性
		jse:  nil,
		Log:  log.Println,
	}

	wb.jse, err = NewOtto(repo)
	if err == nil {
		_, _ = wb.RefreshGUID()
	}
	wb.Log("guid", wb.GUID, err)
	wb.Log("ua", ua)
	return
}

func (wb *WenshuBrowser) GetCookie(name string) string {
	return wb.c.GetCookie(wb.home.String(), name)
}

// :18071 是自定义的一个文书网反向代理，我们通过一个特殊的地址获取它记录下来的cookie
func (wb *WenshuBrowser) LoadCookies() (err error) {
	cookiesraw, sc, err := wb.c.GetTXT("http://localhost:18071/4bd54ad7a845/", "")
	if err != nil {
		return
	}
	var cookies [][2]string
	for _, ck := range strings.Split(string(cookiesraw), "\n") {
		if kv := strings.SplitN(ck, "=", 2); len(kv) == 2 {
			kv[0] = strings.TrimSpace(kv[0])
			kv[1] = strings.TrimSpace(kv[1])
			cookies = append(cookies, [2]string{kv[0], kv[1]})
		}
	}
	wb.c.UseCookie(wb.home.String(), cookies)

	log.Println("fetch-cookie", sc, len(cookies), err)
	return
}

// 逆向 sojson.v5没成功，所以我们通过无头chrome访问文书网的反向代理，通过代理记录wzws_cid和vjkl5
func (wb *WenshuBrowser) WZWSCID(uri string) (string, error) {
	var (
		result []byte
		err    error
		path   = strings.TrimPrefix(uri, "http://wenshu.court.gov.cn")
	)
	result, err = exec.Command(`google-chrome-unstable`, "--headless", "--dump-dom", `http://localhost:18071`+path).Output()
	if err == nil {
		err = wb.LoadCookies()
	}
	return string(result), err
}
func (wb *WenshuBrowser) Home() error {
	_, err := wb.WZWSCID(wb.home.String())
	return err
}

// SetCookies ...
func SetCookies(jar http.CookieJar, u *url.URL, namevalues ...string) {
	var cookies []*http.Cookie
	for i := 0; i < len(namevalues); i = i + 2 {
		n, v := namevalues[i], namevalues[i+1]
		cookies = append(cookies, &http.Cookie{
			Name:  n,
			Value: v,
		})
	}
	jar.SetCookies(u, cookies)
}

// NewOtto ...
// 准备文书网的JS环境，用到的函数在bootstrap.js中定义
func NewOtto(repo string) (jse *otto.Otto, err error) {
	jse = otto.New()
	var load = func(jse *otto.Otto, fn string) {
		var data []byte
		fn = filepath.Join(repo, fn)
		if err == nil {
			data, err = ioutil.ReadFile(fn)
		}
		if err == nil {
			_, err = jse.Run(string(data))
		}
		if err != nil {
			log.Println(fn, err)
		}
	}
	load(jse, "bootstrap.js")  // 导出的API
	load(jse, "sha1.js")       // sha1
	load(jse, "md5.js")        // md5
	load(jse, "base64.js")     // Base64
	load(jse, "ListExtend.js") // CrashVL5X(vjkl5 string)
	load(jse, "b64.js")        // Base64_unzip
	load(jse, "aes.js")        // aes
	load(jse, "rawdeflate.js") // unzip
	load(jse, "rawinflate.js") // zip
	load(jse, "pako.js")       // com.str.Decrypt
	return
}

// JSCall 调用js中的函数，目前所有返回值都转换卫字符串
// JSCall(nil, jse, "somefunc", p1, p2) => somefunc(p1, p2)->string
func JSCall(err error, jse *otto.Otto, function string, params ...interface{}) (string, error) {
	var xs string
	var x otto.Value
	if err == nil {
		x, err = jse.Call(function, nil, params...)
	}
	if err == nil {
		xs, err = x.ToString()
	}
	return xs, err
}

// 执行js脚本
// JSRun(nil, jse, `(function(){return "this is a anonymouse function"})()`)
func JSRun(err error, jse *otto.Otto, script string) (string, error) {
	var ret otto.Value
	if err == nil {
		ret, err = jse.Run(script)
	}
	if err == nil {
		return ret.ToString()
	}
	return "", err
}

// RefreshGUID ...
func (wb *WenshuBrowser) RefreshGUID() (guid string, err error) {
	guid, err = JSCall(err, wb.jse, "CrashGUID")
	wb.GUID = guid
	return guid, err
}

// 只要反向代理成功访问过文书网，代理会记住所有的cookie
// 程序需要vjkl5才能计算vl5x，这里的vjkl5是我们通过代理程序获得的
func (wb *WenshuBrowser) List(ajlx string) (cnt int64, ret []map[string]interface{}, err error) {
	var (
		q       = ants.Q
		param   = sparam("案件类型", ajlx)
		refer   = `http://127.0.0.1:18071/List/List?sorttype=1`
		uri     = `http://127.0.0.1:18071/List/ListContent`
		resp    []byte
		xs      string
		runeval string
		vl5x    string
		vjkl5   string
	)
	refer = ants.MakeURI(refer, q("conditions", "searchWord 1 AJLX  "+param))

	vjkl5 = wb.GetCookie("vjkl5") // List页面会设置cookie vjkl5
	vl5x, err = JSCall(err, wb.jse, "CrashVL5X", vjkl5)
	wb.Log("vjkl5 = ", vjkl5, ", vl5x = ", vl5x, err)

	if err == nil {
		resp, err = Submit(wb.c, uri, refer, "guid", "8e03feb9-3607-3ef93269-041032a56d4d",
			"vl5x", vl5x,
			"number", "loca",
			"Order", "法院层级",
			"Direction", "asc",
			"Index", "1",
			"Page", "10",
			"Param", param,
		)
	}
	fmt.Println(string(resp))
	xs, err = JSRun(err, wb.jse, string(resp))
	if err == nil {
		err = json.Unmarshal([]byte(xs), &ret)
	}
	if err == nil && len(ret) == 0 {
		err = fmt.Errorf("invalid ListContent")
	}
	if err == nil {
		runeval = ants.I2S(ret[0], "RunEval")
		cnt = ants.I2I(ret[0], "Count")
		ret = ret[1:]
	}
	wb.Log("ListContent", cnt, len(ret), err)
	// 文书网通过setTimeout更新com.str._KEY, 同时又在setTimeout中埋了一个循环
	// 所以setTimeout简单判断了下需要eval的内容
	_, err = JSCall(err, wb.jse, "CrashRunEval", runeval)

	for _, kase := range ret {
		did := ants.I2S(kase, "文书ID")
		kase["文书ID"], err = JSCall(err, wb.jse, "CrashDOCID", did)
	}
	xs, runeval, err = wb.CaseContent(ret[0]["文书ID"].(string))
	wb.Log(xs, err)
	wb.Log(runeval, err)
	return cnt, ret, err
}

// CaseContent ...
// google-chrome --headless --no-sandbox --dump-dom 'http://....?DocID=${id}&Keyword=' | goquery-extract ...
func (wb *WenshuBrowser) CaseContent(docid string) (summary, txt string, err error) {
	var (
		sc   int
		data []byte
		p1   = ants.MakeURI("http://wenshu.court.gov.cn/content/content?KeyWord=", ants.Q("DocID", docid))
		p2   = ants.MakeURI("http://wenshu.court.gov.cn/CreateContentJS/CreateContentJS.aspx", ants.Q("DocID", docid))
	)
	data, sc, err = wb.c.GetTXT(p1, wb.home.String())
	wb.Log("content/content", sc, err)
	if err == nil {
		data, sc, err = wb.c.GetTXT(p2, p1)
	}
	if err == nil {
		summary, txt, err = CrashCreateContentJS(data, wb.jse)
	}
	return
}

// CrashCreateContentJS ...
func CrashCreateContentJS(src []byte, jse *otto.Otto) (summary, txt string, err error) {
	re := regexp.MustCompile(`(?s:\$\(function\(\)\s*{(.*?)}\);)`)
	stmts := re.FindAllStringSubmatch(string(src), -1)
	// 应该有四个匿名function
	if len(stmts) != 4 {
		err = fmt.Errorf("extract functions from content.xhr failed")
		return
	}
	// 保留每个函数体, 然后改成返回内容的匿名函数调用
	re2 := regexp.MustCompile(`(?s:\$\("#\w+"\)\.\w+\(.*?\);)`)
	info := stmts[1][1]
	info = re2.ReplaceAllString(info, "")
	// fmt.Println(info)
	info = "(function(){" + info + "});\n return caseinfo;})()"
	x, err := jse.Run(info)
	if err == nil {
		summary, err = x.ToString()
	}

	body := stmts[3][1]
	body = re2.ReplaceAllString(body, "")
	body = strings.ReplaceAll(body, "Content.Content.InitPlugins();", "")
	body = strings.ReplaceAll(body, "Content.Content.KeyWordMarkRed();", "")
	body = "(function (){" + body + "\n return jsonHtml;})()"
	x, err = jse.Run(body)
	if err == nil {
		txt, err = x.ToString()
	}
	return
}

// Submit ...
func Submit(c *ants.Client, uri, refer string, kv ...string) ([]byte, error) {
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
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func form(kv ...string) io.Reader {
	data := params(kv...)
	return strings.NewReader(data.Encode())
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
