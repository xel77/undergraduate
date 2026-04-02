package utili

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func init() {}

func Handle(r map[string]interface{}) map[string]string {
	if r["ipmethon"] == "1" { // 单个 IP 扫描
		return onePortHand(r)
	} else { // IP 范围扫描
		return rangeScan(r)
	}
}

func onePortHand(r map[string]interface{}) map[string]string {
	results := make(map[string]string)

	scanMethod, err := strconv.Atoi(r["scanmethon"].(string))
	if err != nil {
		fmt.Println("Invalid scan method:", err)
		return nil
	}

	if r["portmethon"] == "1" { // 全端口
		if scanMethod == 1 { // TCP
			return tcpConnectScan(r["ip"].(string))
		} else if scanMethod == 2 { // SYN
			return synScan(r["ip"].(string))
		} else { // FIN
			return finScan(r["ip"].(string))
		}
	} else { // 单个端口
		port, err := strconv.Atoi(r["port"].(string))
		if err != nil {
			fmt.Println("错误端口:", err)
			return nil
		}

		return singlePortScan(r["ip"].(string), port, scanMethod, results)
	}
}

// rangeScan 执行 IP 范围扫描
func rangeScan(r map[string]interface{}) map[string]string {
	results := make(map[string]string)
	_, ipNet, err := net.ParseCIDR(r["ip"].(string))
	if err != nil {
		fmt.Println("错误子网:", err)
		return nil
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)
	var mu sync.Mutex

	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		wg.Add(1)
		sem <- struct{}{}
		go func(currentIP net.IP) {
			defer wg.Done()
			defer func() { <-sem }()

			if r["portmethon"] == "1" { // 全端口扫描
				if scanMethod, ok := r["scanmethon"].(string); ok {
					switch scanMethod {
					case "1": // TCP
						tcpConnectScan(currentIP.String())
					case "2": // SYN
						synScan(currentIP.String())
					case "3": // FIN
						finScan(currentIP.String())
					default:
						fmt.Println("未知的扫描方式")
					}
				} else {
					fmt.Println("scanmethon 类型错误: 期望字符串类型")
					return
				}
			} else { // 单个端口扫描
				if portStr, ok := r["port"].(string); ok {
					port, err := strconv.Atoi(portStr)
					if err != nil {
						fmt.Println("错误端口:", err)
						return
					}

					if scanMethodstr, ok := r["scanmethon"].(string); ok {
						mu.Lock()
						scanMethodInt, err := strconv.Atoi(scanMethodstr)
						if err == nil {

							singlePortScan(currentIP.String(), port, scanMethodInt, results)
						}
						mu.Unlock()
					} else {
						fmt.Println("scanmethon 类型错误: 期望整数类型")
					}
				} else {
					fmt.Println("port 类型错误")
					return
				}
			}
		}(ip)
	}

	wg.Wait()
	return results
}

// inc 增加 IP 地址
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

// singlePortScan 执行单个端口的扫描
func singlePortScan(ip string, port int, scanMethod int, results map[string]string) map[string]string {
	address := fmt.Sprintf("%s:%d", ip, port)
	if scanMethod == 1 { // TCP
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err == nil {
			results[address] = "开启\n"
			conn.Close()
		}
	} else if scanMethod == 2 { // SYN
		sendSYN(ip, port, 3*time.Second, results, &sync.Mutex{})
	} else if scanMethod == 3 { // FIN
		sendFIN(address, 3*time.Second, results, &sync.Mutex{})
	}
	return results // 返回扫描结果
}

// tcpConnectScan 使用 TCP 连接扫描端口
func tcpConnectScan(ip string) map[string]string {
	var wg sync.WaitGroup
	results := make(map[string]string)
	mu := &sync.Mutex{}
	sem := make(chan struct{}, 1000)

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量
			address := fmt.Sprintf("%s:%d", ip, p)
			conn, err := net.DialTimeout("tcp", address, 3*time.Second)
			mu.Lock()
			if err == nil {
				results[fmt.Sprintf("%s:%d", ip, p)] = "开启\n"
				conn.Close()
			}
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

// synScan 使用 SYN 扫描端口
func synScan(ip string) map[string]string {
	var wg sync.WaitGroup
	results := make(map[string]string)
	mu := &sync.Mutex{}
	sem := make(chan struct{}, 1000)

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			sendSYN(ip, p, 3*time.Second, results, mu)
		}(port)
	}

	wg.Wait()
	return results
}

// sendSYN 发送 SYN 数据包进行端口扫描
func sendSYN(ip string, port int, timeout time.Duration, results map[string]string, mu *sync.Mutex) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil { // 如果没有错误，表示端口开放
		mu.Lock()
		results[address] = "开启\n"
		mu.Unlock()
		defer conn.Close()
	}
}

// finScan 使用 FIN 扫描端口
func finScan(ip string) map[string]string {
	var wg sync.WaitGroup
	results := make(map[string]string)
	mu := &sync.Mutex{}
	sem := make(chan struct{}, 1000)

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			address := fmt.Sprintf("%s:%d", ip, p)
			sendFIN(address, 3*time.Second, results, mu)
		}(port)
	}

	wg.Wait()
	return results
}

// sendFIN 发送 FIN 数据包进行端口扫描
func sendFIN(address string, timeout time.Duration, results map[string]string, mu *sync.Mutex) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil {
		mu.Lock()
		results[address] = "开启\n"
		mu.Unlock()
		defer conn.Close()
	}
}
