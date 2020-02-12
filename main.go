package main

func main() {
	// 短地址服务
	a := new(App)
	a.Initialize(getEnv())
	a.Run(":8000")
}

/*
	Simple Demo
*/
