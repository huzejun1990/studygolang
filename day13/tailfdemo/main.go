package main

/*import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	fileName := "./xx.log"
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	//打开文件
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail %s faild,err:%v\n", fileName, err)
	}
	//开始读取数据
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen,filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
		}
		fmt.Println("msg", msg.Text)
	}
}*/

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

func main() {

	fileName := `./xx.log`
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail %s failed,err:%v\n", fileName, err)
		return
	}
	//开始读取数据
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines // chan tail.Line
		if !ok {
			fmt.Printf("tail file close reopen,filename:%s\n",
				tails.Filename) //
			time.Sleep(time.Second) //读取出错，等1秒
			continue
		}
		fmt.Println("msg", msg.Text)
	}
}
