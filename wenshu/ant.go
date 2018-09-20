package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type task struct {
	Params        string `json:"params,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusCode    int    `json:"status-code,omitempty"`
	CaseCount     int    `json:"case-count,omitempty"`
	ContentLength int64  `json:"content-length,omitempty"`
}

// ant ...
type ant struct {
	repo      string
	OKParams  sync.Map
	DOCIDs    sync.Map
	docs      chan map[string]interface{}
	okTasks   chan task
	failTasks chan task
	info      func(...interface{})
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
				a.OKParams.Store(t.Params, struct{}{})
			}
		}
		f.Close()
	}
	return nil
}

func (a *ant) Log(param string) {
	id := hex.EncodeToString(sha1.Sum([]byte(param))[:])
}

func (a *ant) Run() (func()) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(2)
	writeTask := func(fn string, pipe chan task) {
		defer wg.Done()
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		didfn := path.Join(a.repo, "doc.ids")
		docids, _ := os.OpenFile(didfn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer docids.Close()

		for {
			select {
			case doc := <-a.docs:
				if docid, _ := doc["_id"].(string); docid != "" {
					a.Save(docid, doc)
					fmt.Fprintln(docids, docid)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return func() {
		cancel()
		wg.Wait()
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
