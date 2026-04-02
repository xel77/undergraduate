package login

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

// Login 结构体，存储登录相关数据
type Login struct {
	ZpData       ZpData       // 二维码登录相关字段
	Client       *http.Client // 带 CookieJar 的 HTTP 客户端
	Pk           string       // 页面标识，如 "header-login"
	Fp           string       // 前端指纹
	DispatchData DispatchResp // dispatcher 接口返回
}

// ZpData 存储二维码登录的返回数据
type ZpData struct {
	QrId      string `json:"qrId"`
	RandKey   string `json:"randKey"`
	SecretKey string `json:"secretKey"`
	ShortKey  string `json:"shortRandKey"`
}

// RandKeyResp 解析获取 RandKey 接口返回结构
type RandKeyResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ZpData  ZpData `json:"zpData"`
}

// ScanStatusResp 解析扫码状态接口返回结构
type ScanStatusResp struct {
	Scaned    bool `json:"scaned"`
	AllWeb    bool `json:"allweb"`
	NewScaned bool `json:"newScaned"`
	ScanedV2  bool `json:"scanedV2"`
}

// ScanLoginResp 解析 scanLogin 接口返回结构
type ScanLoginResp struct {
	Msg    string `json:"msg"`
	Scaned bool   `json:"scaned"`
}

// SecondKeyResp 解析 getSecondKey 接口返回结构
type SecondKeyResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ZpData  struct {
		QrId string `json:"qrId"`
	} `json:"zpData"`
}

// DispatchResp 解析 dispatcher 接口返回结构
type DispatchResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ZpData  struct {
		ToUrl    string  `json:"toUrl"`
		Identity int     `json:"identity"`
		PcToUrl  string  `json:"pcToUrl"`
		Version  float64 `json:"version"`
	} `json:"zpData"`
}

// CookieData 存储登录 Cookie 信息
type CookieData struct {
	Wt2  string `json:"wt2"`
	Wbg  string `json:"wbg"`
	ZpAt string `json:"zp_at"`
	Bst  string `json:"bst"`
}

// 公共请求头
var commonHeaders = map[string]string{
	"Host":               "www.zhipin.com",
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
	"Referer":            "https://www.zhipin.com/web/user/?ka=header-login",
	"Accept-Encoding":    "gzip, deflate, br",
	"Sec-Ch-Ua":          `"Chromium";v="135", "Not-A.Brand";v="8"`,
	"Sec-Ch-Ua-Mobile":   "?0",
	"Sec-Ch-Ua-Platform": `"Windows"`,
	"Accept-Language":    "zh-CN,zh;q=0.9",
}

// newRequest 创建带公共头的 HTTP 请求
func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range commonHeaders {
		req.Header.Set(k, v)
	}
	return req, nil
}

// processResponse 处理 HTTP 响应，自动解压 gzip
func processResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip 解压失败: %v", err)
		}
		defer gz.Close()
		reader = gz
	}
	return io.ReadAll(reader)
}

// Init 初始化 HTTP 客户端及 pk/fp
func (l *Login) Init(pk, fp string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	l.Client = &http.Client{Jar: jar}
	l.Pk = pk
	l.Fp = fp
	return nil
}

