package packsender

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type JobResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ZpData  struct {
		HasMore  bool `json:"hasMore"`
		ResCount int  `json:"resCount"`
		JobList  []struct {
			JobName       string   `json:"jobName"`
			SalaryDesc    string   `json:"salaryDesc"`
			BrandName     string   `json:"brandName"`
			JobExperience string   `json:"jobExperience"`
			JobDegree     string   `json:"jobDegree"`
			Skills        []string `json:"skills"`
			EncryptJobId  string   `json:"encryptJobId"`
			CityName      string   `json:"cityName"`
		} `json:"jobList"`
	} `json:"zpData"`
}

func loadCookies(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var cookieMap map[string]string
	if err := json.NewDecoder(file).Decode(&cookieMap); err != nil {
		return "", err
	}

	var cookieStr strings.Builder
	for k, v := range cookieMap {
		cookieStr.WriteString(fmt.Sprintf("%s=%s; ", k, v))
	}
	return strings.TrimSpace(cookieStr.String()), nil
}

func FetchJobs(lastCity string) error {
	cookieHeader, err := loadCookies("cookies.json")
	fmt.Println("使用的 Cookie:", cookieHeader)
	if err != nil {
		return fmt.Errorf("读取 cookie 失败: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		params := url.Values{}
		params.Add("page", "1")
		params.Add("pageSize", "15")
		params.Add("city", lastCity)
		params.Add("query", "网络安全")
		params.Add("scene", "1")
		params.Add("_", fmt.Sprintf("%d", time.Now().UnixMilli()))

		urlStr := "https://www.zhipin.com/wapi/zpgeek/search/joblist.json?" + params.Encode()

		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}

		req.Header.Set("Cookie", cookieHeader)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("请求失败:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		defer resp.Body.Close()

		var body []byte
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(resp.Body)
			body, _ = io.ReadAll(reader)
		case "deflate":
			reader := flate.NewReader(resp.Body)
			body, _ = io.ReadAll(reader)
		default:
			body, _ = io.ReadAll(resp.Body)
		}

		utf8Body, _, _ := charset.DetermineEncoding(body, "")
		reader := transform.NewReader(bytes.NewReader(body), utf8Body.NewDecoder())
		decodedBody, _ := io.ReadAll(reader)

		fmt.Println("响应数据:", string(decodedBody))

		var jobRes JobResponse
		if err := json.Unmarshal(decodedBody, &jobRes); err != nil {
			fmt.Println("解析 JSON 失败:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if jobRes.Code == 37 && strings.Contains(jobRes.Message, "访问行为异常") {
			fmt.Println("访问异常，3秒后重试...")
			time.Sleep(3 * time.Second)
			continue
		}

		for _, job := range jobRes.ZpData.JobList {
			fmt.Printf("%s [%s] - %s\n", job.JobName, job.SalaryDesc, job.BrandName)
		}

		break // 成功获取后跳出循环
	}

	return nil
}
