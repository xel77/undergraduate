//go:generate rsrc -ico resource/icon.ico -manifest resource/info.manifest -o desktop.syso
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Book struct {
	Title        string
	Author       string
	AuthorNation string
	PublishDate  string
	Publisher    string
	Price        string
	Binding      string
	Rating       float64
	RatingCount  int
	CoverURL     string
	CoverPath    string // 本地封面存储路径
}

func main() {
	baseURL := "https://book.douban.com/latest"
	var books []Book

	// 检测 CSV 是否存在
	csvFilename := "douban_books.csv"
	if _, err := os.Stat(csvFilename); err == nil {
		fmt.Printf("检测到 %s 文件已存在，是否覆盖？(y/n): ", csvFilename)
		input := getUserInput()
		if strings.ToLower(input) != "y" && strings.ToLower(input) != "yes" {
			fmt.Println("终止爬取。")
			return
		}
	}

	autoContinue := false

	for {
		doc, err := fetchHTML(baseURL)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".media.clearfix").Each(func(i int, s *goquery.Selection) {
			book := parseBookItem(s, baseURL)
			books = append(books, book)
		})

		fmt.Printf("\n共解析到 %d 本书籍\n", len(books))

		// 保存数据到 CSV
		if err := saveToCSV(csvFilename, books); err != nil {
			log.Fatalf("保存 CSV 失败: %v", err)
		} else {
			fmt.Printf("数据已成功存入 %s\n", csvFilename)
		}

		// 查找下一页
		nextPage := getNextPageURL(doc, baseURL)
		if nextPage == "" {
			fmt.Println("没有找到下一页，爬取完成。")
			break
		}

		if !autoContinue {
			fmt.Print("检测到下一页，是否继续爬取？(y/n, 输入yy自动爬取到结束): ")
			input := getUserInput()
			if input == "yy" {
				autoContinue = true
			} else if strings.ToLower(input) != "y" && strings.ToLower(input) != "yes" {
				fmt.Println("终止爬取。")
				break
			}
		}

		baseURL = nextPage
	}
	err := runPythonScript("app.py")
	if err != nil {
		fmt.Printf("调用 Python 失败: %v\n", err)
	} else {
		fmt.Println("Python 词云脚本执行成功！")
	}
}

func runPythonScript(script string) error {
	cmd := exec.Command("python", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseBookItem(s *goquery.Selection, baseURL string) Book {
	title := strings.TrimSpace(s.Find(".fleft").First().Text())
	pubInfo := parsePublicationInfo(s.Find(".subject-abstract.color-gray").Text())
	coverURL := resolveURL(baseURL, s.Find("img").AttrOr("src", ""))
	coverPath := "covers/" + filepath.Base(coverURL)

	return Book{
		Title:        title,
		Author:       pubInfo["author"],
		AuthorNation: pubInfo["nation"],
		PublishDate:  pubInfo["date"],
		Publisher:    pubInfo["publisher"],
		Price:        pubInfo["price"],
		Binding:      pubInfo["binding"],
		Rating:       parseRating(s.Find(".font-small.fleft")),
		RatingCount:  parseRatingCount(s.Find(".fleft.ml8.color-gray").Text()),
		CoverURL:     coverURL,
		CoverPath:    coverPath,
	}
}

func parsePublicationInfo(text string) map[string]string {
	info := make(map[string]string)
	parts := regexp.MustCompile(`\s*/\s*`).Split(strings.TrimSpace(text), -1)

	dateIndex := -1
	datePattern := regexp.MustCompile(`^\d{4}(-\d{1,2})?(-\d{1,2})?$`)

	for i := len(parts) - 1; i >= 0; i-- {
		if datePattern.MatchString(parts[i]) {
			dateIndex = i
			break
		}
	}

	if dateIndex == -1 {
		return map[string]string{"author": "", "nation": "", "date": "", "publisher": "", "price": ""}
	}

	authorParts := parts[:dateIndex]
	authorNames, nations := []string{}, []string{}
	for _, part := range authorParts {
		author, nation := parseAuthorNation(part)
		authorNames = append(authorNames, author)
		nations = append(nations, nation)
	}
	info["author"] = strings.Join(authorNames, " / ")
	info["nation"] = strings.Join(nations, " / ")
	info["date"] = parts[dateIndex]

	if len(parts) > dateIndex+1 {
		info["publisher"] = parts[dateIndex+1]
	} else {
		info["publisher"] = ""
	}

	if len(parts) > dateIndex+2 {
		info["price"] = parts[dateIndex+2]
	} else {
		info["price"] = ""
	}

	return info
}

func parseAuthorNation(authorStr string) (author, nation string) {
	re := regexp.MustCompile(`^[【\[\(]\s*([\p{L}]+)\s*[\]\)】]\s*(.+)$`)
	if matches := re.FindStringSubmatch(authorStr); len(matches) == 3 {
		return strings.TrimSpace(matches[2]), strings.TrimSpace(matches[1])
	}

	if regexp.MustCompile(`[a-zA-Z]`).MatchString(authorStr) {
		return strings.TrimSpace(authorStr), "外国"
	}

	return strings.TrimSpace(authorStr), "中国"
}

func parseRating(s *goquery.Selection) float64 {
	text := strings.TrimSpace(s.Text())
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	if matches := re.FindString(text); matches != "" {
		rating, _ := strconv.ParseFloat(matches, 64)
		return rating
	}
	return 0.0
}

func parseRatingCount(text string) int {
	if matches := regexp.MustCompile(`(\d+)`).FindStringSubmatch(text); len(matches) > 0 {
		count, _ := strconv.Atoi(matches[1])
		return count
	}
	return 0
}

func resolveURL(base, path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	baseURL, _ := url.Parse(base)
	relURL, _ := url.Parse(path)
	return baseURL.ResolveReference(relURL).String()
}

func fetchHTML(url string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromReader(resp.Body)
}

func getNextPageURL(doc *goquery.Document, baseURL string) string {
	nextSelection := doc.Find("span.next a")
	if nextSelection.Length() > 0 {
		if nextHref, exists := nextSelection.Attr("href"); exists {
			return resolveURL(baseURL, nextHref)
		}
	}
	return ""
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// 保存数据到 CSV
func saveToCSV(filename string, books []Book) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"书名", "作者", "作者国别", "出版时间", "出版社", "价格", "评分", "评价人数", "封面URL", "封面路径"})

	for _, book := range books {
		writer.Write([]string{book.Title, book.Author, book.AuthorNation, book.PublishDate, book.Publisher, book.Price, strconv.FormatFloat(book.Rating, 'f', 1, 64), strconv.Itoa(book.RatingCount), book.CoverURL, book.CoverPath})
	}

	return nil
}

// 下载封面图片
func downloadCovers(books []Book) {
	os.MkdirAll("covers", os.ModePerm)
	for _, book := range books {
		resp, err := http.Get(book.CoverURL)
		if err == nil {
			data, _ := os.Create(book.CoverPath)
			defer data.Close()
			defer resp.Body.Close()
			_, _ = data.ReadFrom(resp.Body)
		}
	}
}
