package actionmodules

import (
	"BackEnd/config"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
)

func GetData(ctx *gin.Context) {
	// 创建 Cookie Jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// 获取登录页面的 vrf 值和初始 sessionID
	vrf, sessionID := getKey(client)

	// 手动将初始 sessionID 添加到 Cookie Jar 中
	url, _ := url2.Parse("http://hdu.sunnysport.org.cn/")
	cookie := &http.Cookie{Name: "sessionid", Value: sessionID}
	client.Jar.SetCookies(url, []*http.Cookie{cookie})

	// 登录信息
	username := config.SunnySportUserName
	password := config.SunnySportPassWord

	// 构造登录请求的表单数据
	loginData := map[string]string{
		"username": username,
		"vrf":      vrf,
		"password": password,
	}

	// 构造登录请求
	loginURL := "http://hdu.sunnysport.org.cn/login/"
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(encodeFormData(loginData)))
	if err != nil {
		log.Fatalf("Error creating login request: %v", err)
	}

	// 设置请求头
	req.Header.Add("Host", "hdu.sunnysport.org.cn")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Origin", "http://hdu.sunnysport.org.cn")
	req.Header.Add("Referer", "http://hdu.sunnysport.org.cn/login/?next=/runner/index.html")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送登录请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending login request: %v", err)
	}
	defer resp.Body.Close()

	// 检查登录是否成功
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		log.Fatalf("Login failed with status code: %v", resp.StatusCode)
	}

	fmt.Println("Login successful!")

	// 获取阳光长跑数据的请求
	runnerURL := "http://hdu.sunnysport.org.cn/runner/index.html"
	runnerResp, err := client.Get(runnerURL)
	if err != nil {
		log.Fatalf("Error sending runner request: %v", err)
	}
	defer runnerResp.Body.Close()

	// 读取阳光长跑数据
	runnerBody, err := ioutil.ReadAll(runnerResp.Body)
	if err != nil {
		log.Fatalf("Error reading runner response: %v", err)
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(runnerBody))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// 提取长跑数据
	totalMileage := doc.Find("table.table-striped").First().Find("tr").Eq(0).Find("td").Eq(1).Text()
	validMileage := doc.Find("table.table-striped").First().Find("tr").Eq(1).Find("td").Eq(1).Text()
	averageSpeed := doc.Find("table.table-striped").First().Find("tr").Eq(2).Find("td").Eq(1).Text()
	validTimes := doc.Find("table.table-striped").First().Find("tr").Eq(3).Find("td").Eq(1).Text()

	// 创建返回的JSON数据
	data := map[string]string{
		"total_mileage": totalMileage,
		"valid_mileage": validMileage,
		"average_speed": averageSpeed,
		"valid_times":   validTimes,
	}

	// 返回JSON响应
	ctx.JSON(http.StatusOK, data)
}

// 获取登录页面的 vrf 值和初始 sessionID
func getKey(client *http.Client) (string, string) {
	resp, err := client.Get("http://hdu.sunnysport.org.cn")
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL: %s", resp.Status)
	}

	// 解析 HTML 获取 vrf 值
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var vrf string
	var sessionID string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "vrf" {
					for _, attr := range n.Attr {
						if attr.Key == "value" {
							vrf = attr.Val
							break
						}
					}
				}
			}
		}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "sessionid" {
				sessionID = cookie.Value
				break
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return vrf, sessionID
}

// 将表单数据编码为 URL 编码形式
func encodeFormData(data map[string]string) string {
	var encodedData []string
	for key, value := range data {
		encodedData = append(encodedData, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(encodedData, "&")
}
