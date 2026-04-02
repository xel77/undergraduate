// affine.go
package main

import (
	"fmt"
	"strings"
	"unicode"
)

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

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func invMod(a, m int) int {
	a = mod(a, m)
	if gcd(a, m) != 1 {
		return -1
	}
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x
		}
	}
	return -1
}

// Affine encrypt: C = (a * P + b) mod 26
func affineEncrypt(plain string, a, b int) string {
	if invMod(a, 26) == -1 {
		return "ERROR: a not invertible mod26"
	}
	pt := toUpperLetters(plain)
	var sb strings.Builder
	for _, r := range pt {
		x := int(r - 'A')
		c := mod(a*x+b, 26)
		sb.WriteByte(byte(c + 'A'))
	}
	return sb.String()
}

// Affine decrypt: P = a_inv * (C - b) mod 26
func affineDecrypt(cipher string, a, b int) string {
	aInv := invMod(a, 26)
	if aInv == -1 {
		return "ERROR: a not invertible mod26"
	}
	ct := toUpperLetters(cipher)
	var sb strings.Builder
	for _, r := range ct {
		y := int(r - 'A')
		p := mod(aInv*(y-b), 26)
		sb.WriteByte(byte(p + 'A'))
	}
	return sb.String()
}

// 从两对明密文字母求 a,b。如果有解且 a 与 26 互质则返回 (a,b,true)
func affineSolveFromPairs(p0, c0, p1, c1 byte) (int, int, bool) {
	modn := 26
	P0 := int(p0 - 'A')
	P1 := int(p1 - 'A')
	C0 := int(c0 - 'A')
	C1 := int(c1 - 'A')
	dP := mod((P0 - P1), modn)
	dC := mod((C0 - C1), modn)
	// 寻找 a 使得 a * dP ≡ dC (mod26)，且 gcd(a,26)=1
	for a := 0; a < modn; a++ {
		if gcd(a, modn) != 1 {
			continue
		}
		if mod(a*dP-dC, modn) == 0 {
			b := mod(C0-a*P0, modn)
			return a, b, true
		}
	}
	return 0, 0, false
}

func main() {
	// 1.(1) 两人分组示例：security
	plain := "security"
	// 学生 A 的密钥 (示例)
	a1, b1 := 5, 8
	// 学生 B 的密钥 (示例)
	a2, b2 := 3, 7

	ctA := affineEncrypt(plain, a1, b1)
	ctAB := affineEncrypt(ctA, a2, b2)
	// 合成等效密钥： a_total = a2*a1 mod26; b_total = a2*b1 + b2 mod26
	aTotal := mod(a2*a1, 26)
	bTotal := mod(a2*b1+b2, 26)
	ctTotal := affineEncrypt(plain, aTotal, bTotal)

	fmt.Printf("明文: %s\n", plain)
	fmt.Printf("学生A 用 a=%d b=%d 加密 -> %s\n", a1, b1, ctA)
	fmt.Printf("学生B 在收到的密文上用 a=%d b=%d 再加密 -> %s\n", a2, b2, ctAB)
	fmt.Printf("等效合成 a_total=%d b_total=%d, 直接用等效密钥加密 -> %s\n\n", aTotal, bTotal, ctTotal)

	// 2.(2) 已知明文攻击
	ctGiven := "EDSGICKXHUKLZVEQZVKXWKZUKCVUH"
	// 已知明文前两字符 "IF"
	p0, p1 := byte('I'), byte('F')
	// 密文前两字符
	c0, c1 := byte(ctGiven[0]), byte(ctGiven[1])

	aSolved, bSolved, ok := affineSolveFromPairs(p0, c0, p1, c1)
	fmt.Printf("给定密文: %s\n", ctGiven)
	if !ok {
		fmt.Println("无法从已知前两字母推出可逆仿射密钥 (a,b)。")
		return
	}
	plainRecovered := affineDecrypt(ctGiven, aSolved, bSolved)
	fmt.Printf("已知明文前两字母: IF\n")
	fmt.Printf("求得仿射密钥 a=%d b=%d\n", aSolved, bSolved)
	fmt.Printf("解密得到明文: %s\n", plainRecovered)
}
