package log4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : log4go_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/4 22:36
* 修改历史 : 1. [2022/4/4 22:36] 创建文件 by NST
*/

func TestConfiguration(t *testing.T) {
	xc := LoadConfiguration("./examples/example.xml")
	fmt.Println(xc)
}
