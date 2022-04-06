package log4go

import (
	"fmt"
	"strconv"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : utils.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/5 21:17
* 修改历史 : 1. [2022/4/5 21:17] 创建文件 by NST
*/

func BuildFormatString(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func BuildString(a ...any) string {
	return fmt.Sprint(a...)
}

func LeftPad(str string, limit int, placeholder rune) string {
	len := limit - len(str)
	if len >= limit {
		return str
	}

	arr := make([]any, 0, len)
	for idx := 0; idx < len; idx++ {
		arr = append(arr, string(placeholder))
	}
	arr = append(arr, str)

	return BuildString(arr...)
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}
