package utils

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
)

type Request struct {
	Method     string
	Headers    map[string][]string
	Url        string
	Body       string
	BodyReader io.Reader
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
}

func CallHTTPService(request Request) Response {
	client := http.Client{}
	defer client.CloseIdleConnections()
	//ask to X-CSRF-Token and cookies
	if request.Method != "GET" {
		req, _ := http.NewRequest("HEAD", request.Url, nil)
		request.Headers["X-CSRF-Token"] = []string{"fetch"}
		req.Header = request.Headers
		res, err := client.Do(req)
		if err != nil {
			return Response{
				StatusCode: 500,
				Body:       err.Error(),
			}
		}
		request.Headers["X-CSRF-Token"] = []string{res.Header.Get("X-Csrf-Token")}
	}
	/*
		if strings.Contains(request.Url, "/simulations") {
			request.Url = "https://WebServiceTest.cfapps.us10-001.hana.ondemand.com"
		}
	*/
	var body io.Reader

	if request.BodyReader == nil {
		body = bytes.NewBuffer([]byte(request.Body))
	} else {
		body = request.BodyReader
	}

	req, _ := http.NewRequest(request.Method, request.Url, body)
	req.Header = request.Headers
	res, err := client.Do(req)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       err.Error(),
		}
	} else {
		bodyString, _ := io.ReadAll(res.Body)
		return Response{
			StatusCode: res.StatusCode,
			Headers:    res.Header,
			Body:       string(bodyString),
		}
	}

}

func DetectNEOFoundry(tenant string) string {
	//Detect NEO or Foundry
	var re = regexp.MustCompile(`(?mi)^.*integrationsuite.*.cfapps`)
	if !re.MatchString(tenant) {
		tenant += "/itspaces"
	}
	return tenant
}
