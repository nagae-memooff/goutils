package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"strings"
)

func DoHTTPRequest(method, url string, headers, params *map[string]string) (code int, response string, err error) {
	params_x := _url.Values{}

	if params != nil {
		for key, value := range *params {
			params_x.Set(key, value)
		}
	}

	req, err := http.NewRequest(method, url, strings.NewReader(params_x.Encode()))
	if err != nil {
		return
	}

	if headers != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for key, value := range *headers {
			req.Header.Set(key, value)
			// req.Header.Set("Authorization", "Bearer YMCss6Ei-fvq6TnMG4hF2V5uC6c57rQPM2AZBzCvxdA30gKK")
		}
	}

	//   if method == "POST" {
	// req.PostForm.Set("login_name", "mudy@22.com")
	// req.PostForm.Set("password", "11111111")
	//   }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 打印下状态码，看下效果
	// fmt.Printf("返回的状态码是： %v\n", resp.StatusCode)
	// fmt.Printf("返回的信息是： %v\n", resp.StatusCode)

	// fmt.Printf("%s", bytes)
	code = resp.StatusCode
	response = fmt.Sprintf("%s", bytes)

	return
}
