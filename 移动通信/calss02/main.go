package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// 映射规则：00 -> -3, 01 -> -1, 10 -> 1, 11 -> 3
var mapping = map[string]int{
	"00": -3,
	"01": -1,
	"10": 1,
	"11": 3,
}

// 计算能量 E = I^2 + Q^2
func calculateEnergy(I, Q int) int {
	return I*I + Q*Q
}

// 归一化 IQ 值
func normalizeIQ(a, avgPower float64) float64 {
	// 除以根号平均功率
	normalizer := 1 / math.Sqrt(avgPower)
	return a * normalizer
}
func floatToBinary(value float64, precision int) string {
	// 将浮点数转换为整数
	value = math.Abs(value)
	value *= math.Pow(2, float64(precision))
	intValue := int(value)

	// 转换为二进制字符串
	binary := strconv.FormatInt(int64(intValue), 2)

	// 确保二进制字符串长度为指定的精度
	for len(binary) < precision {
		binary = "0" + binary
	}

	return binary
}

// 添加噪声：小范围的噪声，模拟偏移
func addNoise(value float64, noiseLevel float64) float64 {
	// 使用均匀分布生成噪声
	return value + (rand.Float64()*2-1)*noiseLevel
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 协议为 WebSocket 协议
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 随机数种子
	rand.Seed(time.Now().UnixNano())

	// 定义符号
	binarySymbols := []string{
		"0011", "0001", "0101", "0111", "1001", "1011", "1111", "1101",
		"1000", "1010", "1110", "1100", "0110", "0100", "0000", "0010",
	}

	totalEnergy := 0
	numSymbols := len(binarySymbols)

	// 计算总能量
	for _, symbol := range binarySymbols {
		realPart := symbol[:2]
		imagPart := symbol[2:]
		IValue := mapping[realPart]
		QValue := mapping[imagPart]
		energy := calculateEnergy(IValue, QValue)
		totalEnergy += energy
	}

	avgPower := float64(totalEnergy) / float64(numSymbols)

	// 控制信号变化频率
	for {
		for _, symbol := range binarySymbols {
			realPart := symbol[:2]
			imagPart := symbol[2:]
			IValue := mapping[realPart]
			QValue := mapping[imagPart]

			// 归一化 IQ 值
			normalizedI := normalizeIQ(float64(IValue), avgPower)
			normalizedQ := normalizeIQ(float64(QValue), avgPower)

			// 添加噪声
			noiseLevel := 0.1 // 控制噪声级别
			normalizedI = addNoise(normalizedI, noiseLevel)
			normalizedQ = addNoise(normalizedQ, noiseLevel)
			binaryI := floatToBinary(normalizedI, 10)
			binaryQ := floatToBinary(normalizedQ, 10)
			// 将数据打包并发送给客户端
			data := map[string]interface{}{
				"symbol":   symbol,
				"I":        normalizedI,
				"Q":        normalizedQ,
				"I_binary": binaryI,
				"Q_binary": binaryQ,
			}

			err := conn.WriteJSON(data)
			if err != nil {
				fmt.Println("Write error:", err)
				return
			}

			// 控制发送频率

		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 WebSocket 服务器
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
