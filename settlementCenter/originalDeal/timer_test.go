package originalDeal

import (
	"fmt"
	"os"
	"testing"
)

func TestTimer1(t *testing.T) {
	Timer1()
}

func TestTimer(t *testing.T) {
	crontab := NewCrontab()
	// 实现接口的方式添加定时任务
	task := &testTask{}
	if err := crontab.AddByID("1", "* * * * *", task); err != nil {
		fmt.Printf("error to add crontab task:%s", err)
		os.Exit(-1)
	}

	// 添加函数作为定时任务
	taskFunc := func() {
		fmt.Println("hello world")
	}
	if err := crontab.AddByFunc("2", "* * * * *", taskFunc); err != nil {
		fmt.Printf("error to add crontab task:%s", err)
		os.Exit(-1)
	}
	crontab.Start()
	select {}
}

type testTask struct {
}

func (t *testTask) Run() {
	fmt.Println("hello world")
}

