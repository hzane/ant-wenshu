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
	"encoding/json"
	"strconv"
)

func main() {
	guid := GUID()

	client := tools.NewHTTPClient2(time.Second*15, 2, nil, nil)
	client.RewriteRequest = appendHeader
	number := GetCode(client, guid)
	fmt.Println(guid, number)
	vjkl5 := VJKL5(client, guid, number)
	vl5x, err := vl5x(vjkl5)
	fmt.Println(vjkl5, vl5x, err)
	ListContent(client, vjkl5, vl5x, number, guid, 1, 5, "全文检索:农业科学院")
}

func appendHeader(req *http.Request) *http.Request {
	req.Header.Set("Origin", host)
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", host)
	}
	//	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	//	req.Header.Set("Accept-Encoding", "gzip, deflate")

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

type wenshu struct {
	guid   string
	number string
	vjkl5  string
	vl5x   string
}

/*
curl 'http://wenshu.court.gov.cn/list/list/?sorttype=1&number=BHDXJYU9&guid=ac02df2c-f81d-1aeaf27d-3b4da5454e8e&conditions=searchWord+QWJS+++全文检索:农业科学院'
-H 'Connection: keep-alive'
-H 'Pragma: no-cache'
-H 'Cache-Control: no-cache'
-H 'Upgrade-Insecure-Requests: 1'
-H 'DNT: 1'
-H 'User-Agent: Mozilla/...'
-H 'Accept: text/html,application/...'
-H 'Referer: http://wenshu.court.gov.cn/Index'
-H 'Accept-Encoding: gzip, deflate'
-H 'Accept-Language: en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6'
-H 'Cookie: _gscu_2116842793=3533706665kvkj18;
Hm_lvt_d2caefee2de09b8a6ea438d74fd98db2=1535337067,1535357887,1535591984;
_gscbrs_2116842793=1; ASP.NET_SessionId=hrixbudtagxscgszhqgofjtd;
vjkl5=bdef436f9aff6a8857019b181bde5a953144d58e;
Hm_lpvt_d2caefee2de09b8a6ea438d74fd98db2=1535593695;
_gscs_2116842793=35591987bihxht34|pv:6' --compressed*/
const listURL = `http://wenshu.court.gov.cn/list/list/?sorttype=1&number=%v&guid=%v&conditions=searchWord+QWJS+++%v`

// VJKL5 ...
func VJKL5(client *tools.Client, guid, number string) (string) {
	uri := fmt.Sprintf(listURL, number, guid, url.QueryEscape("全文检索:农业科学院"))

	req, _ := http.NewRequest("GET", uri, nil)
	//req.Header.Set("X-Requested-With", "XMLHttpRequest")

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
func ListContent(client *tools.Client, vjkl5, vl5x, number, guid string,
	index, page int,
	param string) {
	uri := "http://wenshu.court.gov.cn/List/ListContent"
	refer := fmt.Sprintf(listURL, number, guid, url.QueryEscape(param))
	body := url.Values{}
	body.Set("Index", strconv.Itoa(index))
	body.Set("Page", strconv.Itoa(page))
	body.Set("Order", "法院层级")
	body.Set("Direction", "asc")
	body.Set("vl5x", vl5x)
	body.Set("number", "&gui")
	body.Set("guid", guid)
	body.Set("Param", param)
	req, _ := http.NewRequest("POST", uri, bytes.NewBufferString(body.Encode()))
	req.Header.Set("Refer", refer)
	req.Header.Add("Cookie", "vjkl5="+vjkl5)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		c := DecodeListContent(string(b))
		log.Println(c)
	}
}

