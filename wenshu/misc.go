package wenshu

import (
	"bytes"
	"encoding/json"
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

	"github.com/robertkrimen/otto"
	"gitlab.com/hearts.zhang/tools"
)

var info = log.Println

func verb(method, uri string, resp *http.Response, err error) {
	var status string
	var cl int64
	if resp != nil {
		status, cl = resp.Status, resp.ContentLength
	}
	info(cl, status, err, method, uri)
}

// SubmitForm ...
func SubmitForm(client *tools.Client, uri string, body url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("POST", uri, bytes.NewBufferString(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	verb("POST", uri, resp, err)
	return resp, err
}

type opt = func(url.Values)

func set(k, v string) opt {
	return func(body url.Values) {
		body.Set(k, v)
	}
}
func Submit(client *tools.Client, uri string, params ...opt) (*http.Response, error) {
	body := url.Values{}
	for _, set := range params {
		set(body)
	}
	return SubmitForm(client, uri, body)
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

func pretty(jo interface{}) string {
	b, _ := json.MarshalIndent(jo, "", "  ")
	return string(b)
}

func prettys(js string) string {
	var ret bytes.Buffer
	_ = json.Indent(&ret, []byte(js), "", "  ")
	return ret.String()
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

// GUID ...
func GUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	return fmt.Sprintf("%x-%x-%x%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
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
