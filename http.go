package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"strings"
)

// 包装一个普通通用http请求。 post/put/patch使用表单提交
func DoHTTPRequest(method, url string, headers, params *map[string]string) (code int, response string, err error) {
	var (
		req *http.Request

		params_x = _url.Values{}
		url_x    string
	)

	if params != nil {
		for key, value := range *params {
			params_x.Set(key, value)
		}
	}

	switch method {
	case "GET", "DELETE":
		// 如果是get和delete请求，则也把params放到url上来
		if strings.Index(url, "?") > 0 {
			url_x = fmt.Sprintf("%s&%s", url, params_x.Encode())
		} else {
			url_x = fmt.Sprintf("%s?%s", url, params_x.Encode())
		}

		req, err = http.NewRequest(method, url_x, nil)
		if err != nil {
			return
		}
	case "POST", "PUT", "PATCH":
		req, err = http.NewRequest(method, url, strings.NewReader(params_x.Encode()))
		if err != nil {
			return
		}
	default:
		req, err = http.NewRequest(method, url, strings.NewReader(params_x.Encode()))
		if err != nil {
			return
		}
	}

	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
			// req.Header.Set("Authorization", "Bearer YMCss6Ei-fvq6TnMG4hF2V5uC6c57rQPM2AZBzCvxdA30gKK")
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

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

	code = resp.StatusCode
	response = fmt.Sprintf("%s", bytes)

	return
}

// 包装一个适合传递json的通用http请求。body参数可以为结构体/map/string/nil。 post/put/patch使用json提交，params参数被添加到url上（如zoom api）
func DoHTTPJsonRequest(method, url string, body interface{}, headers, params *map[string]string) (code int, response string, err error) {
	var (
		req *http.Request

		params_x = _url.Values{}
		body_x   io.Reader
		url_x    string
	)

	if params != nil {
		for key, value := range *params {
			params_x.Set(key, value)
		}
	}

	switch body.(type) {
	case nil:
		body_x = nil
	case string:
		body_x = strings.NewReader(body.(string))
	case []byte:
		body_x = bytes.NewReader(body.([]byte))
	default:
		// TODO
		body_bytes, err := json.Marshal(body)
		if err != nil {
			return code, response, err
		}

		body_x = bytes.NewReader(body_bytes)
	}

	if strings.Index(url, "?") > 0 {
		url_x = fmt.Sprintf("%s&%s", url, params_x.Encode())
	} else {
		url_x = fmt.Sprintf("%s?%s", url, params_x.Encode())
	}

	req, err = http.NewRequest(method, url_x, body_x)

	if err != nil {
		return
	}

	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
			// req.Header.Set("Authorization", "Bearer YMCss6Ei-fvq6TnMG4hF2V5uC6c57rQPM2AZBzCvxdA30gKK")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

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

	code = resp.StatusCode
	response = fmt.Sprintf("%s", bytes)

	return
}
