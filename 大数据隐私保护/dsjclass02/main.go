package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/fernet/fernet-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type SSE struct {
	K1        []byte            // 密钥 K1
	K2        []byte            // 密钥 K2
	K3        []byte            // 密钥 K3
	ArrayA    map[int][]byte    // 存储加密节点的数组
	TableT    map[int][3]interface{} // 查找表
}

func NewSSE() *SSE {
	return &SSE{
		ArrayA: make(map[int][]byte),
		TableT: make(map[int][3]interface{}),
	}
}

// Setup算法：生成密钥
func (sse *SSE) Setup() {
	sse.K1 = generateKey()
	sse.K2 = generateKey()
	sse.K3 = generateKey()
}

// 生成一个新的密钥
func generateKey() []byte {
	key := make([]byte, 32)
	_, err := os.ReadFile("/dev/urandom", key)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// 伪随机置换函数
func (sse *SSE) PseudoRandomPerm(key []byte, value string) int {
	combined := append(key, []byte(value)...)
	hash := sha256.Sum256(combined)
	return int(binary.BigEndian.Uint32(hash[:4])) % 10000
}

// 伪随机函数
func (sse *SSE) PseudoRandomFunc(key []byte, value string) []byte {
	combined := append(key, []byte(value)...)
	hash := sha256.Sum256(combined)
	return hash[:]
}

// BuildIndex算法：构建索引
func (sse *SSE) BuildIndex(documents map[string]map[string]bool) {
	ctr := 1
	for keyword, docIDs := range documents {
		docIDList := make([]string, 0, len(docIDs))
		for docID := range docIDs {
			docIDList = append(docIDList, docID)
		}
		KiPrev := generateKey()

		// 随机生成 gamma 和 eta
		gamma := sse.PseudoRandomPerm(sse.K3, keyword)
		eta := sse.PseudoRandomFunc(sse.K2, keyword)

		firstAddr := sse.PseudoRandomPerm(sse.K1, fmt.Sprintf("%d", ctr))
		sse.TableT[gamma] = [3]interface{}{firstAddr, KiPrev, eta}

		for j, docID := range docIDList {
			KiCurr := generateKey()
			var nextAddr int
			if j < len(docIDList)-1 {
				nextAddr = sse.PseudoRandomPerm(sse.K1, fmt.Sprintf("%d", ctr+1))
			}

			node := map[string]interface{}{
				"doc_id":   docID,
				"next_key": KiCurr,
				"next_addr": nextAddr,
			}

			// 加密节点
			fernetKey := fernet.Key(KiPrev)
			encryptedNode, err := fernetKey.Encrypt([]byte(fmt.Sprintf("%v", node)))
			if err != nil {
				log.Fatal(err)
			}

			sse.ArrayA[sse.PseudoRandomPerm(sse.K1, fmt.Sprintf("%d", ctr))] = encryptedNode

			KiPrev = KiCurr
			ctr++
		}
	}
}

// GenTrapdoor算法：生成陷门
func (sse *SSE) GenTrapdoor(keyword string) (int, []byte) {
	gamma := sse.PseudoRandomPerm(sse.K3, keyword)
	eta := sse.PseudoRandomFunc(sse.K2, keyword)
	return gamma, eta
}

// Search算法：搜索文档
func (sse *SSE) Search(trapdoor [2]interface{}) []string {
	gamma := trapdoor[0].(int)
	eta := trapdoor[1].([]byte)

	var result []string
	if val, exists := sse.TableT[gamma]; exists {
		addr, key, storedEta := val[0].(int), val[1].([]byte), val[2].([]byte)
		if !bytes.Equal(storedEta, eta) {
			return result
		}

		for addr != nil {
			fernetKey := fernet.Key(key)
			encryptedNode := sse.ArrayA[addr]
			decryptedNode, err := fernetKey.Decrypt(encryptedNode)
			if err != nil {
				log.Fatal(err)
			}

			node := map[string]interface{}
			_ = json.Unmarshal(decryptedNode, &node)
			result = append(result, node["doc_id"].(string))

			addr = node["next_addr"]
			key = node["next_key"]
		}
	}
	return result
}

// 从文件中构建文档索引
func BuildDocumentsFromFiles() map[string]map[string]bool {
	documents := make(map[string]map[string]bool)
	documentsFolder := "D:\\study\\python-learn\\Big_Data\\file_test2"

	files, err := ioutil.ReadDir(documentsFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		content, err := ioutil.ReadFile(documentsFolder + "\\" + file.Name())
		if err != nil {
			log.Println("读取文件错误:", err)
			continue
		}

		keywords := ExtractKeywords(string(content))
		for _, keyword := range keywords {
			if _, exists := documents[keyword]; !exists {
				documents[keyword] = make(map[string]bool)
			}
			documents[keyword][file.Name()] = true
		}
	}

	return documents
}

// 从文本内容中提取关键词
func ExtractKeywords(content string) []string {
	// 这里只做示例，可以扩展为你自己的TF-IDF或TextRank方法
	return []string{"example", "keyword", "test"} // 这里只是返回示例
}

func main() {
	// 从文件构建文档集合
	documents := BuildDocumentsFromFiles()

	if len(documents) == 0 {
		fmt.Println("没有成功构建文档索引")
		return
	}

	// 初始化并运行SSE
	sse := NewSSE()
	sse.Setup()

	// 构建索引
	sse.BuildIndex(documents)

	for {
		var keyword string
		fmt.Print("请输入要搜索的关键词 (输入 'q' 退出): ")
		fmt.Scanln(&keyword)
		if keyword == "q" {
			break
		}

		// 生成陷门并搜索
		trapdoor := sse.GenTrapdoor(keyword)
		result := sse.Search(trapdoor)

		if len(result) > 0 {
			fmt.Printf("\n找到以下文件包含关键词 '%s':\n", keyword)
			for i, file := range result {
				fmt.Printf("%d. %s\n", i+1, file)
			}
		} else {
			fmt.Printf("\n未找到包含关键词 '%s' 的文件\n", keyword)
		}
	}
}