// DecodeListContent ...
func DecodeListContent(src string) string {
	vm := otto.New()
	val, err := vm.Run(src)
	if err == nil && val.IsString() {
		return val.String()
	}
	return ""
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

func listlistdemo() {
	x := `"[{\"RunEval\":\"w61aw5vCjsKCMBDDvRbCjA9tMMO7A8OEJz9hHycNMcOowq48wqzCmMKKT8OGf18KLsOLwqVywpHDkjR6EjLChHYuZ05nCk1YHsOiw53DvhzDicO4wpTCrj9TGR/Cvz/CvmTDssKzOWzDpSbDmcOtwpnDr8O5JCDCnBYrwpAAw7EeFcOyCsKJwrxJwrtiV8OCw5pCwoB3MAvCgcOVA2NYXMKQwoDDvMKRNcKYQMOqYAfCnMKADhADwp7CkBwSRkpYAQgUwoLCg8KCB0HCuF5Ew4nDscKcw4pLwpQmchHChBTCisOsYsOMwr/DnsKUWnldb8KcCkvDisKEED7Ds8OyCU51wp9qwqJ/Jn9cw77Dv8KdQEwNZxrDpMKpGRLClMOdGgHDq8Oqw57DncOgccKwBm7CrW4nwqbCqjfDvXgZKkfDncKmwqgFQcKbQ8KHw68Sw6QwRivCj8K9wrAHasKMURvCrcO7wpzDgQTCq8KpwqZGw6zDu8O7w6TCr8OrwoouMBvDmFrCjGbDrcKad243FcODBVAHwq/DtgjDkwU2U8O2FsK8wptkYTbCgi0Xw4rCjGQYb8ORwq7CnsK0w57CrcK2w4rDl8Oowq7DsMOgwrU+T1HCjMO8aMKwT8OtwoTCiMO2w6HCvEY1I8K2wrvCscOdAHBHw5HDnD0GHcKowpzDpEtzfmolR8K1wpfCl37CtMOqwpAmwpwAHSIPeMKKwojCnMOxw6AX\",\"Count\":\"2211\"},{\"裁判要旨段原文\":\"本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻\",\"案件类型\":\"2\",\"裁判日期\":\"2015-06-26\",\"案件名称\":\"山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书\",\"文书ID\":\"DcOOw4kRw4AwCATDgcKUOBYBTyEgw7/CkMOsAMKma8O2wrA2wr3CgzPCt8KDwpUkA8KXR29XJjTDlgNHV0wlw6rDocOKIcOSQ27Dng8hwqDDncKuDiU7w6zCvCbDi8OhYcO3D8KiF1nDvsOIwqnDuMOawr7CtEYRwrFTwrIWw5lkwokudgPDmT7DkcOzwpJQw4XDlDnDhsKRAlPDoCfCuSXCnsOVwqhfw6J/w6duw79OQsK6w4fDnhA6ZGrCpQ0vZmDDkcK8w6jCsx8=\",\"审判程序\":\"再审审查与审判监督\",\"案号\":\"（2014）民申字第1782号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，沁州黄公司于2011年6月23日向本院申请再审，并于2011年11月5日本院再审审查期间提出撤回再审申请，本院于2011年11月24日作出（2011）民申字第922号民事裁定，准许沁州黄公司撤回再审申请。鉴于沁州黄公司于2012年9月27日第二次申请\",\"案件类型\":\"2\",\"裁判日期\":\"2013-12-20\",\"案件名称\":\"再审申请人山西沁州黄小米（集团）有限公司与被申请人沁县吴阁老土特产有限公司确认不侵害商标权、侵害商标权纠纷再审审查民事裁定书\",\"文书ID\":\"DcKNwrcNw4BADMOEVlJ4wqVSccO/wpFsw6AKwoIFTwXCq8KiJ8KubTNzw7Iyw5dGYcK2G8OmHMKdwrnCn8KZwo7Dly8WZV3Ch2NawoULwqFIw7fDh04XeGLDpVoOKgLDrMOCw5fDsMOUw4PCiXfDk8KYLkpeVmDCoMOFUcOww5zDhsOfwo9NfHYkwoLCsQHCtcOJw4FaGH8NNcKNw4ANwqnDhTnDn8K/w7/CksOBVihec8Kyw7oQw7LCtsOow47DksK+QW0bS8K9w74A\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第1643号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为：龙茂公司损失发生后，新疆维吾尔自治区种子管理站（以下简称种子管理站）和新疆维吾尔自治区种子质量监督检验站（以下简称种子质检站）曾对“抗病86”、“97728”甜瓜种子质量进行了田间实地检测，但因种子补过异品种而无法鉴定。后经新疆维吾尔自治区产品质量监\",\"案件类型\":\"2\",\"裁判日期\":\"2013-08-15\",\"案件名称\":\"乌鲁木齐市龙茂实业有限公司与新疆农业科学院园艺作物研究所、新疆农科院园艺科技开发公司、黄再兴、佘建华财产损害赔偿纠纷申请再审民事裁定书\",\"文书ID\":\"FcKNQQJEMQTDhcKuRMOLw4MSw6XDvkfCmj/Cu2zCksOcworDjFlqwpXCuMKYwrLCl8KQZyIywqp3wqMvMXA4D0TCmcKaw4rCuMK6w6I8ZMKvwrZrwprDugsuRcKow6RcPsKnwrsdDcO7w4Bdwq4uMMKkRnbDtnFDw4vDsMK3DxA+wqc4IDDDs8K7JBRPwqzCosK3w7JtcsOXw7fCtF7Djyluw58tw6ZbwpAAW8Oqw7I0wrlRw7IswoReEMOVw7VBwptwfMKBZsO/AQ==\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第242号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"综上所述，本院认为，英巍合作社认为案涉高速公路因距离其养殖场过近、导致养殖场功能丧失必须搬迁的理由不能成立，不予支持，原审判决对此认定正确，应予维持。\\n关于英巍合作社养殖场的种猪发病死亡是否是由案涉高速公路导致的问题。就此问题，双方当事人分别提供了专家意见，\",\"案件类型\":\"2\",\"裁判日期\":\"2014-02-07\",\"案件名称\":\"辽宁英巍良种猪专业合作社与辽宁省高等级公路建设局相邻关系纠纷二审民事判决书\",\"文书ID\":\"FcKMwrcRBEEMw4NawpJWPsKUw63Cv8Kkwr8PQXAQw7XClgfCu0oUcmfCjnMjRcOYZjAubUBdw5bCtsO2wrjDmRDDsRXDl3VbwrbCglrChMK8w4tVDmnCgBQ4woTCoG8nwpADI8KGw7/CjcOhwrowwpVkTMOmwoYpwr84w7DDpB7DgcOywo5mfsOdwpoBcA5Mwoh9wp8Jc8KkwpXCsnXCr8KPBMKHw7ZEeQnDlT3DrsK1wpnDolvClXHDgcKbMcOvAMOYcX8uPw==\",\"审判程序\":\"二审\",\"案号\":\"（2013）民一终字第83号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院经审查认为：孙宝田诉《四平日报》社和《城市晚报》上级主管单位吉林日报社侵害名誉权民事诉讼案件，人民法院已经依法作出驳回诉讼请求的生效判决。在此之后，孙宝田就该事项又向吉林省人民政府申请行政复议，请求责令吉林省新闻出版局履行职责，责令《城市晚报》、《四平日报\",\"案件类型\":\"4\",\"裁判日期\":\"2015-07-28\",\"案件名称\":\"孙宝田与吉林省人民政府行政复议申诉行政裁定书\",\"文书ID\":\"DcKOwrkNw4BADMODVjrDv3bDqXfDv8KRwpJSAiEqwrd2w6DDvMKVY8K6DSrDmMOjw6Yaw6DCicOuwr16wrAaw74zwrTCj3Row4DCh8OoQjXCqnfDq8OuCkcvwoTDpTIvUBZxw6wfeS5Ww7vDh8KAEGDDn8ORPDFsScKLwrbDuEvDjxB6w70LEFjDucOpUzYRw6HDpcOqQsKVw4BLwp0+bkDDgRzDvm/DssKSwrTCqS7Dk2RNwqJMw5lICx7CkipmMMKEwq0rfg==\",\"案号\":\"（2015）行监字第32号\",\"法院名称\":\"最高人民法院\"}]"`

	x = DecodeListContent(x)
	var result []map[string]interface{}
	_ = json.Unmarshal([]byte(x), &result)
	bv, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(bv))
}

