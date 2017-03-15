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
	"runtime"
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

var (
	headerslice headerSlice
	m           = flag.String("m", "GET", "")
	headers     = flag.String("h", "", "")
	body        = flag.String("d", "", "")
	bodyFile    = flag.String("D", "", "")
	accept      = flag.String("A", "", "")
	contentType = flag.String("T", "text/html", "")
	authHeader  = flag.String("a", "", "")
	hostHeader  = flag.String("host", "", "")

	output = flag.String("o", "", "")

	c = flag.Int("c", 50, "")
	n = flag.Int("n", 200, "")
	q = flag.Int("q", 0, "")
	t = flag.Int("t", 20, "")

	h2 = flag.Bool("h2", false, "")

	cpus = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")

	disableCompression = flag.Bool("disable-compression", false, "")
	disableKeepAlives  = flag.Bool("disable-keepalive", true, "")
	proxyAddr          = flag.String("x", "", "")

	enableTrace = flag.Bool("more", false, "")
)

var usage = `Usage: hey [options...] <url>

Options:
  -n  Number of requests to run. Default is 200.
  -c  Number of requests to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -q  Rate limit, in seconds (QPS).
  -o  Output type. If none provided, a summary is printed.
      "csv" is the only supported alternative. Dumps the response
      metrics in comma-separated values format.

  -m  HTTP method, one of GET, POST, PUT, DELETE, HEAD, OPTIONS.
  -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
      For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
  -t  Timeout for each request in seconds. Default is 20, use 0 for infinite.
  -A  HTTP Accept header.
  -d  HTTP request body.
  -D  HTTP request body from file. For example, /home/user/file.txt or ./file.txt.
  -T  Content-type, defaults to "text/html".
  -a  Basic authentication, username:password.
  -x  HTTP Proxy address as host:port.
  -h2 Enable HTTP/2.

  -host	HTTP Host header.

  -disable-compression  Disable compression.
  -disable-keepalive    Disable keep-alive, prevents re-use of TCP
                        connections between different HTTP requests.
  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
  -more                 Provides information on DNS lookup, dialup, request and
                        response timings.
`



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
	for _, h := range headerslice {
		match, err := parseInputWithRegexp(h, headerRegexp)
		if err != nil {
			usageAndExit(err.Error())
		}
		header.Set(match[1], match[2])
	}

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
	//if *bodyFile != "" {
	//	slurp, err := ioutil.ReadFile(*bodyFile)
	//	if err != nil {
	//		errAndExit(err.Error())
	//	}
	//	bodyAll = slurp
	//}

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
