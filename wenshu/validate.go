package main

import (
	"bytes"
	"errors"
	"github.com/otiai10/gosseract"
	"gitlab.com/hearts.zhang/tools"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// ValidateCode ...
func ValidateCode(client *tools.Client) (err error) {
	// req, _ := http.NewRequest("POST", ValidateCodeURL, )
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
	text = strings.Map(func(r rune) rune {
		if r == 'B' {
			return '8'
		}
		if r == 'O' {
			return '0'
		}
		return r
	}, text)
	code := url.Values{
		"ValidateCode": []string{text},
	}
	// Html_Pages/VisitRemind.html
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