/*
[
  {
    "Count": "2211",
    "RunEval": "w61aw5vCjsKCMBDDvRbCjA9tMMO7A8OEJz9hHycNMcOowq48wqzCmMKKT8OGf18KLsOLwqVywpHDkjR6EjLChHYuZ05nCk1YHsOiw53DvhzDicO4wpTCrj9TGR/Cvz/CvmTDssKzOWzDpSbDmcOtwpnDr8O5JCDCnBYrwpAAw7EeFcOyCsKJwrxJwrtiV8OCw5pCwoB3MAvCgcOVA2NYXMKQwoDDvMKRNcKYQMOqYAfCnMKADhADwp7CkBwSRkpYAQgUwoLCg8KCB0HCuF5Ew4nDscKcw4pLwpQmchHChBTCisOsYsOMwr/DnsKUWnldb8KcCkvDisKEED7Ds8OyCU51wp9qwqJ/Jn9cw77Dv8KdQEwNZxrDpMKpGRLClMOdGgHDq8Oqw57DncOgccKwBm7CrW4nwqbCqjfDvXgZKkfDncKmwqgFQcKbQ8KHw68Sw6QwRivCj8K9wrAHasKMURvCrcO7wpzDgQTCq8KpwqZGw6zDu8O7w6TCr8OrwoouMBvDmFrCjGbDrcKad243FcODBVAHwq/DtgjDkwU2U8O2FsK8wptkYTbCgi0Xw4rCjGQYb8ORwq7CnsK0w57CrcK2w4rDl8Oowq7DsMOgwrU+T1HCjMO8aMKwT8OtwoTCiMO2w6HCvEY1I8K2wrvCscOdAHBHw5HDnD0GHcKowpzDpEtzfmolR8K1wpfCl37CtMOqwpAmwpwAHSIPeMKKwojCnMOxw6AX"
  },
  {
    "审判程序": "再审审查与审判监督",
    "文书ID": "DcOOw4kRw4AwCATDgcKUOBYBTyEgw7/CkMOsAMKma8O2wrA2wr3CgzPCt8KDwpUkA8KXR29XJjTDlgNHV0wlw6rDocOKIcOSQ27Dng8hwqDDncKuDiU7w6zCvCbDi8OhYcO3D8KiF1nDvsOIwqnDuMOawr7CtEYRwrFTwrIWw5lkwokudgPDmT7DkcOzwpJQw4XDlDnDhsKRAlPDoCfCuSXCnsOVwqhfw6J/w6duw79OQsK6w4fDnhA6ZGrCpQ0vZmDDkcK8w6jCsx8=",
    "案件名称": "山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书",
    "案件类型": "2",
    "案号": "（2014）民申字第1782号",
    "法院名称": "最高人民法院",
    "裁判日期": "2015-06-26",
    "裁判要旨段原文": "本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻"
  },
  {
    "审判程序": "再审",
    "文书ID": "DcKNwrcNw4BADMOEVlJ4wqVSccO/wpFsw6AKwoIFTwXCq8KiJ8KubTNzw7Iyw5dGYcK2G8OmHMKdwrnCn8KZwo7Dly8WZV3Ch2NawoULwqFIw7fDh04XeGLDpVoOKgLDrMOCw5fDsMOUw4PCiXfDk8KYLkpeVmDCoMOFUcOww5zDhsOfwo9NfHYkwoLCsQHCtcOJw4FaGH8NNcKNw4ANwqnDhTnDn8K/w7/CksOBVihec8Kyw7oQw7LCtsOow47DksK+QW0bS8K9w74A",
    "案件名称": "再审申请人山西沁州黄小米（集团）有限公司与被申请人沁县吴阁老土特产有限公司确认不侵害商标权、侵害商标权纠纷再审审查民事裁定书",
    "案件类型": "2",
    "案号": "（2013）民申字第1643号",
    "法院名称": "最高人民法院",
    "裁判日期": "2013-12-20",
    "裁判要旨段原文": "本院认为，沁州黄公司于2011年6月23日向本院申请再审，并于2011年11月5日本院再审审查期间提出撤回再审申请，本院于2011年11月24日作出（2011）民申字第922号民事裁定，准许沁州黄公司撤回再审申请。鉴于沁州黄公司于2012年9月27日第二次申请"
  },
  {
    "审判程序": "再审",
    "文书ID": "FcKNQQJEMQTDhcKuRMOLw4MSw6XDvkfCmj/Cu2zCksOcworDjFlqwpXCuMKYwrLCl8KQZyIywqp3wqMvMXA4D0TCmcKaw4rCuMK6w6I8ZMKvwrZrwprDugsuRcKow6RcPsKnwrsdDcO7w4Bdwq4uMMKkRnbDtnFDw4vDsMK3DxA+wqc4IDDDs8K7JBRPwqzCosK3w7JtcsOXw7fCtF7Djyluw58tw6ZbwpAAW8Oqw7I0wrlRw7IswoReEMOVw7VBwptwfMKBZsO/AQ==",
    "案件名称": "乌鲁木齐市龙茂实业有限公司与新疆农业科学院园艺作物研究所、新疆农科院园艺科技开发公司、黄再兴、佘建华财产损害赔偿纠纷申请再审民事裁定书",
    "案件类型": "2",
    "案号": "（2013）民申字第242号",
    "法院名称": "最高人民法院",
    "裁判日期": "2013-08-15",
    "裁判要旨段原文": "本院认为：龙茂公司损失发生后，新疆维吾尔自治区种子管理站（以下简称种子管理站）和新疆维吾尔自治区种子质量监督检验站（以下简称种子质检站）曾对“抗病86”、“97728”甜瓜种子质量进行了田间实地检测，但因种子补过异品种而无法鉴定。后经新疆维吾尔自治区产品质量监"
  },
  {
    "审判程序": "二审",
    "文书ID": "FcKMwrcRBEEMw4NawpJWPsKUw63Cv8Kkwr8PQXAQw7XClgfCu0oUcmfCjnMjRcOYZjAubUBdw5bCtsO2wrjDmRDDsRXDl3VbwrbCglrChMK8w4tVDmnCgBQ4woTCoG8nwpADI8KGw7/CjcOhwrowwpVkTMOmwoYpwr84w7DDpB7DgcOywo5mfsOdwpoBcA5Mwoh9wp8Jc8KkwpXCsnXCr8KPBMKHw7ZEeQnDlT3DrsK1wpnDolvClXHDgcKbMcOvAMOYcX8uPw==",
    "案件名称": "辽宁英巍良种猪专业合作社与辽宁省高等级公路建设局相邻关系纠纷二审民事判决书",
    "案件类型": "2",
    "案号": "（2013）民一终字第83号",
    "法院名称": "最高人民法院",
    "裁判日期": "2014-02-07",
    "裁判要旨段原文": "综上所述，本院认为，英巍合作社认为案涉高速公路因距离其养殖场过近、导致养殖场功能丧失必须搬迁的理由不能成立，不予支持，原审判决对此认定正确，应予维持。\n关于英巍合作社养殖场的种猪发病死亡是否是由案涉高速公路导致的问题。就此问题，双方当事人分别提供了专家意见，"
  },
  {
    "文书ID": "DcKOwrkNw4BADMODVjrDv3bDqXfDv8KRwpJSAiEqwrd2w6DDvMKVY8K6DSrDmMOjw6Yaw6DCicOuwr16wrAaw74zwrTCj3Row4DCh8OoQjXCqnfDq8OuCkcvwoTDpTIvUBZxw6wfeS5Ww7vDh8KAEGDDn8ORPDFsScKLwrbDuEvDjxB6w70LEFjDucOpUzYRw6HDpcOqQsKVw4BLwp0+bkDDgRzDvm/DssKSwrTCqS7Dk2RNwqJMw5lICx7CkipmMMKEwq0rfg==",
    "案件名称": "孙宝田与吉林省人民政府行政复议申诉行政裁定书",
    "案件类型": "4",
    "案号": "（2015）行监字第32号",
    "法院名称": "最高人民法院",
    "裁判日期": "2015-07-28",
    "裁判要旨段原文": "本院经审查认为：孙宝田诉《四平日报》社和《城市晚报》上级主管单位吉林日报社侵害名誉权民事诉讼案件，人民法院已经依法作出驳回诉讼请求的生效判决。在此之后，孙宝田就该事项又向吉林省人民政府申请行政复议，请求责令吉林省新闻出版局履行职责，责令《城市晚报》、《四平日报"
  }
]
*/
