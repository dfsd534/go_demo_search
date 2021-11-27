package main

import (
	_ "github.com/dfsd534/go_demo_search/matchers"
	"github.com/dfsd534/go_demo_search/search"
	"log"
	"os"
)

func init() {
	//将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

func main() {
	//使用特定的项目做搜索
	search.Run("president")
}
