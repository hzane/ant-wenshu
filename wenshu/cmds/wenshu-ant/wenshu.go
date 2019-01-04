package wenshu_ant

import (
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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
	"gitlab.com/hearts.zhang/tools"
)



func main() {
	if config.ant {
		AntAll(config.repo, config.tree)
		return
	}
	client, guid := tools.NewHTTPClient(), GUID()
	Home(client)
	Criminal(client) // 种上cookie

	if config.content {
		AntContent(client, config.repo)
	}

	if config.showCookie {
		vjkl5 := GetVJKL5FromCookie(client)
		fmt.Println(vjkl5)
	}

	number := GetCode(client, guid)
	if config.showCode {
		fmt.Println(guid, number)
	}

	if config.params != "" {
		_, cases, cnt, _ := ListContent(client, number, guid, config.pageNo, config.pageSize, config.params)
		for _, cese := range cases {
			json.NewEncoder(os.Stdout).Encode(cese)
		}
		info("total count", cnt)
	}

	if config.caseID != "" {
		doc, _ := CaseContent(client, config.caseID)
		json.NewEncoder(os.Stdout).Encode(doc)
	}
	if config.createTree {
		createTree(client, guid, number)
	}
}

// http://wenshu.court.gov.cn/List/List?sorttype=1&conditions=searchWord+1+AJLX++案件类型:刑事案件
// 执行案件
// 赔偿案件
// 行政案件
// 民事案件
func createTree(client *tools.Client, guid, number string) {
	Home(client)
	Criminal(client)

	_, err := TreeList(client)
	infoe(err, "tree-list")

	params := "案件类型:刑事案件"
	items := map[string][]Tag{}
	err = treeRoot(client, params, guid, number, items)
	infoe(err, "tree-root")

	err = criminalCauseExpand(client, params, guid, number, items)
	infoe(err, "cause-expand")

	err = courtExpand(client, params, guid, number, items)
	infoe(err, "court-expand")

	for key, items := range items {
		for _, item := range items {
			fmt.Println(key, item.key, item.cnt)
		}
	}

}


// 列表的嵌套结构
func treeItemConvert(item map[string]interface{}, ret map[string][]Tag) {
	f, _ := item["Field"].(string)     // 审判程序...
	k, _ := item["Key"].(string)       // 民事案件
	v, _ := item["IntValue"].(float64) // 案件数量
	if f != "" && k != "" && v != 0 {
		ret[f] = append(ret[f], Tag{k, int(v)})
	}
	children, _ := item["Child"].([]interface{})
	for _, child := range children {
		treeItemConvert(child.(map[string]interface{}), ret)
	}
}

func treeExpand(client *tools.Client,
	uri string,
	params string,
	parval string,
	guid, number string,
	ret map[string][]Tag) (err error) {

	resp, err := Submit(client, uri, set("Param", params), set("parval", parval))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := JSONStringBody(resp.Body)

	var items []map[string]interface{}
	err = json.Unmarshal([]byte(content), &items)

	for _, item := range items {
		treeItemConvert(item, ret)
	}
	return
}

// 高级法院 +中级法院 +基层法院
func courtExpand(client *tools.Client, params, guid, number string, ret map[string][]Tag) error {
	for _, f := range ret["法院地域"] {
		t := map[string][]Tag{}
		_ = treeExpand(client, CourtTreeContentURL, params+",法院地域:"+f.key, f.key, guid, number, t)
		for _, f := range t["中级法院"] {
			ret["中级法院"] = append(ret["中级法院"], f)
			_ = treeExpand(client, CourtTreeContentURL, params+",中级法院:"+f.key, f.key, guid, number, ret)
			randSleep()
		}

	}
	return nil
}

// 一级案由 +二级案由 +三级案由
func criminalCauseExpand(client *tools.Client, params, guid, number string, ret map[string][]Tag) error {
	err := treeExpand(client, ReasonTreeContentURL, params+",一级案由:刑事案由", "刑事案由", guid, number, ret)
	for _, filter := range ret["二级案由"] {
		_ = treeExpand(client, ReasonTreeContentURL, params+",二级案由:"+filter.key, filter.key, guid, number, ret)
	}
	return err
}

// treeRoot ...
// 初次请求检索树能够返回一个大纲
// 将嵌套结构展开成平铺结构
func treeRoot(client *tools.Client, params string, guid, number string, ret map[string][]Tag) (err error) {
	resp, err := Submit(client, TreeContentURL,
		set("Param", params),
		set("vl5x", VL5X(client)),
		set("guid", guid),
		set("number", number), )
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, _ := JSONStringBody(resp.Body)

	var items []map[string]interface{}
	err = json.Unmarshal([]byte(content), &items)

	for _, item := range items {
		treeItemConvert(item, ret)
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

func CaseDetail(raw string) (doc map[string]interface{}, err error) {
	vm := otto.New()

	r1 := regexp.MustCompile(`JSON.stringify\((\{.*?\})\);`)
	result := r1.FindStringSubmatch(raw)
	if len(result) == 0 {
		return
	}
	tmp := result[0] // tmp is javascript json object, but we treat it as a json string
	tmp, err = vmRunS(vm, tmp)
	err = json.Unmarshal([]byte(tmp), &doc)

	r1 = regexp.MustCompile(`dirData\s?=\s?(\{.*?\});if`)
	tmp = r1.FindStringSubmatch(raw)[1]
	tmp, err = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, tmp))
	err = json.Unmarshal([]byte(tmp), &doc)

	r1 = regexp.MustCompile(`jsonHtmlData\s?=\s?("\{.*\}");`)
	tmp = r1.FindStringSubmatch(raw)[1]
	tmp, err = vmRunS(vm, "("+tmp+")")

	err = json.Unmarshal([]byte(tmp), &doc)
	return
}



