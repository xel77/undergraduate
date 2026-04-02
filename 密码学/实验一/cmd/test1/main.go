// caesar.go
package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// ---------- helpers ----------
func toUpperLetters(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToUpper(r))
		}
	}
	return b.String()
}

func mod(a, m int) int {
	a %= m
	if a < 0 {
		a += m
	}
	return a
}

// ---------- Caesar ----------
func caesarEncryptPlainUpper(plain string, key int) string {
	key = mod(key, 26)
	pt := toUpperLetters(plain)
	var b strings.Builder
	for _, r := range pt {
		idx := int(r - 'A')
		b.WriteByte(byte(mod(idx+key, 26) + 'A'))
	}
	return b.String()
}

func caesarDecryptUpper(cipher string, key int) string {
	// decrypt by shifting backwards
	return caesarEncryptPlainUpper(cipher, -key)
}

func caesarBrute(cipher string) []string {
	ct := toUpperLetters(cipher)
	res := make([]string, 26)
	for k := 0; k < 26; k++ {
		var b strings.Builder
		for _, r := range ct {
			idx := int(r - 'A')
			b.WriteByte(byte(mod(idx-k, 26) + 'A'))
		}
		res[k] = b.String()
	}
	return res
}

// English letter frequency (percent)
var engFreq = map[rune]float64{
	'A': 8.167, 'B': 1.492, 'C': 2.782, 'D': 4.253, 'E': 12.702,
	'F': 2.228, 'G': 2.015, 'H': 6.094, 'I': 6.966, 'J': 0.153,
	'K': 0.772, 'L': 4.025, 'M': 2.406, 'N': 6.749, 'O': 7.507,
	'P': 1.929, 'Q': 0.095, 'R': 5.987, 'S': 6.327, 'T': 9.056,
	'U': 2.758, 'V': 0.978, 'W': 2.360, 'X': 0.150, 'Y': 1.974, 'Z': 0.074,
}

func chiSquaredScore(text string) float64 {
	text = toUpperLetters(text)
	N := float64(len(text))
	if N == 0 {
		return 1e9
	}
	counts := make(map[rune]float64)
	for _, r := range text {
		counts[r]++
	}
	var chi float64
	for ch, freq := range engFreq {
		expected := freq/100.0 * N
		obs := counts[ch]
		if expected > 0 {
			d := obs - expected
			chi += d * d / expected
		}
	}
	return chi
}

// frequency attack: try all shifts and pick the one with lowest chi-squared
func caesarFreqAttack(cipher string) (bestKey int, bestPlain string, bestScore float64) {
	bestScore = math.Inf(1)
	bestKey = 0
	bestPlain = ""
	for k := 0; k < 26; k++ {
		pt := caesarDecryptUpper(cipher, k)
		score := chiSquaredScore(pt)
		if score < bestScore {
			bestScore = score
			bestKey = k
			bestPlain = pt
		}
	}
	return bestKey, bestPlain, bestScore
}

// naive pretty spacing for known phrases (just for display)
func prettySpacing(s string) string {
	words := []string{
		"IF", "YOU", "CAN", "READ", "THIS", "THANK", "A", "TEACHER",
		"LOOK", "UP", "IN", "THE", "AIR", "IT", "S", "BIRD", "PLANE", "SUPERMAN",
		"ALGORITHMS", "ARE", "QUITE", "GENERAL", "DEFINITIONS", "OF", "ARITHMETIC", "PROCESSES",
	}
	res := s
	for _, w := range words {
		res = strings.ReplaceAll(res, w, " "+w+" ")
	}
	return strings.Join(strings.Fields(res), " ")
}

func main() {
	fmt.Println("===  两个同学互相加密 ===")
	plain := "security"
	// 学生 A 选 keyA，加密发送给学生 B
	keyA := 7 // 示例
	cipherA := caesarEncryptPlainUpper(plain, keyA)
	fmt.Printf("学生A 明文: %s  选 keyA=%d  加密后密文: %s\n", plain, keyA, cipherA)

	// 学生 B 收到密文再用自己的 keyB 再加密发回（等效总key = keyA + keyB mod26）
	keyB := 11 // 示例
	cipherAB := caesarEncryptPlainUpper(cipherA, keyB)
	fmt.Printf("学生B 在收到的密文上再用 keyB=%d 加密 -> %s\n", keyB, cipherAB)
	fmt.Printf("等效总 key = (keyA + keyB) mod 26 = %d\n", mod(keyA+keyB, 26))


	fmt.Println("===  穷尽密钥攻击 ===")
	target := "BEEAKFYDJXUQYHYJIQRYHTYJIQFBQDUYJIIKFUHCQD"
	fmt.Println("截获密文:", target)
	fmt.Println("\n列出所有 26 个候选 (k -> 明文)：")
	cands := caesarBrute(target)
	for k, cand := range cands {
		fmt.Printf("k=%2d -> %s\n", k, cand)
	}
	// 简单人工筛选含常用单词的候选
	fmt.Println("\n人工关键词筛选的可能通顺候选 (含 THE/LOOK/YOU/SUPERMAN 等关键词)：")
	for k, cand := range cands {
		if strings.Contains(cand, "THE") || strings.Contains(cand, "LOOK") || strings.Contains(cand, "YOU") || strings.Contains(cand, "SUPERMAN") {
			fmt.Printf("k=%2d -> %s\n", k, cand)
		}
	}
	fmt.Println("\n=== 频率统计攻击（卡方） ===")
	target2 := "FQLTWNYMRXFWJVZNYJLJSJWFQIJKNSNYNTSXTKFWNYMRJYNHUWTHJXXJX"
	fmt.Println("截获密文:", target2)
	bestK, bestPT, score := caesarFreqAttack(target2)
	fmt.Printf("卡方频率攻击结果: 最优 key = %d  score=%.2f\n", bestK, score)
	fmt.Printf("明文 (无空格): %s\n", bestPT)
	fmt.Printf("明文（尝试加空格以便阅读）: %s\n", prettySpacing(bestPT))
}
