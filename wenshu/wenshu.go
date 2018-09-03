package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract"

	"github.com/robertkrimen/otto"
	"gitlab.com/hearts.zhang/tools"
	// _ "github.com/robertkrimen/otto/underscore"
)

func mainx() {
	rand.Seed(time.Now().Unix())
	guid := GUID()
	client := tools.NewHTTPClient2(time.Second*15, 2, nil, nil)

	number := GetCode(client, guid)
	log.Println(guid, number)
	docids, _ := ListContent(client, number, guid, 1, 5, "全文检索:农业科学院")
	for _, docid := range docids {
		fmt.Println(docid)
	}
}

func treeListDemo() {
	vm := otto.New()
	x, _ := vmRunS(vm, treelist)
	x, _ = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, x))
	fmt.Println(x)
}

// http://wenshu.court.gov.cn/List/List?sorttype=1&conditions=searchWord+1+AJLX++案件类型:刑事案件
// 执行案件
// 赔偿案件
// 行政案件
// 民事案件
func main() {
	rand.Seed(time.Now().Unix())
	guid := GUID()
	client := tools.NewHTTPClient2(time.Second*15, 2, nil, nil)
	Home(client)
	Criminal(client)

	number := GetCode(client, guid)
	log.Println(guid, number)
	h, err := TreeList(client)
	_ = h

	var listKeywords = func(params string, items map[string][]filter) {
		for _, keyword := range items["关键词"] {
			_ = keyword
		}
	}
	_ = listKeywords

	params := "案件类型:刑事案件"
	items := map[string][]filter{}
	err = treeRoot(client, params, guid, number, items)
	err = criminalCauseExpand(client, params, guid, number, items)
	err = courtExpand(client, params, guid, number, items)
	_ = err

	for _, year := range items["裁判年份"] {
		params := params + ",裁判年份:" + year.key
		for _, typi := range items["文书类型"] {
			params := params + ",文书类型:" + typi.key
			ids, err := ListContent(client, number, guid, 1, 20, params+",法院层级:最高法院")

			log.Println("list-content", params, len(ids), err)

			for _, instance := range items["审判程序"] {
				params := params + ",审判程序:" + instance.key

				for _, high := range items["法院地域"] {
					params := params + ",法院层级:高级法院,法院地域:" + high.key
					ids, err = ListContent(client, number, guid, 1, 20, params)
					log.Println("list-c", params, len(ids), err)
				}
				for _, intermediate := range items["中级法院"] {
					_ = intermediate
					// params := params + ",法院层级:中级法院,法院地域:" + intermediate.key
				}
				for _, basic := range items["基层法院"] {
					_ = basic
					// params := params + ",法院层级:基层法院,法院地域:" + basic.key

				}

			}
		}
	}
}

type filter struct {
	key string
	cnt int
}

