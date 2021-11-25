package search

import (
	"log"
	"sync"
)

//注册用于搜索的匹配器的映射
var matchers = make(map[string]Matcher)

//Run 执行搜索逻辑
func Run(searchTerm string) {
	feeds, err := RetrieveFeeds()

	if err != nil {
		log.Fatalln()
	}

	//创建一个无缓冲的通道，接收匹配后的结果
	results := make(chan *Result)

	//构造一个waitGroup,以便处理所有的数据源
	var waitGroup sync.WaitGroup

	//设置需要等待处理
	//每个数据源的goroutine的数据
	waitGroup.Add(len(feeds))

	//为每个数据源启动一个goroutine来查找结果
	for _, feed := range feeds {
		//获取一个匹配器用户查找
		matchers, exists := matchers[feed.Type]
		if !exists {
			matchers = matchers["default"]
		}

		//启动一个goroutine来执行搜索
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matchers, feed)
	}

	//启动一个goroutine来监控是否所有的工作都做完了
	go func() {
		//等待所有的任务完成
		waitGroup.Wait()

		//用关闭通道的方式，通知Display函数
		//可以退出程序了
		close(results)
	}()

	//启动函数，显示返回的结果，
	//并且在最后一个结果显示完成返回
	Display(results)
}

// Register Register调用时，会注册一个匹配器，提供给后面的程序使用
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
