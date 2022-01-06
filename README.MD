# websocket ＧＷ框架

## 使用方式:

```
	gwnet.SetConfig(&iface.Config{
		PingTime:       30,
		MaxConn:        100000,
		WorkerPoolSize: 10,
		MessageType:    2,
		WriteDeadlineDelay: 30,
	})
	
	s := gwnet.NewServer(gwnet.WithPacket(pack.NewDataPack()))
	
	//會根據封包的第一個int來分type進不同function處理
	s.AddRouter(0, &router.StartGame{})
	s.AddRouter(1 &router.EnterTable{})
	s.AddRouter(2, &router.LeaveTable{})
	
	//如果有自己的logger可以用此方法換掉zap的預設物件
	zap.ReplaceGlobals(logger)
	
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery(), middleware.Logs)
	g.GET("/ws", s.Start)
	g.Run(":8080")
	
	
```
