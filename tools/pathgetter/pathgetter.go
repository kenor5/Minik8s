package pathgetter

import (
	"os"
	"path/filepath"

)


// 获取当前运行函数的绝对路径
// 如果出现错误，可能是go run 和 go build 的路径不同导致的
func GetCurrentAbPath() string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
	return exPath
}
