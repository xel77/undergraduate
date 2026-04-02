package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
	"math/rand"
	"time"
)

// powerCompute 计算功率增量
func powerCompute(bit int, BER, Npsd float64) float64 {
	bitNum := math.Pow(2, float64(bit)) - 1
	Gamma := math.Sqrt(2) * erfcInv(2*BER/4)
	powerAdd := Npsd / 3 * (Gamma * Gamma) * bitNum
	return powerAdd
}

// erfcInv 反误差函数的近似实现
func erfcInv(x float64) float64 {
	return math.Sqrt(-2 * math.Log(x))
}

// HughesHartogs 主算法
func HughesHartogs(NSubcarriers, Rt, M int, BER, Npsd float64, H []float64) ([]int, []float64) {
	bitAlloc := make([]int, NSubcarriers)
	powerAlloc := make([]float64, NSubcarriers)
	powerAdd := make([]float64, NSubcarriers)

	// 初始化 powerAdd 数组
	for i := 0; i < NSubcarriers; i++ {
		powerAdd[i] = (powerCompute(bitAlloc[i]+1, BER, Npsd) - powerCompute(bitAlloc[i], BER, Npsd)) / math.Pow(H[i], 2)
	}

	// 找到最小功率增量的索引
	_, indexMin := findMinAndIndex(powerAdd)
	bitAlloc[indexMin]++

	// 计算总比特数
	bitTotal := 0
	for _, bits := range bitAlloc {
		bitTotal += bits
	}

	// 循环直到比特分配达到目标比特数 Rt
	for bitTotal < Rt {
		if bitAlloc[indexMin] < M {
			// 更新 powerAdd 数组
			for i := 0; i < NSubcarriers; i++ {
				powerAdd[i] = (powerCompute(bitAlloc[i]+1, BER, Npsd) - powerCompute(bitAlloc[i], BER, Npsd)) / math.Pow(H[i], 2)
			}

			// 找到最小功率增量及其索引
			_, indexMin = findMinAndIndex(powerAdd)

			// 更新比特分配
			bitAlloc[indexMin]++

			// 更新比特总数
			bitTotal = 0
			for _, bits := range bitAlloc {
				bitTotal += bits
			}
		} else {
			// 如果已经达到最大调制方式，设定功率增量为无穷大
			powerAdd[indexMin] = math.Inf(1)

			// 重新计算最小功率增量和索引
			_, indexMin = findMinAndIndex(powerAdd)

			// 更新比特分配
			bitAlloc[indexMin]++

			// 更新比特总数
			bitTotal = 0
			for _, bits := range bitAlloc {
				bitTotal += bits
			}
		}
	}

	// 计算功率分配
	for i := 0; i < NSubcarriers; i++ {
		powerAlloc[i] = powerCompute(bitAlloc[i], BER, Npsd) / math.Pow(H[i], 2)
	}

	return bitAlloc, powerAlloc
}

// findMinAndIndex 返回切片中的最小值及其索引
func findMinAndIndex(slice []float64) (float64, int) {
	minValue := slice[0]
	minIndex := 0
	for i, value := range slice[1:] {
		if value < minValue {
			minValue = value
			minIndex = i + 1
		}
	}
	return minValue, minIndex
}

// generateRayleighDistribution 生成 Rayleigh 分布
func generateRayleighDistribution(NSubcarriers int) []float64 {
	rand.Seed(time.Now().UnixNano())
	H := make([]float64, NSubcarriers)
	scale := 1.0 // Rayleigh分布的尺度参数

	// 通过生成两个标准正态分布的随机变量 X 和 Y 来生成 Rayleigh 分布
	for i := 0; i < NSubcarriers; i++ {
		X := rand.NormFloat64()           // 标准正态分布
		Y := rand.NormFloat64()           // 标准正态分布
		H[i] = scale * math.Sqrt(X*X+Y*Y) // Rayleigh分布
	}

	return H
}