func extract(item map[string]interface{}, ret map[string][]filter) {
	f, _ := item["Field"].(string)
	k, _ := item["Key"].(string)
	v, _ := item["IntValue"].(float64)
	if f != "" && k != "" && v != 0 {
		ret[f] = append(ret[f], filter{k, int(v)})
		log.Println(f, k, v)
	}
	children, _ := item["Child"].([]interface{})
	for _, child := range children {
		extract(child.(map[string]interface{}), ret)
	}
}
func verb(method, uri string, resp *http.Response, err error) {
	var status string
	var code int
	var cl int64
	if resp != nil {
		status, code, cl = resp.Status, resp.StatusCode, resp.ContentLength
	}
	log.Println(cl, code, status, err, method, uri)
}
func treeExpand(client *tools.Client, uri, params, parval, guid, number string, ret map[string][]filter) (err error) {
	body := url.Values{}
	body.Set("Param", params)
	body.Set("parval", parval)

	// uri := `http://wenshu.court.gov.cn/List/TreeContent`
	req, _ := http.NewRequest("POST", uri, bytes.NewBufferString(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	verb("POST", uri, resp, err)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := JSONStringBody(resp.Body)

	var items []map[string]interface{}
	err = json.Unmarshal([]byte(content), &items)

	for _, item := range items {
		extract(item, ret)
	}
	return
}

func courtExpand(client *tools.Client, params, guid, number string, ret map[string][]filter) error {
	// uri := `https://wenshu.court.gov.cn/List/CourtTreeContent`
	for _, f := range ret["法院地域"] {
		t := map[string][]filter{}
		_ = treeExpand(client, CourtTreeContentURL, params+",法院地域:"+f.key, f.key, guid, number, t)
		for _, f := range t["中级法院"] {
			ret["中级法院"] = append(ret["中级法院"], f)
			//			_ = treeExpand(client, CourtTreeContentURL, params+",中级法院:"+f.key, f.key, guid, number, ret)
		}

	}
	return nil
}

func criminalCauseExpand(client *tools.Client, params, guid, number string, ret map[string][]filter) error {
	// uri := `https://wenshu.court.gov.cn/List/ReasonTreeContent`
	err := treeExpand(client, ReasonTreeContentURL, params+",一级案由:刑事案由", "刑事案由", guid, number, ret)
	for _, filter := range ret["二级案由"] {
		_ = treeExpand(client, ReasonTreeContentURL, params+",二级案由:"+filter.key, filter.key, guid, number, ret)
	}
	return err
}

// treeRoot ...
func treeRoot(client *tools.Client, params string, guid, number string, ret map[string][]filter) (err error) {
	// uri := `http://wenshu.court.gov.cn/List/TreeContent`
	body := url.Values{}
	body.Set("Param", params)
	body.Set("vl5x", VL5X(client))
	body.Set("guid", guid)
	body.Set("number", number)

	// uri := `http://wenshu.court.gov.cn/List/TreeContent`
	req, _ := http.NewRequest("POST", TreeContentURL, bytes.NewBufferString(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	verb("POST", TreeContentURL, resp, err)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := JSONStringBody(resp.Body)

	var items []map[string]interface{}
	err = json.Unmarshal([]byte(content), &items)

	for _, item := range items {
		extract(item, ret)
	}
	return
}

// TreeList ...
func TreeList(client *tools.Client) (ret string, err error) {
	// uri := `http://wenshu.court.gov.cn/List/TreeList`
	req, _ := http.NewRequest("POST", TreeListURL, nil)
	resp, err := client.Do(req)
	verb("POST", TreeListURL, resp, err)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return JSONStringBody(resp.Body)
}

// JSONStringBody ...
func JSONStringBody(r io.Reader) (ret string, err error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	vm := otto.New()
	ret, err = vmRunS(vm, string(body))
	if err != nil {
		return
	}
	ret, err = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, ret)) // -> json string
	return
}

// VL5X ...
func VL5X(client *tools.Client) string {
	vjkl5 := GetVJKL5FromCookie(client)
	ret, _ := vl5x(vjkl5)
	return ret
}

// GetVJKL5FromCookie ...
func GetVJKL5FromCookie(client *tools.Client) string {
	uri, _ := url.Parse(host)
	cookies := client.Jar.Cookies(uri)
	for _, ck := range cookies {
		if ck.Name == "vjkl5" {
			return ck.Value
		}
	}
	return ""
}

// Criminal ...
func Criminal(client *tools.Client) {
	if resp, err := client.Get(CriminalURL, host); err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	vjkl5 := GetVJKL5FromCookie(client)
	log.Println(vjkl5)
}

// Home ...
func Home(client *tools.Client) {
	// for set-cookie
	if resp, err := client.Get(host, ""); err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}
func caseExportDemo() {
	x := casecontent
	r1 := regexp.MustCompile(`stringify\((\{.*?\})\);`)
	caseinfo := r1.FindStringSubmatch(x)[1] // caseinfo is javascript json object, but we treat it as a json string
	r1 = regexp.MustCompile(`\\/Date\((\d+)\)`)
	caseinfo = r1.ReplaceAllString(caseinfo, `$1`)
	var ccase map[string]interface{}
	_ = json.Unmarshal([]byte(caseinfo), &ccase)
	// fmt.Println(ccase)

	vm := otto.New()
	r1 = regexp.MustCompile(`dirData\s=\s(\{.*?\});if`)
	caseinfo = r1.FindStringSubmatch(x)[1]
	caseinfo, _ = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, caseinfo))
	_ = json.Unmarshal([]byte(caseinfo), &ccase)
	// fmt.Println(ccase)

	r1 = regexp.MustCompile(`jsonHtmlData\s?=\s?("\{.*\}");`)
	caseinfo = r1.FindStringSubmatch(x)[1]
	caseinfo, _ = vmRunS(vm, caseinfo)
	fmt.Println(prettys(caseinfo))
	// fmt.Println(caseinfo)
	// _ = json.Unmarshal([]byte(caseinfo), &ccase)

	// fmt.Println(pretty(ccase))
}

