// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command hey is an HTTP load generator.
package hey

import (
	"flag"
	"fmt"
	"net/http"
	gourl "net/url"
	"os"
	"regexp"
	"strings"

	"github.com/rakyll/hey/requester"
)


type Hey struct{
	Url string
	Num int  // 请求数
	Con int // 压力数
	Time int

	ContentType string //Content-type
	Method string //请求方法
	Body string //请求内容
	Accept string
	AuthHeader string
	Output string
	ProxyAddr string
	HostHeader string
	Headers string

	DisableCompression bool
	DisableKeepAlives bool
	H2 bool
	EnableTrace bool
}

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
	authRegexp   = `^(.+):([^\s].+)`
)

type headerSlice []string

func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}


func (this *Hey)RunHey() {
	//flag.Usage = func() {
	//	fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	//}

	//flag.Var(&headerslice, "H", "")

	//flag.Parse()
	//if flag.NArg() < 1 {
	//	usageAndExit("")
	//}

	//runtime.GOMAXPROCS(-1)
	num := this.Num
	conc := this.Con
	q := 0

	if num <= 0 || conc <= 0 {
		usageAndExit("-n and -c cannot be smaller than 1.")
	}

	if num < conc {
		usageAndExit("-n cannot be less than -c.")
	}

	if this.Time == 0 {
		this.Time = 20
	}

	url := this.Url
	method := strings.ToUpper(this.Method)

	// set content-type
	header := make(http.Header)
	header.Set("Content-Type", this.ContentType)
	// set any other additional headers
	if this.Headers != "" {
		usageAndExit("Flag '-h' is deprecated, please use '-H' instead.")
	}
	// set any other additional repeatable headers
	//for _, h := range headerslice {
	//	match, err := parseInputWithRegexp(h, headerRegexp)
	//	if err != nil {
	//		usageAndExit(err.Error())
	//	}
	//	header.Set(match[1], match[2])
	//}

	if this.Accept != "" {
		header.Set("Accept", this.Accept)
	}

	// set basic auth if set
	var username, password string
	if this.AuthHeader != "" {
		match, err := parseInputWithRegexp(this.AuthHeader, authRegexp)
		if err != nil {
			usageAndExit(err.Error())
		}
		username, password = match[1], match[2]
	}

	var bodyAll []byte
	if this.Body != "" {
		bodyAll = []byte(this.Body)
	}

	if this.Output != "csv" && this.Output != "" {
		usageAndExit("Invalid output type; only csv is supported.")
	}

	var proxyURL *gourl.URL
	if this.ProxyAddr != "" {
		var err error
		proxyURL, err = gourl.Parse(this.ProxyAddr)
		if err != nil {
			usageAndExit(err.Error())
		}
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		usageAndExit(err.Error())
	}
	req.Header = header
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}

	// set host header if set
	if this.HostHeader != "" {
		req.Host = this.HostHeader
	}

	this.DisableKeepAlives = true
	(&requester.Work{
		Request:            req,
		RequestBody:        bodyAll,
		N:                  num,
		C:                  conc,
		QPS:                q,
		Timeout:            this.Time,
		DisableCompression: this.DisableCompression,
		DisableKeepAlives:  this.DisableKeepAlives,
		H2:                 this.H2,
		ProxyAddr:          proxyURL,
		Output:             this.Output,
		EnableTrace:        this.EnableTrace,
	}).Run()
}

func errAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func parseInputWithRegexp(input, regx string) ([]string, error) {
	re := regexp.MustCompile(regx)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not parse the provided input; input = %v", input)
	}
	return matches, nil
}
