package nethandlers

import (
	"BackEnd/actionmodules/WordKiller/models"
	"BackEnd/actionmodules/WordKiller/ui"
	"BackEnd/actionmodules/WordKiller/wordhandlers"
	"bytes"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var keyRegexp = regexp.MustCompile("<p id=\"login-croypto\">(.*?)</p>")
var executionRegexp = regexp.MustCompile("<p id=\"login-page-flowkey\">(.*?)</p>")

func GetHeaders(token string) http.Header {
	ticket := generateTicket(21) // 自定义函数，模拟 JavaScript 中的 ticket 函数
	headers := http.Header{}
	headers.Set("Skl-Ticket", ticket)
	headers.Set("X-Auth-Token", token)
	headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 4.2.1; M040 Build/JOP40D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Connection", "keep-alive")
	headers.Set("Referer", "https://skl.hdu.edu.cn/")
	return headers
}
func PostHeader(token string) http.Header {
	headers := http.Header{}
	headers.Set("Host", "skl.hdu.edu.cn")
	headers.Set("Sec-Ch-Ua", `"Chromium";v="117", "Not;A=Brand";v="8"`)
	headers.Set("Skl-Ticket", generateTicket(21))
	headers.Set("Sec-Ch-Ua-Mobile", "?0")
	headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 4.2.1; M040 Build/JOP40D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36")
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("X-Auth-Token", token)
	headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	headers.Set("Origin", "https://skl.hduhelp.com")
	headers.Set("Sec-Fetch-Site", "cross-site")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Referer", "https://skl.hduhelp.com/")
	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	return headers
}
func generateTicket(length int) string {
	const NL = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	ticket := make([]byte, length)
	for i, b := range bytes {
		ticket[i] = NL[b&63] // 与 `& 63` 相同，确保索引在 0-63 范围内
	}
	return string(ticket)
}
func Exam(token string, week int, mode string, delay int) string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error retrieving working directory:", err)
		return "服务器内部错误"
	}
	client := &http.Client{}
	startTime := time.Now().UnixMilli()
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/paper/new?type=%s&week=%d&startTime=%d", mode, week, startTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = GetHeaders(token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//将body保存到json文件中
	filename := fmt.Sprintf("paper_%s.json", time.Now().Format("20060102150405"))
	fullPath := filepath.Join(wd, "actionmodules", "WordKiller", "papers", filename)
	err = ioutil.WriteFile(fullPath, body, 0644)
	fmt.Println("[*]存储试题信息中...")
	if err != nil {
		fmt.Println("[x]文件存储错误:", err)
		return "Error:文件存储错误"
	}

	var paper models.Paper
	if err := json.Unmarshal(body, &paper); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[*]等待提交中...")
	go ui.DisplayProgress(time.Duration(delay) * time.Second)
	time.Sleep(time.Duration(delay) * time.Second)

	ans := wordhandlers.GetAnswer(paper)
	data, err := json.Marshal(ans)
	//fmt.Println(string(data))
	if err != nil {
		log.Fatal(err)
	}

	saveURL := "https://skl.hdu.edu.cn/api/paper/save"
	req, err = http.NewRequest("POST", saveURL, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header = PostHeader(token)
	resp, err = client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	//fmt.Println("[*]Response status:", resp.Status)
	//fmt.Println("[*]Response body:", string(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("[+]提交成功！")
	time.Sleep(2 * time.Second)
	fmt.Println("[+]系统处理中...")
	time.Sleep(2 * time.Second)
	fmt.Println("[+]正在查看成绩...")
	detailURL := fmt.Sprintf("https://skl.hdu.edu.cn/api/paper/detail?paperId=%s", paper.PaperId)
	req, err = http.NewRequest("GET", detailURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = GetHeaders(token)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var score models.ExamScore
	if err := json.Unmarshal(body, &score); err != nil {
		log.Fatal(err)
	}
	//fmt.Println("[*]Response status:", resp.Status)
	//fmt.Println("[*]Response body:", string(body))
	fmt.Printf("[+]本次成绩: %d\n", score.Mark)
	return strconv.Itoa(score.Mark)

}
func Get_X_Auth_Token(username string, password string) string {
	var key, execution []byte
	casLoginURL := "https://skl.hdu.edu.cn/api/userinfo?type=&index="
	ssologinurl := "https://sso.hdu.edu.cn/login"
	//Get URL0
	client := &http.Client{ //禁用自动跳转，手动处理Redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, casLoginURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}
	url0 := gjson.Get(string(body), "url").String()
	//println("URL0:" + url0)

	//Get URL1
	req, err = http.NewRequest(http.MethodGet, url0, nil) //第二步，请求url，再拿到header里面的location URL进行一次请求
	addHeaders(req)
	resp, err = client.Do(req)
	url1 := resp.Header.Get("Location")
	//fmt.Println("URL1:", url1)
	//拿到和解析HTML
	req, err = http.NewRequest(http.MethodGet, url1, nil)
	addHeaderssso(req)
	resp, err = client.Do(req)
	body, err = io.ReadAll(resp.Body)
	cookie := resp.Cookies()
	session := cookie[0].Value
	//fmt.Println("Cookie: SESSION=" + session)
	//fmt.Println(string(body))
	tmp := keyRegexp.FindSubmatch(body)
	key = tmp[1]
	tmp = executionRegexp.FindSubmatch(body)
	execution = tmp[1]
	bytes.Trim(key, " \"\r\n")

	//获取password
	encryptedPasswd, err := EncryptPasswd(key, password)

	postData := url.Values{}
	postData.Set("username", username)
	postData.Set("passwordPre", password)
	postData.Set("password", encryptedPasswd)
	postData.Set("type", "UsernamePassword")
	postData.Set("_eventId", "submit")
	postData.Set("geolocation", "")
	postData.Set("execution", string(execution))
	// missing spelling from hdu, so what can I say?
	postData.Set("croypto", string(key))
	req, err = http.NewRequest(http.MethodPost, ssologinurl, bytes.NewBufferString(postData.Encode()))
	addHeadersToRequest(req, session, url1)
	resp, err = client.Do(req)
	url2 := resp.Header.Get("Location")
	//fmt.Println("URL2:", url2)
	req, err = http.NewRequest(http.MethodGet, url2, nil)
	addHeadersToRequest(req, session, url1)
	resp, err = client.Do(req)
	X_Auth_Token := resp.Header.Get("X-Auth-Token")
	//fmt.Println("X-Auth-Token:", X_Auth_Token)
	return X_Auth_Token
}
func addHeaders(req *http.Request) {
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="117", "Not;A=Brand";v="8"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5938.132 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "close")
}
func addHeaderssso(req *http.Request) {
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="117", "Not;A=Brand";v="8"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5938.132 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "close")
}
func addHeadersToRequest(req *http.Request, sessionID string, refurl string) {
	req.Header.Set("Cookie", "SESSION="+sessionID)
	req.Header.Set("Content-Length", "5324")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://sso.hdu.edu.cn")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5938.132 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", refurl)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "close")
}
func EncryptPasswd(key []byte, password string) (string, error) {
	var keyBytes [8]byte
	if _, err := base64.StdEncoding.Decode(keyBytes[:], key); err != nil {
		return "", fmt.Errorf("decode key: %v", err)
	}
	cipher, err := des.NewCipher(keyBytes[:])
	if err != nil {
		return "", fmt.Errorf("des cipher: %v", err)
	}
	// padding
	text := pkcs7Padding([]byte(password), cipher.BlockSize())
	// ecb mode
	for i := 0; i < len(text); i += cipher.BlockSize() {
		cipher.Encrypt(text[i:], text[i:])
	}
	return base64.StdEncoding.EncodeToString(text), nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}
