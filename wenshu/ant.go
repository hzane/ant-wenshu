package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab.com/hearts.zhang/tools"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type task struct {
	ID            string `json:"id,omitempty"`
	Params        string `json:"params,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusCode    int    `json:"status-code,omitempty"`
	CaseCount     int    `json:"case-count,omitempty"`
	PageNo        int    `json:"page-no,omitempty"`
	PageSize      int    `json:"page-size,omitempty"`
	ContentLength int64  `json:"content-length,omitempty"`
	error         error
}

// ant ...
type ant struct {
	wg        sync.WaitGroup
	client    *tools.Client
	guid      string
	number    string
	repo      string
	OKParams  sync.Map
	DOCIDs    sync.Map
	items     map[string][]Parameter
	docs      chan map[string]interface{}
	okTasks   chan *task
	failTasks chan *task
	qTasks    chan *task
	info      func(...interface{})
}

func (a *ant) LoadTree(fn string) (err error) {
	tf, err := os.Open(fn)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(tf)
	for scanner.Scan() {
		if fields := strings.Fields(scanner.Text()); len(fields) == 3 {
			a.items[fields[0]] = append(a.items[fields[0]], Parameter{key: fields[1]})
		}
	}
	tf.Close()
	return
}

func (a *ant) Load() error {
	didsfn := path.Join(a.repo, "doc.ids")
	if f, err := os.Open(didsfn); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			a.DOCIDs.Store(scanner.Text(), struct{}{})
		}
		f.Close()
	}
	okfn := path.Join(a.repo, "tasks.ok")
	if f, err := os.Open(okfn); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var t task
			if err := json.Unmarshal(scanner.Bytes(), &t); err == nil {
				a.OKParams.Store(t.ID, struct{}{})
			}
		}
		f.Close()
	}
	return nil
}

// NewAnt ...
func NewAnt(repo string) (*ant, func()) {
	a := &ant{
		repo:      repo,
		guid:      GUID(),
		client:    tools.NewHTTPClient(),
		items:     map[string][]Parameter{},
		okTasks:   make(chan *task, 1e2),
		failTasks: make(chan *task, 1e2),
		qTasks:    make(chan *task, 1e6),
		docs:      make(chan map[string]interface{}, 1e2),
		info:      log.Println,
	}
	Home(a.client)
	Criminal(a.client) // 种上cookie
	a.number = GetCode(a.client, a.guid)
	vjkl5 := GetVJKL5FromCookie(a.client)
	a.info(a.guid, a.number, vjkl5)

	ctx, cancel := context.WithCancel(context.Background())
	os.MkdirAll(a.repo, 0755)

	a.wg.Add(2)
	writeTask := func(fn string, pipe chan *task) {
		defer a.wg.Done()
		mode, perm := os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)
		var f io.WriteCloser
		var err error
		if f, err = os.OpenFile(path.Join(a.repo, fn), mode, perm); err != nil {
			a.info(err)
			f = &NopWriter{}
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		for {
			select {
			case t := <-pipe:
				encoder.Encode(t)
			case <-ctx.Done():
				return
			}
		}
	}
	go writeTask("tasks.ok", a.okTasks)
	go writeTask("tasks.fail", a.failTasks)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		didfn := path.Join(a.repo, "doc.ids")
		docids, _ := os.OpenFile(didfn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer docids.Close()

		for {
			select {
			case doc := <-a.docs:
				if docid, _ := doc["_id"].(string); docid != "" {
					if _, loaded := a.DOCIDs.LoadOrStore(docid, struct{}{}); !loaded {
						a.Save(docid, doc)
						fmt.Fprintln(docids, docid)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			select {
			case t := <-a.qTasks:
				a.Do(t)
			case <-ctx.Done():
				return
			}
		}
	}()

	return a, cancel
}

func (a *ant) Bootstrap() {

	params := "案件类型:刑事案件"

	//裁判年份	文书类型	审判程序	法院地域	中级法院	基层法院
	//一级案由	二级案由	三级案由	关键词	法院层级

	for _, year := range a.items["裁判年份"] {
		params := params + ",裁判年份:" + year.key + ",文书类型:" + "判决书"

		for _, instance := range a.items["审判程序"] {
			params := params + ",审判程序:" + instance.key

			for _, high := range a.items["法院地域"] {
				t := newTask(params + ",法院层级:高级法院,法院地域:" + high.key)
				if _, ok := a.OKParams.Load(t.ID); !ok {
					a.qTasks <- t
				}
				time.Sleep(time.Second * 10)
				//causeExpand(items, params)
			}
			/*
			for _, intermediate := range a.items["中级法院"] {
				t := newTask(params + ",法院层级:中级法院,法院地域:" + intermediate.key)
				if _, ok := a.OKParams.Load(t.ID); !ok {
					a.qTasks <- t
				}
				//causeExpand(items, params)
			}
			for _, basic := range a.items["基层法院"] {
				t := newTask(params + ",法院层级:基层法院,法院地域:" + basic.key)
				if _, ok := a.OKParams.Load(t.ID); !ok {
					a.qTasks <- t
				}
			}*/
		}
	}
	a.info("bootstrap done")
}

func newTask(params string) *task {
	id := sha1.Sum([]byte(params))
	ret := &task{
		ID:       hex.EncodeToString(id[:]),
		Params:   params,
		PageNo:   1,
		PageSize: 5,
	}
	return ret
}
func (a *ant) Wait() {
	a.wg.Wait()
}

func (a *ant) Stop() {

}

func (a *ant) Do(t *task) {
	a.info(t.Params)
	_, cases, cnt, err := ListContent(a.client, a.number, a.guid, t.PageNo, t.PageSize, t.Params)
	t.CaseCount, t.error = cnt, err
	a.info(len(cases), cnt, err)

	if err != nil {
		a.failTasks <- t
	} else {
		a.okTasks <- t
	}
	for _, cise := range cases {
		a.docs <- cise
	}
	if cnt > 100 {
		a.causeExpand(a.items, t.Params)
		return
	}
	if t.CaseCount > t.PageNo*t.PageSize {
		t := &task{
			ID:       t.ID,
			Params:   t.Params,
			PageNo:   t.PageNo + 1,
			PageSize: t.PageSize,
		}
		a.qTasks <- t
	}
}

func (a *ant) causeExpand(items map[string][]Parameter, params string) {
	for _, cause := range items["一级案由"] {
		t := newTask(params + ",一级案由:" + cause.key)
		if _, ok := a.OKParams.Load(t.ID); !ok {
			a.qTasks <- t
		}
	}
	for _, cause := range items["二级案由"] {
		t := newTask(params + ",二级案由:" + cause.key)
		if _, ok := a.OKParams.Load(t.ID); !ok {
			a.qTasks <- t
		}
	}
	for _, cause := range items["三级案由"] {
		t := newTask(params + ",三级案由:" + cause.key)
		if _, ok := a.OKParams.Load(t.ID); !ok {
			a.qTasks <- t
		}
	}
	for _, cause := range items["关键词"] {
		t := newTask(params + ",关键词:" + cause.key)
		if _, ok := a.OKParams.Load(t.ID); !ok {
			a.qTasks <- t
		}
	}
}

func (a *ant) Save(docid string, doc map[string]interface{}) {
	parts := strings.Split(docid, "-")
	fp := path.Join(append([]string{a.repo, "summary"}, parts...)...) + ".json"
	err := os.MkdirAll(path.Dir(fp), 0755)
	if data, err := json.Marshal(doc); err == nil {
		err = ioutil.WriteFile(fp, data, 0644)
	}
	if err != nil {
		a.info(err)
	}
	return
}

type NopWriter struct{}

func (_ *NopWriter) Close() error {
	return nil
}
func (_ *NopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
