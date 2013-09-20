package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func SignV4(req *http.Request, service, region, accesskey, secret string) {
	creq, signedHeaders := req2canonical(req)
	sha := sha256.New()
	sha.Write([]byte(creq))
	chash := sha.Sum(nil)

	t := time.Now().UTC()
	// for the test suite date has already been set
	if len(req.Header["Date"]) > 0 {
		t, _ = time.Parse(time.RFC1123, req.Header["Date"][0])
	}
	req.Header.Set("x-amz-date", t.Format("20060102T150405Z"))
	string2sign := "AWS4-HMAC-SHA256\n" + t.Format("20060102T150405Z") + "\n" +
		t.Format("20060102") + "/" + region + "/" + service + "/aws4_request\n" +
		fmt.Sprintf("%x", chash)
	sig := doHmac(doHmac(doHmac(doHmac(("AWS4"+secret), (t.Format("20060102"))), (region)), (service)), ("aws4_request"))

	sig = fmt.Sprintf("%x", doHmac(sig, string2sign))
	sig = "AWS4-HMAC-SHA256 Credential=" + accesskey + "/" + t.Format("20060102") + "/" + region + "/" + service + "/aws4_request, SignedHeaders=" + signedHeaders + ", Signature=" + sig
	req.Header.Set("Authorization", sig)
}

func doHmac(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(data))
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(h.Sum(nil))
}

// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func req2canonical(req *http.Request) (string, string) {
	/*
	     CanonicalRequest =
	     HTTPRequestMethod + '\n' +
	     CanonicalURI + '\n' +
	     CanonicalQueryString + '\n' +
	     CanonicalHeaders + '\n' +
	     SignedHeaders + '\n' +
	   		HexEncode(Hash(Payload))
	*/

	//  HTTPRequestMethod + '\n' +
	canon := req.Method + "\n"

	//  CanonicalURI + '\n' +
	// take care of unicode in path
	p := unUnicode(path.Clean(req.URL.Path))
	// undo some of path.Clean's work
	if len(p) > 1 && strings.HasSuffix(req.URL.Path, "/") {
		p += "/"
	}
	if p == "" {
		p = "/"
	}
	// space is special
	p = strings.Replace(p, " ", "%20", -1)
	canon += p + "\n"

	//  CanonicalQueryString + '\n' +
	queries := make([]string, 0, 20)
	if len(req.URL.RawQuery) > 0 {
		for _, q := range strings.Split(req.URL.RawQuery, "&") {
			kv := strings.SplitN(strings.Replace(q, "+", "%20", -1), "=", 2)
			if len(kv) == 2 {
				queries = append(queries, unUnicode(kv[0])+"="+unUnicode(kv[1]))
			} else {
				queries = append(queries, unUnicode(kv[0])+"=")
			}
		}
	}
	sort.StringSlice(queries).Sort()
	canon += strings.Join(queries, "&") + "\n"

	//  CanonicalHeaders + '\n' +
	headers := make([]string, 0, len(req.Header))
	headerList := make([]string, 0, len(req.Header))
	for k, v := range req.Header {
		if strings.Index(strings.ToLower(k), "x-") == 0 {
			if strings.ToLower(k) != "x-amz-date" {
				continue
			}
		}
		b := make([]string, 0, 10)
		for _, s := range v {
			b = append(b, strings.TrimSpace(s))
		}
		sort.StringSlice(b).Sort()
		headers = append(headers, strings.ToLower(k)+":"+strings.Join(b, ","))
		headerList = append(headerList, strings.ToLower(k))
	}
	sort.StringSlice(headers).Sort()
	canon += strings.Join(headers, "\n") + "\n"
	canon += "\n"

	//  SignedHeaders + '\n' +
	sort.StringSlice(headerList).Sort()
	canon += strings.Join(headerList, ";") + "\n"

	//  HexEncode(Hash(Payload))
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(b)) // reset body to orig after read
	sha := sha256.New()
	sha.Write(b)
	hash := sha.Sum(nil)
	canon += fmt.Sprintf("%x", hash)
	return canon, strings.Join(headerList, ";")
}

// extended UTF-8 characters must be in the form %XY%ZA%BC
func unUnicode(in string) string {
	b := []byte(in)
	out := ""
	for len(b) > 0 {
		_, size := utf8.DecodeRune(b)
		if size > 1 {
			for i := 0; i < size; i++ {
				out += fmt.Sprintf("%%%X", b[i])
			}
		} else {
			out += string(b[0])
		}
		b = b[size:]
	}
	return out
}
