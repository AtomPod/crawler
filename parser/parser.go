package parser

import "io"

//Parser 解析器
type Parser interface {
	//用于解析请求的数据，request为当前请求的Request，
	//r为URL请求后获取的数据，如果返回值中error不为nil
	//则返回错误，不继续执行，接着判断是否*Result不为nil，如果
	//*Result不为nil，则继续执行请求解析
	Parse(request *Request, r io.Reader) (*Result, error)
}
