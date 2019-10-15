# crawler  <img src="https://github.com/phantom-atom/images/blob/master/561364cd63c69f0d1405dc659c9c139b6c7f1cd92d37a-mL4DlD_fw658.gif"  width = 256 align='right'/>

### crawler为一个简单的golang爬虫实现，只需要简单的实现解析接口即可，其中解析接口为：
    type Parser interface {
      //参数：
      //  request：为该次请求对应的Request接口提，
      //  r：为该次请求对应的数据
      //返回值
      //  *Result：执行后的返回值(并非数据返回值)
      //  error：错误信息(当不为nil时，不再继续执行*Result的请求)
      Parse(request *Request, r io.Reader) (*Result, error)
    }
### 此外，提供一个函数别名，用于直接使用函数作为解析器，该函数位于parser中
    type ParserFunc func(*Request, io.Reader) (*Result, error)
### 请求结构体
    //Request 请求结构体
    type Request struct {
      //请求的URL
      URL *url.URL

      //解析器，用于解析URL请求后的数据
      Parser Parser

      //等待组，用于等待解析完成，该变量
      //可以通过每次解析的时候，传递给下一个Request
      //这样只要在最开始调用的时候传入，可以等待所有请求结束
      WaitGroup *sync.WaitGroup
      
      //数据收集器，可用于数据接收使用
	  Collector collector.Collector
    }
### 例子
    //这里创建最多1024个goroutine同时处理请求，设置请求队列大小为2048
    engine := crawler.NewEngine(context.Background(), 1024, 2048)
    engine.Run()
    defer engine.Close()

    //将解析的地址转换为URL格式
    u, err := url.Parse("http://example.com")
    if (err != nil) {
      log.Fatalln(err)
    }
    //创建一个WaitGroup，用于等待该次处理结束
    wg := &sync.WaitGroup{}
    
    dataCh := make(chan int , 64)
    
    go func() {
       for v := range dataCh {
            log.Println(v)
       }
    }()
    
    //创建一个请求
    engine.Process(&parser.Request{
      URL:       u,
      Parser:    parser.ParserFunc(parserFunc),
      WaitGroup: wg,  //WaitGroup可以用于在解析后传入到下一个请求，即可以等待该系列结束
      Collector: typed.NewCollector(dataCh), //传入数据收集器
    })

    wg.Wait()
    close(dataCh)
