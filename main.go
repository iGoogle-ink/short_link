package main

func main() {
	// 短地址服务
	a := new(App)
	a.Initialize()
	a.Run(":8000")
}
