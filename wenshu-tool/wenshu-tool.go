package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	gq "github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	info = log.Println
	inf  = log.Printf
)

type (
	doc = map[string]interface{}
)

func main() {
	wb := NewWenshu(config.repo)
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

	courts, err := wb.CourtTreeContent("民事案件", "中级法院", "北京市第一中级人民法院")
	panice(err)
	printd(courts)
	// entries, err := wb.TreeList()
	// panice(err)
	// printd(entries)

	// courts, err = wb.GetChildAllCourt(courts)
	// panice(err)
	// printd(courts)

	// cases, err := wb.ListContent()
	// panice(err)
	// printd(cases)

	// entries, err = wb.TreeContent("案件类型", "刑事案件")
	// panice(err)
	// printd(entries)

	// dic, err := wb.GetDicValue()
	// fmt.Println(dic, err)

	// doc, err := wb.CreateContentJS("8252121f-8260-4241-b707-018d52d151ca")
	// panice(err)
	// printd([]map[string]interface{}{doc})
}
func init() {
	flag.StringVar(&config.repo, "repo", "/repo/wenshu", "")

	flag.Parse()
	_ = os.MkdirAll(config.repo, 0755)
}

var config struct {
	repo string
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
