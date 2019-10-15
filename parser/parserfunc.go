package parser

import "io"

//ParserFunc 解析函数->解析器
type ParserFunc func(*Request, io.Reader) (*Result, error)

//Parse 解析函数
func (p ParserFunc) Parse(req *Request, r io.Reader) (*Result, error) {
	return p(req, r)
}
