package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/phantom-atom/crawler"
	"github.com/phantom-atom/crawler/parser"
)

//parserFunc 解析函数
func parserFunc(req *parser.Request, r io.Reader) (*parser.Result, error) {
	r = transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("div").Each(
		func(i int, s *goquery.Selection) {
			fmt.Println(s.Text())
		})

	nextPage := doc.Find("a").First()
	result := &parser.Result{}

	if href, exists := nextPage.Attr("href"); exists {
		//解析该url
		u, err := url.Parse(href)
		if err != nil {
			return nil, err
		}
		//如果查找到一个连接，那么添加这个连接到执行队列中
		//并将当前的执行器与等待组赋值给新的连接，以便执行相同的处理
		//以及等待执行的组
		result.AddRequest(&parser.Request{
			URL:       u,
			Parser:    req.Parser,
			WaitGroup: req.WaitGroup,
		})
	}
	return result, nil
}

func main() {
	engine := crawler.NewEngine(context.Background(), 1024, 1024)
	engine.Run()
	defer engine.Close()

	//获取解析的地址
	u, _ := url.Parse("http://example.com")

	wg := &sync.WaitGroup{}

	//添加一个处理
	engine.Process(&parser.Request{
		URL:       u,
		Parser:    parser.ParserFunc(parserFunc),
		WaitGroup: wg,
	})

	wg.Wait()
}
