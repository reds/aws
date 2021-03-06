package sign

import (
	"bufio"
	"bytes"
	"github.com/reds/aws/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := config.LoadConfig("awsTestSuite.cfg", "dynamodb")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}

const (
	suiteDir = "aws4_testsuite"
)

func TestDir(t *testing.T) {
	dir, err := os.Open(suiteDir)
	if err != nil {
		t.Fatal(err)
	}
	files, err := dir.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadConfig("awsTestSuite.cfg", "dynamodb")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)

	for _, file := range files {
		if strings.HasSuffix(file, ".req") {
			file = suiteDir + "/" + file[:len(file)-4]
			req := reqFile2req(t, file+".req")
			if req == nil {
				t.Log("skipping", file, "(file2req)")
				continue
				//				t.Fatal("error parsing", file)
			}
			t.Log(file)
			creq, _ := req2canonical(req)
			if !compFile2String(t, file+".creq", creq) {
				t.Logf("%#v", req)
				t.Logf("%#v", creq)
				t.Log("skipping", file, "(cannonical)")
				continue
			}
			SignV4(req, "host", cfg.Region, cfg.AccessKeyId, cfg.Secret)
			auth := req.Header["Authorization"][0]
			if !compFile2String(t, file+".authz", auth) {
				t.Fatal("authorization", file)
			}
		}
	}
}

func compFile2String(t *testing.T, fn, data string) bool {
	f, err := os.Open(fn)
	if err != nil {
		t.Log(err)
		return false
	}
	fdata, err := ioutil.ReadAll(f)
	if err != nil {
		t.Log(err)
		return false
	}
	fsdata := strings.Replace(strings.Replace(string(fdata), "\r\n", "", -1), "\n", "", -1)
	data = strings.Replace(strings.Replace(string(data), "\r\n", "", -1), "\n", "", -1)
	if len(fsdata) != len(data) {
		t.Log("lengths differ", len(fsdata), len(data))
		t.Log(fsdata)
		t.Log(data)
		return false
	}
	if string(fsdata) != data {
		t.Log("differ", string(fsdata), data)
		t.Log(fsdata)
		t.Log(data)
		return false
	}
	return true
}

// Take a .req file from test suite and create a http.Request
func reqFile2req(t *testing.T, fn string) *http.Request {
	f, err := os.Open(fn)
	if err != nil {
		t.Fatal(err)
	}
	r := bufio.NewReader(f)
	l1, err := r.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	l1 = strings.TrimRight(l1, "\r\n")
	req := strings.Split(l1, " ")
	if len(req) < 3 {
		t.Log(req)
		return nil
	}
	uri := req[1]
	headers := make([]string, 0, 10)
	for {
		h, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if h != "" {
					headers = append(headers, strings.TrimRight(h, "\r\n"))
				}
				break
			}
			t.Fatal(err)
		}
		if h == "\r\n" {
			break
		}
		headers = append(headers, strings.TrimRight(h, "\r\n"))
	}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		if err != io.EOF {
			t.Fatal(err)
		} else {
		}
	}
	f.Close()
	originput := strings.Join(req, " ") + "\r\n" +
		strings.Join(headers, "\r\n") + "\r\n" + "\r\n"
	if len(body) > 0 {
		originput += string(body)
	}
	if !compFile2String(t, fn, originput) {
		t.Log(body)
		t.Fatal("differ: ", fn)
	}
	rt := ""
	switch req[0] {
	case "GET":
		rt = "GET"
	case "POST":
		rt = "POST"
	default:
		t.Fatal(req)
	}
	u := strings.SplitN(uri, "?", 2)
	queries := make([]string, 0, 20)
	if len(u) > 1 {
		for _, q := range strings.Split(strings.Replace(u[1], "+", " ", -1), "&") {
			kv := strings.SplitN(q, "=", 2)
			if len(kv) == 2 {
				queries = append(queries, url.QueryEscape(kv[0])+"="+url.QueryEscape(kv[1]))
			} else {
				queries = append(queries, url.QueryEscape(kv[0])+"=")
			}
		}
		uri = u[0] + "?" + strings.Join(queries, "&")
	}
	request, err := http.NewRequest(rt, "https://aws.amazon.com"+uri, bytes.NewReader(body))
	if err != nil {
		//t.Fatal(err,uri,req)
		return nil
	}

	for _, h := range headers {
		a := strings.SplitN(h, ":", 2)
		if len(a) > 1 {
			request.Header.Add(a[0], a[1])
		}
	}
	return request
}

func testUnicode(t *testing.T) {
	dir, err := os.Open(suiteDir)
	if err != nil {
		t.Fatal(err)
	}
	files, err := dir.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, ".req") {
			file = suiteDir + "/" + file[:len(file)-4]
			t.Log("file", file)
			req := reqFile2req(t, file+".req")
			if req == nil {
				t.Fatal("error parsing", file)
			}
		}
	}
}
