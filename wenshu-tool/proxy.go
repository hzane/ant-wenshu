package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/http/httpguts"
)

const magic = "4bd54ad7a845"

func main() {
	http.HandleFunc("/"+magic+"/", expose) // 获取cookies
	http.HandleFunc("/", wenshu)           // 反向代理
	_ = http.ListenAndServe(":18071", nil)
}
func expose(w http.ResponseWriter, r *http.Request) {
	cookies, _ := config.cookies.Load().([]*http.Cookie)
	w.Header().Set("content-type", "text/plain")
	for _, ck := range cookies {
		_, _ = fmt.Fprintln(w, ck.Name, "=", ck.Value)
	}
}

type BackupBody struct {
	data     []byte
	internal io.ReadCloser
}

func (bb *BackupBody) Read(p []byte) (n int, err error) {
	n, err = bb.internal.Read(p)
	if n > 0 {
		bb.data = append(bb.data, p[:n]...)
	}
	return
}
func (bb *BackupBody) Close() error {
	return bb.internal.Close()
}

func wenshu(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = config.base.Host // 将请求地址改成真实的文书地址
	r.URL.Scheme = config.base.Scheme
	r.Host = config.base.Host

	ext := path.Ext(strings.ToLower(r.URL.Path)) // 忽略所有的资源文件
	ignore := ext != "" && strings.Contains(".css.jpg.png.gif.js.jpeg.jpe.ico", ext)
	ignore = ignore || strings.HasPrefix(r.URL.Path, "/Index/") // 这个地址有个定时器，我们也不关心这个

	uwc := r.Header.Get("use-wenshu-cookie") == magic
	if uwc { // 来自我们的程序的请求，完成cookie替换
		r.Header.Del("cookie")
		r.Header.Del("use-wenshu-cookie")
		for _, cookie := range config.cookies.Load().([]*http.Cookie) {
			r.AddCookie(cookie)
			log.Println("cookie-force", cookie.Name, cookie.Value)
		}
		// UA和cookie必须保持一致
		r.Header.Set("user-agent", config.ua.Load().(string))
	}
	// 通过trailer返回给 浏览器所有cookie
	trailer := r.Header.Get("use-wenshu-trailer") == magic

	proxy := httputil.NewSingleHostReverseProxy(config.base)
	proxy.ServeHTTP(w, r)
	if ignore {
		return
	}
	log.Println(r.Method, r.URL)

	// 记录请求中的cookie和服务器返回的set-cookie
	// cookie名称有可能重复
	cookies := append(r.Cookies(), readSetCookies(w.Header())...)
	if !uwc {
		config.cookies.Store(cookies)
		config.ua.Store(r.Header.Get("user-agent"))
	}

	if trailer {
		for _, cookie := range cookies {
			w.Header().Add("wenshu-cookie", cookie.String())
		}
	}
}

func init() {
	flag.Parse()
	config.base, _ = url.Parse("http://wenshu.court.gov.cn")
	config.ua.Store("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3766.2 Safari/537.36")
}

var config struct {
	base    *url.URL
	cookies atomic.Value
	ua      atomic.Value
}

// copied from http.go, 后面的代码不重要
func readSetCookies(h http.Header) []*http.Cookie {
	cookieCount := len(h["Set-Cookie"])

	if cookieCount == 0 {
		return []*http.Cookie{}
	}

	cookies := make([]*http.Cookie, 0, cookieCount)

	for _, line := range h["Set-Cookie"] {
		parts := strings.Split(strings.TrimSpace(line), ";")

		if len(parts) == 1 && parts[0] == "" {
			continue
		}
		parts[0] = strings.TrimSpace(parts[0])
		j := strings.Index(parts[0], "=")

		if j < 0 {
			continue
		}

		name, value := parts[0][:j], parts[0][j+1:]
		if !isCookieNameValid(name) {
			continue
		}

		value, ok := parseCookieValue(value, true)
		if !ok {
			continue
		}

		c := &http.Cookie{
			Name:  name,
			Value: value,
			Raw:   line,
		}

		for i := 1; i < len(parts); i++ {
			parts[i] = strings.TrimSpace(parts[i])
			if len(parts[i]) == 0 {
				continue
			}

			attr, val := parts[i], ""
			if j := strings.Index(attr, "="); j >= 0 {
				attr, val = attr[:j], attr[j+1:]
			}
			lowerAttr := strings.ToLower(attr)
			val, ok = parseCookieValue(val, false)
			if !ok {
				c.Unparsed = append(c.Unparsed, parts[i])
				continue
			}

			switch lowerAttr {
			case "samesite":
				lowerVal := strings.ToLower(val)
				switch lowerVal {
				case "lax":
					c.SameSite = http.SameSiteLaxMode
				case "strict":
					c.SameSite = http.SameSiteStrictMode
				default:
					c.SameSite = http.SameSiteDefaultMode
				}
				continue
			case "secure":
				c.Secure = true
				continue
			case "httponly":
				c.HttpOnly = true
				continue
			case "domain":
				c.Domain = val
				continue
			case "max-age":
				secs, err := strconv.Atoi(val)
				if err != nil || secs != 0 && val[0] == '0' {
					break
				}

				if secs <= 0 {
					secs = -1
				}
				c.MaxAge = secs
				continue

			case "expires":
				c.RawExpires = val
				exptime, err := time.Parse(time.RFC1123, val)

				if err != nil {
					exptime, err = time.Parse("Mon, 02-Jan-2006 15:04:05 MST", val)
					if err != nil {
						c.Expires = time.Time{}
						break
					}
				}
				c.Expires = exptime.UTC()
				continue
			case "path":
				c.Path = val
				continue
			}

			c.Unparsed = append(c.Unparsed, parts[i])
		}

		cookies = append(cookies, c)
	}

	return cookies
}
func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, httpguts.IsTokenRune) < 0
}
func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}