// GetLoginQr 获取 randkey 并下载二维码
func (l *Login) GetLoginQr() error {
	// 1. 请求 randkey
	req, err := newRequest("POST", "https://www.zhipin.com/wapi/zppassport/captcha/randkey", nil)
	if err != nil {
		return fmt.Errorf("创建 randkey 请求失败: %v", err)
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Traceid", fmt.Sprintf("F-%d", time.Now().UnixNano()))

	resp, err := l.Client.Do(req)
	if err != nil {
		return fmt.Errorf("请求 randkey 失败: %v", err)
	}
	data, err := processResponse(resp)
	if err != nil {
		return fmt.Errorf("读取 randkey 响应失败: %v", err)
	}

	var rr RandKeyResp
	if err := json.Unmarshal(data, &rr); err != nil {
		return fmt.Errorf("解析 randkey JSON 失败: %v", err)
	}
	if rr.Code != 0 {
		return fmt.Errorf("randkey 返回 code=%d, message=%s", rr.Code, rr.Message)
	}
	l.ZpData = rr.ZpData
	fmt.Printf("[RandKey] qrId=%s\n", l.ZpData.QrId)

	// 2. 下载二维码图片
	qrURL := fmt.Sprintf("https://www.zhipin.com/wapi/zpweixin/qrcode/getqrcode?content=%s&w=200&h=200", l.ZpData.QrId)
	req2, err := newRequest("GET", qrURL, nil)
	if err != nil {
		return fmt.Errorf("创建二维码请求失败: %v", err)
	}
	req2.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	resp2, err := l.Client.Do(req2)
	if err != nil {
		return fmt.Errorf("请求二维码失败: %v", err)
	}
	img, err := processResponse(resp2)
	if err != nil {
		return err
	}
	if err := os.WriteFile("qrcode.png", img, 0644); err != nil {
		return err
	}
	fmt.Println("二维码已保存为 qrcode.png")
	return nil
}

// PollQrLogin 轮询扫码并自动完成登录流程
func (l *Login) PollQrLogin(timeout time.Duration) error {
	scanTicker := time.NewTicker(2 * time.Second)       // 扫码状态检查间隔
	qrRefreshTicker := time.NewTicker(30 * time.Second) // 二维码刷新间隔
	defer scanTicker.Stop()
	defer qrRefreshTicker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-scanTicker.C:
			// 1. 检查扫码状态
			status, err := l.getScanStatus()
			if err != nil {
				return err
			}
			if !status.Scaned {
				fmt.Println("[Scan] 等待扫码...")
				continue
			}
			fmt.Println("[Scan] 二维码已扫码，检查登录确认...")

			// 2. 检查登录确认
			scanLogin, err := l.getScanLogin()
			if err != nil {
				return err
			}
			if !scanLogin {
				fmt.Println("[ScanLogin] 等待用户确认登录...")
				continue
			}
			fmt.Println("[ScanLogin] 用户已确认登录，获取新 qrId...")

			// 3. 获取新 qrId
			code, newQrId, err := l.getSecondKey()
			if err != nil {
				return err
			}
			if code != 0 {
				return fmt.Errorf("getSecondKey 返回 code=%d", code)
			}
			//l.ZpData.QrId = newQrId
			fmt.Printf("[SecondKey] 新 qrId=%s\n", newQrId)

			// 4. 调用 dispatcher 完成登录
			dr, err := sendDispatcher(l.Client, l.ZpData.QrId, l.Pk, l.Fp)
			if err != nil {
				return err
			}
			l.DispatchData = dr
			fmt.Printf("[Dispatcher] 登录成功，跳转至: %s\n", dr.ZpData.ToUrl)
			return nil

		case <-qrRefreshTicker.C:
			// 刷新二维码
			fmt.Println("[QR] 刷新二维码...")
			if err := l.GetLoginQr(); err != nil {
				fmt.Printf("[QR] 刷新失败: %v\n", err)
				continue
			}
			fmt.Printf("[QR] 二维码已更新，新 qrId=%s\n", l.ZpData.QrId)

		case <-timer.C:
			return fmt.Errorf("自动登录超时: %v", timeout)
		}
	}
}

// getScanStatus 请求 scan 接口
func (l *Login) getScanStatus() (*ScanStatusResp, error) {
	url := fmt.Sprintf("https://www.zhipin.com/wapi/zppassport/qrcode/scan?uuid=%s", l.ZpData.QrId)
	req, err := newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 scan 请求失败: %v", err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Traceid", fmt.Sprintf("F-%d", time.Now().UnixNano()))
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("scan 请求失败: %v", err)
	}
	data, err := processResponse(resp)
	if err != nil {
		return nil, err
	}
	var status ScanStatusResp
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("解析 scan JSON 失败: %v", err)
	}
	return &status, nil
}

