package main

import (
	"fmt"
	"time"
)

func main()  {
	//主线程阻塞
	timer1:=time.NewTimer(time.Second*5)
	<-timer1.C
	println("test")
	//主线程不阻塞
	timer2 := time.NewTimer(time.Second)
	go func() {
		//等触发时的信号
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()
	//由于上面的等待信号是在新线程中，所以代码会继续往下执行，停掉计时器
	time.Sleep(time.Second*5)
}

