package main

import (
	"fmt"
	"log"
	"os"
	"time"

	packsender "boss/PackSender"
	"boss/login"
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// 从目标网站获取职位数据
func fetchJobs(lastCity string) error {
	if err := packsender.FetchJobs(lastCity); err != nil {
		return fmt.Errorf("获取职位失败: %v", err)
	}
	return nil
}

// 获取并保存 cookie
func getCookies() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		log.Fatal(err)
	}

	// 安全验证页 + 回调 URL
	securityURL := "https://www.zhipin.com/web/common/security-check.html?seed=2XEA8YBKDjNGiMnmTH0oLc3yvVyUHPQMsQ%2BofNIf%2BqQ%3D&name=ec959783&ts=1746962468667&callbackUrl=https%3A%2F%2Fwww.zhipin.com%2Fweb%2Fgeek%2Fjob%3Fquery%3D%E7%BD%91%E7%BB%9C%E5%AE%89%E5%85%A8%26city%3D101230100"
	targetURL := "https://www.zhipin.com/web/geek/job?query=网络安全&city=101230100"
	var cookies []*network.Cookie

	err := chromedp.Run(ctx,
		chromedp.Navigate(securityURL), // 先访问安全页
		chromedp.Sleep(3*time.Second),  // 等待跳转或执行校验
		chromedp.Navigate(targetURL),   // 再访问目标页
		chromedp.Sleep(5*time.Second),  // 等待页面加载完成
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		return fmt.Errorf("获取 cookie 失败: %v", err)
	}

	// 筛选你需要保存的字段
	newCookies := make(map[string]string)
	for _, c := range cookies {
		if c.Name == "__zp_stoken__" {
			fmt.Println("__zp_stoken__ =", c.Value)
			newCookies[c.Name] = c.Value
		}
	}

	if len(newCookies) == 0 {
		return fmt.Errorf("未获取到目标 Cookie，跳过写入")
	}

	// 读取已有 cookie.json
	cookieFile := "cookies.json"
	existingCookies := make(map[string]string)

	file, err := os.OpenFile(cookieFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打开 cookie.json 文件失败: %v", err)
	}
	defer file.Close()

	if stat, _ := file.Stat(); stat.Size() > 0 {
		if err := json.NewDecoder(file).Decode(&existingCookies); err != nil {
			return fmt.Errorf("解析 cookie.json 失败: %v", err)
		}
	}

	// 合并新的 cookie，不覆盖已有的
	for k, v := range newCookies {
		if _, exists := existingCookies[k]; !exists {
			existingCookies[k] = v
		}
	}

	// 重写 cookie.json
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingCookies); err != nil {
		return fmt.Errorf("写入 cookie.json 失败: %v", err)
	}

	fmt.Println("cookie.json 更新完成")
	return nil
}

// 调用 Node.js 脚本获取 fp 参数
func getFP() error {
	cmd := exec.Command("node", "get-fp.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("正在调用 Node.js 脚本获取 fp 参数...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Node.js 执行失败: %v", err)
	}

	fpBytes, err := os.ReadFile("fp.txt")
	if err != nil {
		return fmt.Errorf("读取 fp 失败: %v", err)
	}
	fp := strings.TrimSpace(string(fpBytes))
	fmt.Println("使用 fp:", fp)

	// 这里可以根据需要进一步处理 fp
	return nil
}

func main() {
	// 先更新 cookie
	if err := getCookies(); err != nil {
		fmt.Println("获取 Cookie 错误:", err)
		return
	}

	// 然后执行 Node.js 获取 fp
	if err := getFP(); err != nil {
		fmt.Println("获取 FP 错误:", err)
		return
	}

	var l login.Login
	fpBytes, err := os.ReadFile("fp.txt")
	if err != nil {
		log.Fatalf("读取 fp 失败: %v", err)
	}
	fp := strings.TrimSpace(string(fpBytes))
	fmt.Println("使用 fp:", fp)

	// 初始化 login，设置页面标识
	if err := l.Init("header-login", fp); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	// 获取二维码
	if err := l.GetLoginQr(); err != nil {
		panic(err)
	}

	// 轮询登录，设置 5 分钟超时
	if err := l.PollQrLogin(30 * time.Hour); err != nil {
		panic(err)
	}

	// 输出登录结果
	fmt.Printf("登录成功")

	// 获取职位数据
	if err := fetchJobs("101230100"); err != nil {
		fmt.Println("获取职位数据失败:", err)
		return
	}

	fmt.Println("所有操作完成")
}
