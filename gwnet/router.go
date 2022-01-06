package gwnet

import "github.com/kevinlincg/gw_websocket/gwiface"

// BaseRouter 要實現router時，先內崁這個基礎類別，然後根據需要對這個基本類別的function重寫
type BaseRouter struct{}

// 這裡之所以以BaseRouter的實作都是空的，
// 是因為有可能有的Router不會有PreHandle或PostHandle
// 所以Router全部繼承BaseRouter的好處是，不用重新實作一次PreHandle和PostHandle也可以

// PreHandle -
func (br *BaseRouter) PreHandle(req gwiface.Request) {}

// Handle -
func (br *BaseRouter) Handle(req gwiface.Request) {}

// PostHandle -
func (br *BaseRouter) PostHandle(req gwiface.Request) {}
