package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// json 数据返回
func ReturnJson(c *gin.Context, status_code int, message string, data interface{}) {
	c.JSON(
		status_code,
		gin.H{
			"status_code": status_code,
			"message":     message,
			"data":        data,
		})
}

// 请求数据接口
func Curl(method, urlVal, data string) (result []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		if len(urlArr) == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}
	//var res map[string]interface{}

	//可以添加多个cookie
	//cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "111", HttpOnly: true}
	//req.AddCookie(cookie1)
	//req.Header.Add("X-Xsrftoken", "1111")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8") //设置Content-Type

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//if err := json.Unmarshal([]byte(string(b)), &res); err != nil {
	//	return res, err
	//}
	return body, nil
}

//将get请求的参数进行转义
func getParseParam(param string) string {
	return url.PathEscape(param)
}