func prettys(js string) string {
	var ret bytes.Buffer
	_ = json.Indent(&ret, []byte(js), "", "  ")
	return ret.String()
}

func pretty(jo interface{}) string {
	b, _ := json.MarshalIndent(jo, "", "  ")
	return string(b)
}

// CaseContent ...
// http://wenshu.court.gov.cn/content/content?DocID=d8952be5-e5a2-4b8b-b554-cccf5824617f&KeyWord=%E5%86
func CaseContent(docid string) {

}

// AESKey ...
// https://www.jianshu.com/p/1dc99e3d927c
func AESKey(runeval string) (key string, err error) {
	vm := otto.New()
	compile(vm, config.js, "docid.js")
	js, err := vm.Run(fmt.Sprintf(`GetJs("%v")`, runeval))
	if err != nil {
		return
	}
	jss, err := js.ToString()
	if err != nil {
		return
	}
	statements := strings.Split(jss, ";;")
	statements[0] = statements[0] + ";" // $hidescript=...
	js, err = vm.Run(statements[0])     // Tm('._KEY="6942871305;,*Mh)
	if err != nil {
		return
	}
	jss, err = js.ToString()
	if err != nil {
		return
	}
	log.Println(jss)

	r := regexp.MustCompile(`_\[_\]\[_\]\((.*?)\)\(\);`)
	xs := r.FindStringSubmatch(statements[1])[1]
	xs = strings.Replace(xs, "$hidescript", strconv.Quote(jss), -1)
	js, err = vm.Run(xs)
	if err != nil {
		return
	}
	jss, err = js.ToString()
	if err != nil {
		return
	}
	// setTimeout('com.str._KEY="a69e42871c4f499c930c755edbf6d7d1";',8000*Math.random());
	r = regexp.MustCompile(`_KEY="(.*)";'`)
	key = r.FindStringSubmatch(jss)[1]
	return
}

const host = "http://wenshu.court.gov.cn"

// ...
var (
	CriminalURL          = host + "/List/List?sorttype=1&conditions=searchWord+1+AJLX++案件类型:刑事案件"
	GetCodeURL           = host + "/ValiCode/GetCode"
	TreeListURL          = host + "/List/TreeList"
	ListContentURL       = host + "/List/ListContent"
	TreeContentURL       = host + "/List/TreeContent"
	ReasonTreeContentURL = host + "/List/ReasonTreeContent"
	CourtTreeContentURL  = host + "/List/CourtTreeContent"
	ValidateCodeURL      = host + "/User/ValidateCode"
	CheckCodeURL         = host + "/Content/CheckVisitCode"
	VisitRemindURL       = host + "/Html_Pages/VisitRemind.html"
)

