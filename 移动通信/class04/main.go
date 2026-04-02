package main

import (
	"fmt"
	"strings"
)

func xorOperation(dividend string, divisor string) string {
	result := ""
	for i := 0; i < len(dividend); i++ {
		if dividend[i] == divisor[i] {
			result += "0"
		} else {
			result += "1"
		}
	}
	return result
}

func crcEncode(data string, generator string) string {
	n := len(generator) - 1
	augmentedData := data + strings.Repeat("0", n)

	currentDividend := augmentedData[:len(generator)]
	for i := 0; i < len(data); i++ {
		if currentDividend[0] == '1' {
			currentDividend = xorOperation(currentDividend, generator)
		} else {
			currentDividend = xorOperation(currentDividend, strings.Repeat("0", len(generator)))
		}

		if i+len(generator) < len(augmentedData) {
			currentDividend = currentDividend[1:] + string(augmentedData[i+len(generator)])
		} else {
			currentDividend = currentDividend[1:]
		}
	}

	remainder := currentDividend
	return data + remainder
}

func crcCheck(receivedData string, generator string) bool {
	currentDividend := receivedData[:len(generator)]
	for i := 0; i < len(receivedData)-len(generator)+1; i++ {
		if currentDividend[0] == '1' {
			currentDividend = xorOperation(currentDividend, generator)
		} else {
			currentDividend = xorOperation(currentDividend, strings.Repeat("0", len(generator)))
		}

		if i+len(generator) < len(receivedData) {
			currentDividend = currentDividend[1:] + string(receivedData[i+len(generator)])
		} else {
			currentDividend = currentDividend[1:]
		}
	}

	return currentDividend == strings.Repeat("0", len(generator)-1)
}

func CRC() {
	for {
		fmt.Println(strings.Repeat("=", 30))
		var mode string
		fmt.Print("请输入运行模式（0：退出，1：CRC编码，2：CRC验证）：")
		fmt.Scanln(&mode)
		if mode == "0" {
			break
		} else if mode == "1" {
			var data, generator string
			fmt.Print("请输入数据（例如：10110011）：")
			fmt.Scanln(&data)
			fmt.Print("请输入生成多项式（例如：11001）：")
			fmt.Scanln(&generator)
			encodedData := crcEncode(data, generator)
			fmt.Printf("编码结果：%s\n", encodedData)
		} else if mode == "2" {
			var receivedData, generator string
			fmt.Print("请输入接收数据：")
			fmt.Scanln(&receivedData)
			fmt.Print("请输入生成多项式：")
			fmt.Scanln(&generator)
			if crcCheck(receivedData, generator) {
				fmt.Println("验证成功：数据未出错。")
			} else {
				fmt.Println("验证失败：数据有错误！")
			}
		} else {
			fmt.Println("无效输入，请重新选择模式！")
		}
	}
}

// 卷积编码器
type ConvolutionalEncoder struct {
	register []byte // 移位寄存器，长度为生成多项式长度减1
	G1       []byte // 生成多项式G1
	G2       []byte // 生成多项式G2
}

// 新建卷积编码器
func NewConvolutionalEncoder(g1, g2 []byte) *ConvolutionalEncoder {
	register := make([]byte, len(g1))
	encoder := &ConvolutionalEncoder{
		register: register,
		G1:       g1,
		G2:       g2,
	}
	return encoder
}

func (encoder *ConvolutionalEncoder) enqueue(bit byte) {
	// 依次将数组元素向左移动一位
	for i := 0; i < len(encoder.register)-1; i++ {
		encoder.register[i] = encoder.register[i+1]
	}
	// 将新的比特放入队列尾部
	encoder.register[len(encoder.register)-1] = bit
}

// 编码单个比特
func (encoder *ConvolutionalEncoder) Encode() string {
	g1 := make([]byte, 3)
	g2 := make([]byte, 3)
	for i, r := range encoder.register {
		if encoder.G1[i] == 1 {
			g1 = append(g1[1:], r)
		}
		if encoder.G2[i] == 1 {
			g2 = append(g2[1:], r)
		}
	}
	g1Result := byte(0)
	for i, bit := range g1 {
		if i == 0 {
			g1Result = bit
			continue
		}
		g1Result ^= bit
	}

	// 对 g2 中的元素进行异或操作
	g2Result := byte(0)
	for i, bit := range g2 {
		if i == 0 {
			g2Result = bit
			continue
		}
		g2Result ^= bit
	}
	return fmt.Sprintf("%b%b", g1Result, g2Result)
}

func stringToBytes(s string) []byte {
	var result []byte
	for _, c := range s {
		// 将每个字符转换为 0 或 1，并添加到字节数组
		if c == '1' {
			result = append(result, 1)
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func junaji() {
	var inputStr string
	var g1Str, g2Str string

	// 手动输入生成多项式 G1, G2 和输入数据
	fmt.Println("请输入生成多项式 G1 (例如: 111): ")
	fmt.Scanln(&g1Str)
	fmt.Println("请输入生成多项式 G2 (例如: 101): ")
	fmt.Scanln(&g2Str)
	fmt.Println("请输入输入数据 (例如: 101): ")
	fmt.Scanln(&inputStr)

	// 将输入的字符串转为字节数组
	input := stringToBytes(inputStr)
	G1 := stringToBytes(g1Str)
	G2 := stringToBytes(g2Str)

	encoder := NewConvolutionalEncoder(G1, G2)

	var ans string
	for _, bit := range input {
		encoder.enqueue(bit)
		ans += encoder.Encode()
	}
	// 填充寄存器，以完成编码过程
	for i := 0; i < len(encoder.register)-1; i++ {
		encoder.enqueue(0)
		ans += encoder.Encode()
	}

	// 打印输入和输出
	fmt.Printf("输入: %v\n输出: %s\n", input, ans)
}
func main() {
	fmt.Print("0:CRC编码，1:卷积：")
	var cho string
	fmt.Scanln(&cho)
	if cho == "0" {
		CRC()
	} else if cho == "1" {
		junaji()
	} else {
		fmt.Errorf("未知选项")
	}
}