// getScanLogin 检查登录确认状态
func (l *Login) getScanLogin() (bool, error) {
	url := fmt.Sprintf("https://www.zhipin.com/wapi/zppassport/qrcode/scanLogin?qrId=%s", l.ZpData.QrId)
	req, err := newRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("创建 scanLogin 请求失败: %v", err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Traceid", fmt.Sprintf("F-%d", time.Now().UnixNano()))
	resp, err := l.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("scanLogin 请求失败: %v", err)
	}
	data, err := processResponse(resp)
	if err != nil {
		return false, err
	}
	var r ScanLoginResp
	if err := json.Unmarshal(data, &r); err != nil {
		return false, fmt.Errorf("解析 scanLogin JSON 失败: %v", err)
	}
	return r.Scaned, nil
}

// getSecondKey 获取新 qrId
func (l *Login) getSecondKey() (int, string, error) {
	url := fmt.Sprintf("https://www.zhipin.com/wapi/zppassport/captcha/getSecondKey?uuid=%s", l.ZpData.QrId)
	req, err := newRequest("GET", url, nil)
	if err != nil {
		return -1, "", fmt.Errorf("创建 getSecondKey 请求失败: %v", err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Traceid", fmt.Sprintf("F-%d", time.Now().UnixNano()))
	resp, err := l.Client.Do(req)
	if err != nil {
		return -1, "", fmt.Errorf("getSecondKey 请求失败: %v", err)
	}
	data, err := processResponse(resp)
	if err != nil {
		return -1, "", err
	}
	var r SecondKeyResp
	if err := json.Unmarshal(data, &r); err != nil {
		return -1, "", fmt.Errorf("解析 getSecondKey JSON 失败: %v", err)
	}
	return r.Code, r.ZpData.QrId, nil
}

// sendDispatcher 调用 dispatcher 完成登录
func sendDispatcher(client *http.Client, qrId, pk, fp string) (DispatchResp, error) {
	// 构造请求URL
	url := fmt.Sprintf("https://www.zhipin.com/wapi/zppassport/qrcode/dispatcher?qrId=%s&pk=%s&fp=%s",
		qrId, pk, fp)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DispatchResp{}, fmt.Errorf("创建 dispatcher 请求失败: %v", err)
	}

	// 设置请求头
	req.Header = http.Header{
		"Host":               {"www.zhipin.com"},
		"Accept":             {"application/json, text/plain, */*"},
		"Accept-Language":    {"zh-CN,zh;q=0.9"},
		"Accept-Encoding":    {"gzip, deflate, br"},
		"Connection":         {"keep-alive"},
		"Referer":            {"https://www.zhipin.com/web/user/?ka=header-login"},
		"Sec-Ch-Ua":          {`"Chromium";v="135", "Not-A.Brand";v="8"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-origin"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"},
		"Traceid":            {fmt.Sprintf("F-%d", time.Now().UnixNano())},
		"Priority":           {"u=1, i"},
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return DispatchResp{}, fmt.Errorf("dispatcher 请求失败: %v", err)
	}

	// 保存关键 Cookie
	cookieData := CookieData{}
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "wt2":
			cookieData.Wt2 = cookie.Value
		case "wbg":
			cookieData.Wbg = cookie.Value
		case "zp_at":
			cookieData.ZpAt = cookie.Value
		case "bst":
			cookieData.Bst = cookie.Value
		}
	}

	// 将 Cookie 保存到文件
	cookieJSON, err := json.MarshalIndent(cookieData, "", "  ")
	if err != nil {
		fmt.Printf("警告: Cookie 序列化失败: %v\n", err)
	} else {
		if err := os.WriteFile("cookies.json", cookieJSON, 0644); err != nil {
			fmt.Printf("警告: Cookie 保存失败: %v\n", err)
		} else {
			fmt.Println("Cookie 已保存到 cookies.json")
		}
	}

	// 处理响应
	data, err := processResponse(resp)
	if err != nil {
		return DispatchResp{}, err
	}

	// 解析JSON响应
	var dr DispatchResp
	if err := json.Unmarshal(data, &dr); err != nil {
		return DispatchResp{}, fmt.Errorf("解析 dispatcher JSON 失败: %v", err)
	}

	// 检查响应code
	if dr.Code != 0 {
		return dr, fmt.Errorf("dispatcher 返回 code=%d, message=%s", dr.Code, dr.Message)
	}

	return dr, nil
}