// GetCode ...
func GetCode(client *tools.Client, guid string) (number string) {
	data := url.Values{}
	data.Set("guid", guid)

	req, _ := http.NewRequest("POST", GetCodeURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Origin", host)
	req.Header.Set("Referer", host)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
func GUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	return fmt.Sprintf("%x-%x-%x%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// CaseSummary ...
type CaseSummary struct {
	ID       string `json:"_id,omitempty"`
	Name     string `json:"案件名称"`
	CaseType string `json:"案件类型"`
	No       string `json:"案号"`
	Court    string `json:"法院名称"`
	Date     string `json:"裁判日期"`
	Abstract string `json:"裁判要旨段原文"`
}

// ListContent ...
/*
curl 'http://wenshu.court.gov.cn/List/ListContent'
-H 'Pragma: no-cache'
-H 'Origin: http://wenshu.court.gov.cn'
-H 'Accept-Encoding: gzip, deflate'
-H 'Accept-Language: en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6'
-H 'User-Agent: Mozilla/5.0 ...'
-H 'Cache-Control: no-cache'
-H 'X-Requested-With: XMLHttpRequest'
-H 'Cookie: _gscu_2116842793=...; vjkl5=c3c5bc9aff9f886c014b188efe53fc26b16f626e; ...'
-H 'Connection: keep-alive'
-H 'Referer: http://wenshu.court.gov.cn/list/list/?sorttype=1&number=&guid=042...0&conditions=searchWord+QWJS+++全文检索:农业科学院'
-H 'DNT: 1'
--data 'Param=全文检索:农业科学院&
Index=1&
Page=5&
Order=法院层级&
Direction=asc&
vl5x=4ce429d14932c99fd594b7e9&
number=%26gui&
guid=8bcbcecd-25f9-5922503e-d48918ba0c39' --compressed
*/
func ListContent(client *tools.Client, number, guid string,
	index, page int,
	param string) (ids []string, err error) {

	// uri := "http://wenshu.court.gov.cn/List/ListContent"
	// refer := fmt.Sprintf(listURL, number, guid, url.QueryEscape(param))
	body := url.Values{}
	body.Set("Index", strconv.Itoa(index))
	body.Set("Page", strconv.Itoa(page))
	body.Set("Order", "法院层级")
	body.Set("Direction", "asc")
	body.Set("vl5x", VL5X(client))
	body.Set("number", number)
	body.Set("guid", guid)
	body.Set("Param", param)
	req, _ := http.NewRequest("POST", ListContentURL, bytes.NewBufferString(body.Encode()))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	verb("POST", ListContentURL, resp, err)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	vm := otto.New()
	c, err := vmRunS(vm, string(b)) // javascript 字符串
	if err != nil {
		return
	}

	if strings.HasPrefix(c, "remind") {
		err = errors.New(c)
		ValidateCode(client)
		return
	}

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(c), &result)
	if err != nil {
		return
	}
	cnt, _ := result[0]["Count"].(string)
	runeval, _ := result[0]["RunEval"].(string)
	key, err := AESKey(runeval)
	if err != nil {
		return
	}
	log.Println(cnt, key)
	compile(vm, "docid.js")
	for _, doc := range result[1:] {
		id, _ := doc["文书ID"].(string)
		s, _ := vm.Run(fmt.Sprintf(`DecryptDocID("%v","%v");`, key, id))
		id, _ = s.ToString()
		doc["_id"] = id
		ids = append(ids, id)
		x, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(x))
	}
	return
}

func vmRunS(vm *otto.Otto, src string) (string, error) {
	val, err := vm.Run(src)
	if err == nil {
		return val.ToString()
	}
	return "", err
}

func vl5x(vjkl5 string) (string, error) {
	vm := otto.New()
	compile(vm, path.Join(config.js, "vl5x.js"))

	value, err := vm.Run(fmt.Sprintf(`GetVl5x("%v")`, vjkl5))
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

// ValidateCode ...
func ValidateCode(client *tools.Client) (err error) {
	//req, _ := http.NewRequest("POST", ValidateCodeURL, )
	resp, err := client.Get(ValidateCodeURL, "")
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}
	log.Println(resp.ContentLength, resp.StatusCode, resp.Header.Get("content-type"))

	ocr := gosseract.NewClient()
	defer ocr.Close()
	ocr.SetWhitelist("0123456789")
	err = ocr.SetImageFromBytes(body)
	if err != nil {
		return
	}

	text, err := ocr.Text()
	if err != nil {
		return
	}
	code := url.Values{
		"ValidateCode": []string{text},
	}
	//Html_Pages/VisitRemind.html
	req, _ := http.NewRequest("POST", CheckCodeURL, bytes.NewBufferString(code.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", VisitRemindURL)

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	log.Println("validate-code", resp.StatusCode, resp.Status, text)
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return nil
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
