package main

import (
	"fmt"
	"math"
)

// Cluster 结构定义
type Cluster struct {
	rights  map[string]bool
	members map[string]bool
}

// 计算两个类簇之间的距离
func calculateDistance(c1, c2 *Cluster) float64 {
	// 计算交集大小
	intersectionSize := 0
	for user := range c1.members {
		if _, exists := c2.members[user]; exists {
			intersectionSize++
		}
	}

	// 计算并集大小
	unionSize := len(c1.rights) + len(c2.rights)
	for permission := range c1.rights {
		if _, exists := c2.rights[permission]; exists {
			unionSize--
		}
	}

	// 距离公式：交集 / 权限并集
	if unionSize == 0 {
		return 0 // 防止除以零
	}
	return float64(intersectionSize) / float64(unionSize)
}

// 合并两个类簇
func mergeClusters(c1, c2 *Cluster) *Cluster {
	// 合并权限集（并集）
	mergedRights := make(map[string]bool)
	for permission := range c1.rights {
		mergedRights[permission] = true
	}
	for permission := range c2.rights {
		mergedRights[permission] = true
	}

	// 如果并集为空，保持原来的权限集
	if len(mergedRights) == 0 {
		return c1
	}

	// 合并用户集（交集）
	mergedMembers := make(map[string]bool)
	for user := range c1.members {
		if _, exists := c2.members[user]; exists {
			mergedMembers[user] = true
		}
	}

	// 如果交集为空，保持原来的用户集
	if len(mergedMembers) == 0 {
		return c1
	}

	// 返回合并后的Cluster
	return &Cluster{
		rights:  mergedRights,
		members: mergedMembers,
	}
}

// 判断类簇 c1 是否是 c2 的子集
func isSubset(c1, c2 *Cluster) bool {
	// 判断权限是否是子集
	for permission := range c1.rights {
		if !c2.rights[permission] {
			return false
		}
	}
	return true
}

// 打印类簇信息
func printCluster(c *Cluster) {
	fmt.Printf("权限：%v, 用户：%v\n", c.rights, c.members)
}

func printOrderRelations(orderRelations []string) {
	fmt.Println("偏序关系集合 <：")
	for _, relation := range orderRelations {
		fmt.Println(relation)
	}
}
func printClusters(clusters []*Cluster) {
	for i, cluster := range clusters {
		fmt.Printf("C%d: ", i)
		printCluster(cluster)
	}
}
func isEqualCluster(c1, c2 *Cluster) bool {
	if len(c1.rights) != len(c2.rights) || len(c1.members) != len(c2.members) {
		return false
	}

	for permission := range c1.rights {
		if !c2.rights[permission] {
			return false
		}
	}

	for user := range c1.members {
		if !c2.members[user] {
			return false
		}
	}

	return true
}

func a(clusters []*Cluster) {
	// 偏序关系集合 <，用于存储偏序关系的类簇对
	var orderRelations []string

	// 记录无偏序关系的类簇集合 T_
	var T_ []*Cluster

	// 初始聚类
	for len(clusters) > 1 {
		T_ = nil
		// 计算最短距离的两个类簇
		minDist := math.MaxFloat64
		var c1, c2 *Cluster
		c1Index, c2Index := -1, -1
		for i := 0; i < len(clusters); i++ {
			for j := i + 1; j < len(clusters); j++ {
				dist := calculateDistance(clusters[i], clusters[j])
				if dist < minDist {
					minDist = dist
					c1, c2 = clusters[i], clusters[j]
					c1Index, c2Index = i, j
				}
			}
		}

		// 合并类簇
		mergedCluster := mergeClusters(c1, c2)
		if isEqualCluster(mergedCluster, c1) {
			break
		}

		// 将合并后的新簇追加到簇集合中
		clusters = append(clusters, mergedCluster)

		// 计算偏序关系
		var newOrderRelations []string
		for i := 0; i < len(clusters); i++ {
			for j := i + 1; j < len(clusters); j++ {
				if isSubset(clusters[i], clusters[j]) {
					newOrderRelations = append(newOrderRelations, fmt.Sprintf("%v < %v", clusters[i].rights, clusters[j].rights))
				} else {
					T_ = append(T_, clusters[i], clusters[j])
				}
			}
		}
		// 删除原来的两个簇
		if c1Index > c2Index {
			clusters = append(clusters[:c1Index], clusters[c1Index+1:]...)
			clusters = append(clusters[:c2Index], clusters[c2Index+1:]...)
		} else {
			clusters = append(clusters[:c2Index], clusters[c2Index+1:]...)
			clusters = append(clusters[:c1Index], clusters[c1Index+1:]...)
		}

		orderRelations = append(orderRelations, newOrderRelations...)
		if T_ == nil {
			break
		}

	}

	// 最终输出结果
	fmt.Println("最终聚类结果：")
	printClusters(clusters)
	printOrderRelations(orderRelations)
	fmt.Println("无偏序关系类簇集合 T_：")
	for _, cluster := range T_ {
		printCluster(cluster)
	}
}

func main() {
	// 初始化原始类簇（根据测试数据初始化）
	clusters := []*Cluster{
		{
			rights:  map[string]bool{"read": true},
			members: map[string]bool{"user1": true, "user2": true},
		},
		{
			rights:  map[string]bool{"write": true},
			members: map[string]bool{"user1": true, "user3": true},
		},
		{
			rights:  map[string]bool{"execute": true},
			members: map[string]bool{"user2": true, "user3": true},
		},
	}
	clusters1 := []*Cluster{
		{
			rights:  map[string]bool{"权限1": true},
			members: map[string]bool{"userA": true, "userB": true, "userC": true, "userD": true},
		},
		{
			rights:  map[string]bool{"权限2": true},
			members: map[string]bool{"userA": true, "userB": true, "userD": true},
		},
		{
			rights:  map[string]bool{"权限3": true},
			members: map[string]bool{"userC": true, "userD": true},
		},
		{
			rights:  map[string]bool{"权限4": true},
			members: map[string]bool{"userA": true, "userD": true},
		},
		{
			rights:  map[string]bool{"权限5": true},
			members: map[string]bool{"userC": true, "userD": true},
		},
	}
	fmt.Println("题目一:")
	a(clusters)
	fmt.Println("题目二:")
	a(clusters1)

}
