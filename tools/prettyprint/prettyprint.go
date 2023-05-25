package prettyprint

func PrettyPrint(titles []string, data [][]string) {
	// 计算每列的最大长度
	maxLengths := make([]int, len(titles))
	for i := 0; i < len(titles); i++ {
		maxLengths[i] = len(titles[i])
	}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if maxLengths[j] < len(data[i][j]) {
				maxLengths[j] = len(data[i][j])
			}
		}
	}

	// 打印表头
	PrintLine(maxLengths)
	print("|")
	for i := 0; i < len(titles); i++ {
		// 打印每列的标题
		PrintCell(titles[i], maxLengths[i])
	}
	println()
	PrintLine(maxLengths)
	// 打印数据
	for i := 0; i < len(data); i++ {
		print("|")
		for j := 0; j < len(data[i]); j++ {
			PrintCell(data[i][j], maxLengths[j])
		}
		println()
		 
	}
	PrintLine(maxLengths)
}

func PrintCell(data string, length int) {
	// 打印数据
	print(data)
	// 补齐空格
	for i := 0; i < length-len(data); i++ {
		print(" ")
	}
	print("|")
}

func PrintLine(lengths []int) {
	// println()
	for i := 0; i < len(lengths); i++ {
		for j := 0; j < lengths[i]+1; j++ {
			print("-")
		}
	}
	println()
}
