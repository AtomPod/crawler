package parser

//Result 解析器返回值
type Result struct {
	requests []*Request //请求列表
}

//Requests 获取请求列表
func (r *Result) Requests() []*Request {
	return r.requests
}

//AddRequest 添加一个请求
func (r *Result) AddRequest(req *Request) {
	r.requests = append(r.requests, req)
}
