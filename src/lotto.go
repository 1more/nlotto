package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Input: ID PW")
		os.Exit(1)
	}
	result := checkPensionLotto(args[0], args[1])
	fmt.Println(strings.Index(result, "2016-12-03"))
}

func doLogin(id string, pw string) string {
	//loginUrl := "http://www.nlotto.co.kr/common/cu_login_Result.jsp?event_return_code=&returnUrl=http%3A%2F%2Fnlotto.co.kr%2Fcommon.do%3Fmethod%3Dmain"
	loginUrl := "http://www.nlotto.co.kr/common.do?method=login&returnUrl=http%3A%2F%2Fnlotto.co.kr%2Fcommon.do%3Fmethod%3Dmain"

	client := &http.Client{}
	formData := url.Values{"userId": {id}, "password": {pw}}
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://nlotto.co.kr/common.do?method=main")
	req.Header.Set("Host", "www.nlotto.co.kr")
	req.Header.Set("DNT", "1")
	resp, err := client.Do(req)

	//resp, err := http.PostForm(loginUrl, url.Values{"userId": {id}, "password": {pw}})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}
	respHeader := resp.Header
	cookie1 := respHeader.Get("Set-Cookie")
	cookie2 := respHeader.Get("Set-Cookie2")

	return cookie2 + ";" + cookie1
}

func checkPensionLotto(id string, pw string) string {
	rawCookies := doLogin(id, pw)
	fmt.Println(rawCookies)

	client := &http.Client{}
	header := http.Header{}
	header.Add("Cookies", rawCookies)
	url := "http://www.nlotto.co.kr/lotto520.do?method=pensionBuyReserveList"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}

	req.Header.Set("Cookie", rawCookies)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return parseLottoResult(string(body))
}

func parseLottoResult(result string) string {
	return result
}
