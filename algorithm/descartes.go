package algorithm

/*
 example:

 dim := make([][]interface{}, 3)
    result := [][]interface{}{}
    for i := 0; i < 3; i++ {
        dim[i] = make([]interface{}, 3)
        for j := 0; j < 3; j++ {
            dim[i][j] = j + 1
        }
    }
    Descartes(dim, &result)
    log.Println(fmt.Sprintf("%+v", dim))
    log.Println(fmt.Sprintf("%+v", result))
*/

// 笛卡尔乘积算法
func descartes(dimValue [][]interface{}, result *[][]interface{}, layer int, cursor []interface{}) {
	size := len(dimValue[layer])
	// 递归二维数组
	if layer < len(dimValue)-1 {
		if size == 0 {
			descartes(dimValue, result, layer+1, cursor)
		} else {
			for i := 0; i < size; i++ {
				// 再接着与下一数组计算,所以创建新的游标
				newCursor := append(cursor, dimValue[layer][i])
				descartes(dimValue, result, layer+1, newCursor)
			}
		}
	} else {
		// 递归结束后添加到结果
		if size == 0 {
			if len(cursor) != 0 {
				*result = append(*result, cursor)
			}
		} else {
			for i := 0; i < size; i++ {
				r := append(cursor, dimValue[layer][i])
				*result = append(*result, r)
			}
		}
	}
}

// 笛卡尔乘积算法
func Descartes(dimValue [][]interface{}, result *[][]interface{}) {
	descartes(dimValue, result, 0, []interface{}{})
}

// 笛卡尔乘积算法
func DescartesInts(dim [][]int, result *[][]int) {
	dimNew := make([][]interface{}, len(dim))
	resultNew := make([][]interface{}, len(*result))
	for i, v := range dim {
		dimNew[i] = make([]interface{}, len(v))
		for j, k := range dim[i] {
			dimNew[i][j] = k
		}
	}
	Descartes(dimNew, &resultNew)

	*result = make([][]int, len(resultNew))
	for i, v := range resultNew {
		(*result)[i] = make([]int, len(v))
		for j, k := range resultNew[i] {
			(*result)[i][j] = k.(int)
		}
	}
}

// 笛卡尔乘积算法
func DescartesStrings(dim [][]string, result *[][]string) {
	dimNew := make([][]interface{}, len(dim))
	resultNew := make([][]interface{}, len(*result))
	for i, v := range dim {
		dimNew[i] = make([]interface{}, len(v))
		for j, k := range dim[i] {
			dimNew[i][j] = k
		}
	}
	Descartes(dimNew, &resultNew)

	*result = make([][]string, len(resultNew))
	for i, v := range resultNew {
		(*result)[i] = make([]string, len(v))
		for j, k := range resultNew[i] {
			(*result)[i][j] = k.(string)
		}
	}
}
