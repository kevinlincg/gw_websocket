package gwiface

/*
	Router主邏輯，一個邏輯定義一個struct來處理
*/
type Router interface {
	PreHandle(request Request)
	Handle(request Request)
	PostHandle(request Request)
}
