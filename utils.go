package log4go

import "fmt"

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

func BuildString(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
