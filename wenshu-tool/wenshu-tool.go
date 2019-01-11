package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	_ = wb.QueryReferer(config.word, config.ay, config.fycj, config.spcx, config.ajlx, config.wslx, config.slfy, config.cprq)
	err = wb.ListList(config.word, config.ay, config.fycj, config.spcx, config.ajlx, config.wslx, config.slfy, config.cprq)
	panice(err)

	v := wb.GetCookie("vjkl5")
	info("vjkl5\t: ", v)
	_, _ = wb.js.Run(fmt.Sprintf(`function getCookie(cname) {
    return "%s"
  }`, v))
	vlx, err := wb.VL5X()
	panice(err)
	info("vl5x :\t", vlx)

	n, err := wb.GetCode()
	panice(err)
	info("number :\t", n)

	switch config.method {
	case "tree-list":
		entries, err := wb.TreeList()
		panice(err)
		printd(entries)
	case "child-all-court":
		// courts, err := wb.GetChildAllCourt(courts)
		// panice(err)
		// printd(courts)
	case "list-content":
		cases, err := wb.ListContent(config.index,
			config.order,
			config.direction,
			config.ay,
			config.fycj,
			config.spcx,
			config.ajlx,
			config.wslx,
			config.slfy,
			config.cprq)
		panice(err)
		printd(cases)
	case "tree-content":
		entries, err := wb.TreeContent(config.ay, config.fycj, config.spcx, config.ajlx, config.wslx, config.slfy, config.cprq)
		panice(err)
		printd(entries)
	case "dic-value":
		dic, err := wb.GetDicValue(config.dicType, config.dicKeys)
		fmt.Println(dic, err)
	case "dic":
		doc, err := wb.CreateContentJS(config.docid)
		panice(err)
		printd([]map[string]interface{}{doc})
	case "court-tree-content":
		courts, err := wb.CourtTreeContent("民事案件", "中级法院", "北京市第一中级人民法院")
		panice(err)
		printd(courts)

	}

}
func init() {
	flag.StringVar(&config.repo, "repo", "/repo/wenshu", "")
	flag.StringVar(&config.docid, "doc-id", "8252121f-8260-4241-b707-018d52d151ca", "")
	flag.StringVar(&config.dicType, "dic-type", "AJLX", "")
	flag.StringVar(&config.dicKeys, "dic-keys", "1,1,1,1", "")
	flag.StringVar(&config.word, "word", "", "")
	flag.StringVar(&config.ay, "ay", "", "")
	flag.StringVar(&config.ajlx, "ajlx", "", "")
	flag.StringVar(&config.slfy, "slfy", "", "")
	flag.StringVar(&config.fycj, "fycj", "", "")
	flag.StringVar(&config.spcx, "spcx", "", "")
	flag.StringVar(&config.wslx, "wslx", "", "")
	flag.StringVar(&config.index, "index", "1", "1-")
	flag.StringVar(&config.page, "page", "10", "must be 10")
	flag.StringVar(&config.order, "order", "法院层级", "裁判日期 / 审判程序")
	flag.StringVar(&config.direction, "direction", "asc", "acs / desc")
	flag.StringVar(&config.method, "method", "list-content", "tree-list,child-all-court,list-content,tree-content,tree-list, dic-value, doc")
	flag.IntVar(&config.fromYear, "from-year", 2019, "")
	flag.IntVar(&config.fromMonth, "from-month", 1, "1-12")
	flag.IntVar(&config.fromDay, "from-day", 1, "1-31")
	flag.IntVar(&config.days, "days", 0, "1+")

	flag.Parse()
	_ = os.MkdirAll(config.repo, 0755)
	if config.days > 0 {
		to := time.Date(config.fromYear, time.Month(config.fromMonth), config.fromDay, 0, 0, 0, 0, nil)
		to = to.AddDate(0, 0, config.days)
		ty, tm, td := to.Date()
		config.cprq = fmt.Sprintf("%04d-%02d-%02d TO %04d-%02d-%02d",
			config.fromYear, config.fromMonth, config.fromDay, ty, tm, td)
	}
}

var config struct {
	docid     string
	dicType   string
	dicKeys   string
	method    string
	repo      string
	word      string
	ay        string
	ajlx      string
	slfy      string
	fycj      string
	spcx      string
	wslx      string
	cprq      string
	index     string
	page      string
	order     string
	direction string
	fromYear  int
	fromMonth int
	fromDay   int
	days      int
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