// Chow 算法的实现
func Chow(SNR []float64, H []float64, N_subc int, Rt int, gap float64) ([]int, []float64, int) {
	N_use := 0
	margin := 0.0
	Iterate_count := 0
	total_bits := 0
	bits_alloc := make([]int, N_subc)
	power_alloc := make([]float64, N_subc)
	temp_bits := make([]float64, N_subc)
	bit_round := make([]int, N_subc)
	diff := make([]float64, N_subc)

	// 循环直到达到目标比特数
	for total_bits < Rt && Iterate_count < 1 {
		for i := 0; i < N_subc; i++ {
			temp_bits[i] = math.Log2(1+math.Pow(SNR[i], 2)) / (1 + margin/gap)
			bit_round[i] = int(math.Round(temp_bits[i]))
			diff[i] = temp_bits[i] - float64(bit_round[i])
		}

		// 计算总比特数
		total_bits = 0
		for _, b := range bit_round {
			total_bits += b
		}

		// 如果信道不可用，继续循环
		if total_bits == 0 {
			fmt.Println("The channel is not usable.")
			continue
		}

		// 计算未使用的子载波数目
		N_notuse := 0
		for _, b := range bit_round {
			if b == 0 {
				N_notuse++
			}
		}
		N_use = N_subc - N_notuse

		// 计算信噪比裕量
		margin += 10 * math.Log10(math.Pow(2, float64(total_bits-Rt)/float64(N_use)))
		Iterate_count++
	}

	// 调整比特数以满足目标比特数
	for total_bits > Rt {
		_, indexMin := findMinAndIndex(diff)
		bit_round[indexMin]--
		diff[indexMin]++
		total_bits = 0
		for _, b := range bit_round {
			total_bits += b
		}
	}

	// 调整比特数以达到目标比特数
	for total_bits < Rt {
		maxDiffIndex := -1
		maxDiffValue := -math.MaxFloat64
		for i, d := range diff {
			if bit_round[i] != 0 && d > maxDiffValue {
				maxDiffValue = d
				maxDiffIndex = i
			}
		}
		bit_round[maxDiffIndex]++
		diff[maxDiffIndex]--
		total_bits = 0
		for _, b := range bit_round {
			total_bits += b
		}
	}

	// 计算功率分配
	for i := 0; i < N_subc; i++ {
		bits_alloc[i] = bit_round[i]
		power_alloc[i] = (math.Pow(2, float64(bits_alloc[i])) - 1) / (math.Pow(H[i], 2) * gap)
	}

	return bits_alloc, power_alloc, Iterate_count
}

// 调用 Chow 算法
func chow(N_subc int, P_av float64, gap float64, Rt int, SNR_av float64) {
	// 初始化子载波的信道增益（Rayleigh分布）
	H := generateRayleighDistribution(N_subc)

	// 计算每个子载波的信噪比
	SNR := make([]float64, N_subc)
	for i := 0; i < N_subc; i++ {
		SNR[i] = math.Pow(H[i], 2) / ((P_av * gap) / math.Pow(10, SNR_av/10))
	}

	// 调用 Chow 算法
	bits_alloc, power_alloc, Iterate_count := Chow(SNR, H, N_subc, Rt, gap)

	// 输出结果
	fmt.Println("迭代次数:", Iterate_count)
	fmt.Println("Bit分配:", bits_alloc)
	fmt.Println("功率分配:", power_alloc)

	// 调整功率值，准备绘图
	pow_a := make([]float64, N_subc)
	for i, b := range power_alloc {
		pow_a[i] = b / 6000
	}

	// 生成图像
	makepng(N_subc, bits_alloc, pow_a, "chow算法")
}

