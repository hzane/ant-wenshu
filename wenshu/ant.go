package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"gitlab.com/hearts.zhang/tools"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type task struct {
	ID            string `json:"id,omitempty"`
	Params        string `json:"params,omitempty"`
	Status        string `json:"status,omitempty"`
	Expanded      bool   `json:"expanded,omitempty"`
	StatusCode    int    `json:"status-code,omitempty"`
	CaseCount     int    `json:"case-count,omitempty"`
	PageNo        int    `json:"page-no,omitempty"`
	PageSize      int    `json:"page-size,omitempty"`
	ContentLength int64  `json:"content-length,omitempty"`
	error         error
}

type ant struct {
	client     *tools.Client
	number     string
	guid       string
	repo       string
	tree       string
	tasks      io.WriteCloser
	categories map[string][]Tag
	oks        map[string]*task
}

func AntAll(repo, treefn string) (err error) {
	a := &ant{
		client: tools.NewHTTPClient(),
		guid:   GUID(),
		repo:   repo,
		tree:   treefn,
	}
	Home(a.client)
	Criminal(a.client) // 种上cookie
	a.number = GetCode(a.client, a.guid)
	vjkl5 := GetVJKL5FromCookie(a.client)
	info(a.guid, a.number, vjkl5)

	a.categories, _ = LoadTree(treefn)
	a.oks, _ = LoadListContentOK(path.Join(repo, "tasks.all"))
	info(len(a.oks), "tasks loaded")

	a.tasks, _ = os.OpenFile(path.Join(repo, "tasks.all"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if a.tasks == nil {
		a.tasks = &NopWriter{}
	}

	params := "案件类型:刑事案件"

	// 裁判年份	文书类型	审判程序	法院地域	中级法院	基层法院
	// 一级案由	二级案由	三级案由	关键词	法院层级

	// 	for _, year := range a.categories["裁判年份"] {
	// 	params = params + ",文书类型:" + "判决书"

	// 	for _, instance := range a.categories["审判程序"] {
	// 		params := params + ",审判程序:" + instance.key

	for _, high := range a.categories["法院地域"] {
		t, _ := a.LoadNew(params+",法院层级:高级法院,法院地域:"+high.key, 1, config.pageSize, false)
		a.Query(t)
	}

	for _, intermediate := range a.categories["中级法院"] {
		t, _ := a.LoadNew(params+",法院层级:中级法院,中级法院:"+intermediate.key, 1, config.pageSize, false)
		a.Query(t)
	}
	for _, basic := range a.categories["基层法院"] {
		t, _ := a.LoadNew(params+",法院层级:基层法院,基层法院:"+basic.key, 1, config.pageSize, false)
		a.Query(t)
	}
	// 	}
	// 	}
	return
}

func LoadTree(fn string) (map[string][]Tag, error) {
	tf, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	items := map[string][]Tag{}

	scanner := bufio.NewScanner(tf)
	for scanner.Scan() {
		if fields := strings.Fields(scanner.Text()); len(fields) == 3 {
			items[fields[0]] = append(items[fields[0]], Tag{key: fields[1]})
		}
	}
	tf.Close()
	return items, err
}

func LoadListContentOK(fn string) (ids map[string]*task, err error) {
	ids = map[string]*task{}
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var t task
		if err := json.Unmarshal(scanner.Bytes(), &t); err == nil {
			if t.StatusCode == http.StatusOK {
				ids[t.ID] = &t
			}
		}
	}
	return
}

func LoadListDOCsOK(fn string) (ids map[string]struct{}, err error) {
	ids = map[string]struct{}{}
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ids[scanner.Text()] = struct{}{}
	}
	return
}

func id(params string, pn, psz int) string {
	id := sha1.Sum([]byte(params))
	return hex.EncodeToString(id[:]) + "-" + strconv.Itoa(pn) + "-" + strconv.Itoa(psz)
}

func (a *ant) LoadNew(params string, pn, psz int, expanded bool) (t *task, ok bool) {
	id := id(params, pn, psz)
	if t, ok = a.oks[id]; !ok {
		t = &task{
			ID:       id,
			Params:   params,
			PageNo:   pn,
			PageSize: psz,
			Expanded: expanded,
		}
		t.Expanded = t.Expanded || t.PageNo > 1
	}
	return
}

func (a *ant) Query(t *task) {
	info(t.Params, t.CaseCount, t.PageNo)
	var cases []map[string]interface{}
	var tries int
	for (t.StatusCode == 0 || t.StatusCode >= http.StatusInternalServerError) && tries < 10 {
		t.StatusCode, cases, t.CaseCount, t.error = ListContent(a.client, a.number, a.guid, t.PageNo, t.PageSize, t.Params)
		info("cases", len(cases), "total:", t.CaseCount, "page-no", t.PageNo, t.error)

		a.SaveCases(cases)
		json.NewEncoder(a.tasks).Encode(t)
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(4000)+1200))
		randSleep()
		tries++
		if t.StatusCode >= http.StatusInternalServerError {
			a.guid = GUID()
			a.number = GetCode(a.client, a.guid)
		}
		if t.StatusCode == http.StatusInternalServerError {
			time.Sleep(time.Second * time.Duration(config.mean) * 10)
			a.client = tools.NewHTTPClient()
			Criminal(a.client) // 种上cookie
		}
	}
	if t.PageNo > 1 {
		return
	}
	if t.CaseCount > 20*t.PageSize && !t.Expanded {
		a.causeExpand(t.Params)
		return
	}
	for i := 2; i <= (t.CaseCount+t.PageSize-1)/t.PageSize; i++ {
		t, _ = a.LoadNew(t.Params, i, t.PageSize, true)
		a.Query(t)
	}

	return
}

func (a *ant) causeExpand(params string) {
	for _, cause := range a.categories["一级案由"] {
		t, _ := a.LoadNew(params+",一级案由:"+cause.key, 1, config.pageSize, true)
		a.Query(t)
	}
	for _, cause := range a.categories["二级案由"] {
		t, _ := a.LoadNew(params+",二级案由:"+cause.key, 1, config.pageSize, true)
		a.Query(t)
	}
	for _, cause := range a.categories["三级案由"] {
		t, _ := a.LoadNew(params+",三级案由:"+cause.key, 1, config.pageSize, true)
		a.Query(t)
	}
	for _, cause := range a.categories["关键词"] {
		t, _ := a.LoadNew(params+",关键词:"+cause.key, 1, config.pageSize, true)
		a.Query(t)
	}
}

func (a *ant) SaveCases(docs []map[string]interface{}) {
	for _, doc := range docs {
		if docid, _ := doc["_id"].(string); docid != "" {
			parts := strings.Split(docid, "-")
			fp := path.Join(append([]string{a.repo, "summary"}, parts...)...) + ".json"
			_ = os.MkdirAll(path.Dir(fp), 0755)
			if data, err := json.Marshal(doc); err == nil {
				err = ioutil.WriteFile(fp, data, 0644)
				infoe(err, "save-case")
			}
		}
	}
}

type NopWriter struct{}

func (_ *NopWriter) Close() error {
	return nil
}
func (_ *NopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func infoe(err error, head string) {
	if err != nil {
		info(head, err)
	}
}
