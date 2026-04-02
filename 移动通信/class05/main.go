package main

import (
	"fmt"
	"math"
)

// User 表示系统中的用户
type User struct {
	ID              int     // 用户ID
	SNR             float64 // 信噪比
	QueueDelay      float64 // 队列延迟
	AverageRate     float64 // 平均传输速率
	InstantRate     float64 // 瞬时传输速率
	Throughput      float64 // 吞吐量
	LastServiceTime int     // 上次服务时间
	RemainingTime   int     // 剩余需要服务的时间
}

// 最大载干比算法
func maxSNRScheduling(users []User, totalTime int) {
	if len(users) == 0 {
		return
	}

	fmt.Println("\n最大载干比算法运行过程：")
	currentTime := 0

	// 深拷贝users数组
	workingUsers := make([]User, len(users))
	copy(workingUsers, users)

	for currentTime < totalTime {
		fmt.Printf("\n当前时间: %d\n", currentTime)
		// 选择SNR最大的用户
		maxSNR := -math.MaxFloat64
		selectedUser := -1

		for i, user := range workingUsers {
			if user.RemainingTime > 0 && user.SNR > maxSNR {
				maxSNR = user.SNR
				selectedUser = i
			}
		}

		if selectedUser == -1 {
			break // 所有用户都完成了
		}

		// 为选中的用户分配时间片
		executionTime := 1 // 每次分配1个时间单位
		workingUsers[selectedUser].RemainingTime -= executionTime
		currentTime += executionTime

		fmt.Printf("选择用户%d (SNR=%.2f)执行1个时间单位，剩余时间：%d\n",
			workingUsers[selectedUser].ID, maxSNR, workingUsers[selectedUser].RemainingTime)
	}
}

// 轮询调度算法
func roundRobinScheduling(users []User, timeSlice int) {
	if len(users) == 0 {
		return
	}

	fmt.Println("\n轮询调度算法运行过程：")
	currentTime := 0
	remainingUsers := len(users)

	// 深拷贝users数组以保留原始数据
	workingUsers := make([]User, len(users))
	copy(workingUsers, users)

	for remainingUsers > 0 {
		fmt.Printf("\n时间片开始时间: %d\n", currentTime)
		for i := range workingUsers {
			if workingUsers[i].RemainingTime > 0 {
				executionTime := math.Min(float64(timeSlice), float64(workingUsers[i].RemainingTime))
				workingUsers[i].RemainingTime -= int(executionTime)
				currentTime += int(executionTime)

				fmt.Printf("用户%d执行了%.0f个时间单位，剩余时间：%d\n",
					workingUsers[i].ID, executionTime, workingUsers[i].RemainingTime)

				if workingUsers[i].RemainingTime == 0 {
					remainingUsers--
					fmt.Printf("用户%d在时间%d完成了所有任务\n", workingUsers[i].ID, currentTime)
				}
			}
		}
	}
}

// 比例公平算法
func proportionalFairScheduling(users []User, totalTime int) {
	if len(users) == 0 {
		return
	}

	fmt.Println("\n比例公平算法运行过程：")
	currentTime := 0

	// 深拷贝users数组
	workingUsers := make([]User, len(users))
	copy(workingUsers, users)

	for currentTime < totalTime {
		fmt.Printf("\n当前时间: %d\n", currentTime)
		maxPriority := -math.MaxFloat64
		selectedUser := -1

		for i, user := range workingUsers {
			if user.RemainingTime > 0 {
				priority := user.InstantRate / user.AverageRate
				if priority > maxPriority {
					maxPriority = priority
					selectedUser = i
				}
			}
		}

		if selectedUser == -1 {
			break
		}

		executionTime := 1
		workingUsers[selectedUser].RemainingTime -= executionTime
		currentTime += executionTime

		fmt.Printf("选择用户%d (优先级=%.2f)执行1个时间单位，剩余时间：%d\n",
			workingUsers[selectedUser].ID, maxPriority, workingUsers[selectedUser].RemainingTime)
	}
}

// 最大加权时延优先算法
func mLWDFScheduling(users []User) int {
	if len(users) == 0 {
		return -1
	}

	fmt.Println("\n最大加权时延优先算法运行过程：")
	fmt.Println("当前各用户优先级计算：")

	maxPriority := -math.MaxFloat64
	selectedUser := -1

	for i, user := range users {
		priority := user.QueueDelay * (user.InstantRate / user.AverageRate)
		fmt.Printf("用户%d: 队列延迟=%.2f, 速率比=%.2f, 优先级=%.2f\n",
			user.ID, user.QueueDelay, user.InstantRate/user.AverageRate, priority)

		if priority > maxPriority {
			maxPriority = priority
			selectedUser = i
		}
	}
	fmt.Printf("选择最高优先级%.2f的用户%d\n", maxPriority, selectedUser)
	return selectedUser
}

// 公平吞吐量算法
func fairThroughputScheduling(users []User) int {
	if len(users) == 0 {
		return -1
	}

	fmt.Println("\n公平吞吐量算法运行过程：")
	fmt.Println("当前各用户吞吐量：")
	for _, user := range users {
		fmt.Printf("用户%d: 吞吐量=%.2f\n", user.ID, user.Throughput)
	}

	minThroughput := math.MaxFloat64
	selectedUser := -1

	for i, user := range users {
		if user.Throughput < minThroughput {
			minThroughput = user.Throughput
			selectedUser = i
		}
	}
	fmt.Printf("选择最小吞吐量%.2f的用户%d以保证公平性\n", minThroughput, selectedUser)
	return selectedUser
}

func main() {
	users := []User{
		{ID: 0, SNR: 20.5, QueueDelay: 1.2, AverageRate: 5.0, InstantRate: 7.0, Throughput: 6.0, LastServiceTime: 0, RemainingTime: 5},
		{ID: 1, SNR: 18.3, QueueDelay: 2.1, AverageRate: 4.0, InstantRate: 6.0, Throughput: 4.5, LastServiceTime: 1, RemainingTime: 7},
		{ID: 2, SNR: 22.1, QueueDelay: 0.8, AverageRate: 6.0, InstantRate: 8.0, Throughput: 7.0, LastServiceTime: 2, RemainingTime: 13},
		{ID: 3, SNR: 19.7, QueueDelay: 1.5, AverageRate: 5.5, InstantRate: 7.5, Throughput: 5.5, LastServiceTime: 3, RemainingTime: 15},
	}

	totalTime := 20 // 设置总时间片为20

	// 测试各种调度算法
	maxSNRScheduling(users, totalTime)
	roundRobinScheduling(users, 5)
	proportionalFairScheduling(users, totalTime)
	mLWDFScheduling(users)
	fairThroughputScheduling(users)
}