// 调用 Hughes-Hartogs 算法
func h2h(NSubcarriers int, Rt int, M int, BER float64, Npsd float64) {
	// 生成 Rayleigh 信道增益
	H := generateRayleighDistribution(NSubcarriers)

	// 调用 Hughes-Hartogs 算法
	bitAlloc, powerAlloc := HughesHartogs(NSubcarriers, Rt, M, BER, Npsd, H)

	// 输出比特分配和功率分配
	fmt.Println("Bit分配:")
	fmt.Println(bitAlloc)
	fmt.Println("功率分配:")
	fmt.Println(powerAlloc)

	// 生成图像
	makepng(NSubcarriers, bitAlloc, powerAlloc, "h-h算法")
}

// 创建图表并保存为 PNG 文件
func makepng(NSubcarriers int, bitAlloc []int, powerAlloc []float64, name string) {
	// 创建一个新的图表
	p := plot.New()

	// 创建比特分配的条形图
	barData := make(plotter.Values, NSubcarriers)
	for i := 0; i < NSubcarriers; i++ {
		barData[i] = float64(bitAlloc[i])
	}
	bars, err := plotter.NewBarChart(barData, vg.Points(20))
	if err != nil {
		return
	}
	p.Add(bars)

	// 创建功率分配的折线图，并将功率值扩大 100 倍
	lineData := make(plotter.XYs, NSubcarriers)
	for i := 0; i < NSubcarriers; i++ {
		lineData[i].X = float64(i)
		lineData[i].Y = powerAlloc[i] * 50 // 扩大 100 倍
	}

	// 创建红色的折线图
	line, err := plotter.NewLine(lineData)
	if err != nil {
		return
	}
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	p.Add(line)

	// 设置标题和标签
	p.X.Label.Text = "子载波序列"
	p.Y.Label.Text = "bit分配"

	// 计算纵轴最大值，并设置为 bitAlloc 中的最大值
	maxBitAlloc := float64(0)
	for _, bit := range bitAlloc {
		if float64(bit) > maxBitAlloc {
			maxBitAlloc = float64(bit)
		}
	}

	// 将 Y 轴的最大值设置为最大值的三倍，纵轴范围分成 3 个区间
	p.Y.Min = 0
	p.Y.Max = maxBitAlloc * 1.1

	// 设置 Y 轴的刻度间隔
	interval := maxBitAlloc / 3
	ticks := []plot.Tick{}
	for i := 0; i <= 3; i++ {
		ticks = append(ticks, plot.Tick{Value: float64(i) * interval, Label: fmt.Sprintf("%.0f", float64(i)*interval)})
	}

	// 保存图像为 PNG 文件
	if err := p.Save(8*vg.Inch, 6*vg.Inch, name+".png"); err != nil {
		fmt.Println(err)
	}
}

// main 函数
func main() {
	// 参数初始化
	NSubcarriers := 15                      // 子载波数目
	Pav := 2.0                              // 平均功率
	SNRav := 10.0                           // 信噪比
	Noise := Pav / math.Pow(10, (SNRav/10)) // 噪声功率
	B := 1e6
	Npsd := Noise / (B / float64(NSubcarriers)) // 噪声功率谱密度
	BER := 1e-4                                 // 误比特率
	M := 10                                     // 最大调制方法
	Rt := 128                                   // 目标比特数
	gap := -math.Log(5*BER) / 1.5               // 信噪比裕量计算

	var choice int
	fmt.Println("请选择算法：")
	fmt.Println("1. Hughes-Hartogs 算法")
	fmt.Println("2. Chow 算法")
	fmt.Print("请输入 1 或 2: ")
	_, err := fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("输入错误，请重新输入")
		return
	}

	// 根据用户输入选择算法
	switch choice {
	case 1:
		// 调用 Hughes-Hartogs 算法
		h2h(NSubcarriers, Rt, M, BER, Npsd)
	case 2:
		// 调用 Chow 算法
		chow(NSubcarriers, Pav, gap, Rt, SNRav)
	default:
		fmt.Println("无效的选择，请输入 1 或 2")
	}
}