// CaseContent ...
// http://wenshu.court.gov.cn/content/content?DocID=d8952be5-e5a2-4b8b-b554-cccf5824617f&KeyWord=%E5%86
// http://wenshu.court.gov.cn/CreateContentJS/CreateContentJS.aspx?DocID=d8952be5-e5a2-4b8b-b554-cccf5824617f
func CaseContent(client *tools.Client, docID string) (summary map[string]interface{}, err error) {
	uri, _ := url.Parse(CreateContentJSURL)
	params := uri.Query()
	params.Set("DocID", docID)
	uri.RawQuery = params.Encode()

	ref, _ := url.Parse(ContentURL)
	ref.RawQuery = params.Encode()

	resp, err := client.Get(uri.String(), ref.String())

	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	summary, err = CaseDetail(string(body))
	return
}

// ListContent ...
// 判决书列表
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
	param string) (sc int, cases []map[string]interface{}, cnt int, err error) {

	resp, err := Submit(client, ListContentURL, set("Index", strconv.Itoa(index)),
		set("Page", strconv.Itoa(page)),
		set("Order", "法院层级"),
		set("Direction", "asc"),
		set("vl5x", VL5X(client)),
		set("number", number),
		set("guid", guid),
		set("Param", param), )

	if err != nil {
		return
	}
	defer resp.Body.Close()
	sc = resp.StatusCode

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
		return
	}
	vm := otto.New()
	c, err := vmRunS(vm, string(b)) // javascript 字符串
	if err != nil {
		sc = http.StatusInternalServerError
		return
	}

	if strings.HasPrefix(c, "remind") {
		sc = http.StatusInternalServerError
		err = errors.New(c)
		ValidateCode(client)
		return
	}

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(c), &result)
	if err != nil {
		c = strings.Replace(c, `\"`, `"`, -1)
		err = json.Unmarshal([]byte(c), &result)
	}
	if err != nil {
		info(c)
		sc = http.StatusInternalServerError
	}
	if len(result) == 0 {
		sc = http.StatusInternalServerError
		return
	}
	scnt, _ := result[0]["Count"].(string)
	cnt, _ = strconv.Atoi(scnt)

	runeval, _ := result[0]["RunEval"].(string)
	key, err := AESKey(runeval)
	if err != nil {
		sc = http.StatusInternalServerError
		return
	}
	info("aes-key", key, err)

	compile(vm, "docid.js")
	for _, doc := range result[1:] {
		id, _ := doc["文书ID"].(string)
		s, _ := vm.Run(fmt.Sprintf(`DecryptDocID("%v","%v");`, key, id))
		id, _ = s.ToString()
		doc["_id"] = id
		delete(doc, "文书ID")
	}
	cases = result[1:]
	return
}

func init() {
	rand.Seed(time.Now().Unix())

	flag.BoolVar(&config.showCookie, "show-cookie", false, "")
	flag.BoolVar(&config.showCode, "show-code", false, "")
	flag.BoolVar(&config.createTree, "create-tree", false, "create full tree")
	flag.BoolVar(&config.createParams, "create-params", false, "")
	flag.BoolVar(&config.ant, "ant", false, "")
	flag.BoolVar(&config.content, "ant-content", false, "")
	flag.StringVar(&config.params, "params", "", "list content with params")
	flag.StringVar(&config.caseID, "case-id", "", "show case details with id")
	flag.StringVar(&config.js, "js-dir", ".", "javascript file folder")
	flag.StringVar(&config.repo, "repo", "data", "")
	flag.StringVar(&config.proxies, "proxies", "", "")
	flag.StringVar(&config.tree, "tree", "trees.csv", "")
	flag.IntVar(&config.workers, "workers", 1, "")
	flag.IntVar(&config.pageNo, "page-no", 1, "")
	flag.IntVar(&config.pageSize, "page-size", 15, "")
	flag.Float64Var(&config.stddev, "stddev", 2000.0, "")
	flag.Float64Var(&config.mean, "mean", 1200.0, "")

	flag.Parse()
	config.guid = GUID()
}

var config struct {
	js           string
	repo         string
	proxies      string
	tree         string
	params       string
	caseID       string
	guid         string
	code         string
	workers      int
	pageSize     int
	pageNo       int
	showCookie   bool
	showCode     bool
	createTree   bool
	createParams bool
	ant          bool
	content      bool
	stddev       float64
	mean         float64
}

func randSleep() {
	x := rand.NormFloat64()*config.stddev + config.mean
	if x < 500 {
		x = 500
	}
	if x > config.mean+2*config.stddev {
		x = config.mean + 2*config.stddev
	}
	time.Sleep(time.Millisecond * time.Duration(x))
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