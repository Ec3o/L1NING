package actionmodules

import (
	"BackEnd/actionmodules/WordKiller/nethandlers"
	"BackEnd/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartKill(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}

	log.Println("WebSocket connection opened")
	go wordKillerWebSocket(conn, config.WordKillerUsername, config.WordKillerPassWord, config.WordKillerMode, config.WordKillerWeek, config.WordKillerDelay)
}

func wordKillerWebSocket(conn *websocket.Conn, username, password, mode string, week, delay int) {
	fmt.Println("\n[*]正在使用配置登录...")
	token := nethandlers.Get_X_Auth_Token(username, password)
	fmt.Printf("[*]登录成功, token: %s\n", token)
	fmt.Println("[*]请输入测试周数:")
	fmt.Println("[*]请选择模式:自测[0]/考试[1]")
	fmt.Println("[*]请输入做题时间(单位：s,推荐时长:450s):")
	switch mode {
	case "0":
		fmt.Println("[*]开始自测...")
	case "1":
		fmt.Println("[*]开始考试...")
	default:
		log.Fatalln("[-]模式错误")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from websocket:", err)
				// 添加详细错误日志
				log.Printf("Error details: %v", err)
				return
			}
			log.Printf("Received message: %s", msg)
			if string(msg) == "{\"type\":\"heartbeat\"}" {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"pong\"}")); err != nil {
					log.Println("Error sending pong:", err)
					return
				}
			}
		}
	}()

	for i := 0; i <= 100; i += 10 {
		select {
		case <-done:
			log.Println("Connection closed by client")
			return
		case <-ticker.C:
			if err := conn.WriteJSON(gin.H{"progress": i}); err != nil {
				log.Println("Error writing to websocket:", err)
				return
			}
			log.Printf("Progress sent: %d%%", i)
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error sending ping:", err)
				return
			}
			log.Println("Ping sent")
		}
	}

	sc := nethandlers.Exam(token, week, mode, delay)
	if err := conn.WriteJSON(gin.H{"progress": 100, "score": sc}); err != nil {
		log.Println("Error writing final result to websocket:", err)
		return
	}
	log.Println("Final result sent")
	log.Println("WebSocket connection closed")
	conn.Close()
}
