package actionmodules

import (
	"BackEnd/actionmodules/WordKiller/nethandlers"
	"fmt"
	"log"
	"time"
)

func WordKiller(username string, password string, mode string, week int, delay int) {
	println("\n[*]正在登录...")
	time.Sleep(2 * time.Second)
	token := nethandlers.Get_X_Auth_Token(username, password)
	fmt.Printf("[*]登录成功, token: %s\n", token)
	fmt.Println("[*]请输入测试周数:")
	fmt.Println("[*]请选择模式:自测[0]/考试[1]")
	fmt.Scanln(&mode)
	fmt.Println("[*]请输入做题时间(单位：s,推荐时长:450s):")
	switch mode {
	case "0":
		fmt.Println("[*]开始自测...")
	case "1":
		fmt.Println("[*]开始考试...")
	default:
		log.Fatalln("[-]模式错误")
	}
	nethandlers.Exam(token, week, mode, delay)
	//test.Test()
	fmt.Println("[*]运行结束.按任意键退出...")
}
