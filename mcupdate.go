package mcupdate

import (
	"bufio"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	instance = &logging{}
	bot.RegisterModule(instance)
}

type logging struct {
}

func (m *logging) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "internal.mcupdate",
		Instance: instance,
	}
}

func (m *logging) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *logging) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *logging) Serve(b *bot.Bot) {
	// 注册服务函数部分
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		go update(c,msg)
	})
	go urls(b)
}

func (m *logging) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *logging) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *logging

var logger = utils.GetModuleLogger("internal.logging")


var build strings.Builder

func update(c *client.QQClient, msg *message.GroupMessage){
	beta := ""
	release := ""
	s3 := ""
	s33 := ""
	resp, err := http.Get("https://feedback.minecraft.net/hc/en-us/categories/115000410252-Knowledge-Base")
	if err != nil {
		log.Println("http get error.")
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("http read error.")
		}
		src := string(body)
		docs := strings.NewReader(src)
		doc, err := html.Parse(docs)
		if err != nil {
			fmt.Fprint(os.Stderr, "findlinks: %v\n", err)
			os.Exit(1)
		}
	OuterLoop:
		for _, link := range visit(nil, doc) {
			rds := bufio.NewReader(strings.NewReader(link))
			for {
				lined, err := rds.ReadString('\n')
				res1 := strings.Contains(lined, "-Xbox-One-")
				if err == nil || io.EOF == err {
					if res1 == true {
						beta = strings.Replace(GetBetweenStr(lined, "1-1", "-Xbox"), "-", ".", 3)
						break OuterLoop
					}
				}
				break
			}
		}
	OuterLoops:
		for _, link := range visit(nil, doc) {
			rds := bufio.NewReader(strings.NewReader(link))
			for {
				lined, err := rds.ReadString('\n')
				res2 := strings.Contains(lined, "-Bedrock-")
				if err == nil || io.EOF == err {
					if res2 == true {
						release = strings.Replace(GetBetweenStr(lined, "1-1", "-Bed"), "-", ".", 2)
						break OuterLoops
					}
				}
				break
			}
		}
		for _, link := range visit(nil, doc) {
		rds := bufio.NewReader(strings.NewReader(link))
		for {
			lined, err := rds.ReadString('\n')
			res1 := strings.Contains(lined, "-Xbox-One-")
			if err == nil || io.EOF == err {
				if res1 == true {
					betas := strings.Replace(GetBetweenStr(lined, "1-1", "-Xbox"), "-", ".", 3)
					build.WriteString(betas)
					build.WriteString("、")
					s3 = build.String()
				}
			}
			break
		}
	}
		build.Reset()
		for _, link := range visit(nil, doc) {
			rds := bufio.NewReader(strings.NewReader(link))
			for {
				lined, err := rds.ReadString('\n')
				res2 := strings.Contains(lined, "-Bedrock-")
				if err == nil || io.EOF == err {
					if res2 == true {
						res3 := strings.Contains(lined, "-Minecraft-1-")
						if res3 == true {
							releases := strings.Replace(GetBetweenStr(lined, "1-1", "-Bed"), "-", ".", 2)
							build.WriteString(releases)
							build.WriteString("、")
							s33 = build.String()
						}
					}
				}
				break
			}
		}
	}
	build.Reset()
	out := msg.ToString()
	if out == "查询最新版本" {
		m := message.NewSendingMessage().Append(message.NewText("最新测试版："+beta+"\n最新稳定版："+release))
		c.SendGroupMessage(msg.GroupCode, m)
	}
	if out == "查询历史版本" {
		m := message.NewSendingMessage().Append(message.NewText("近五个版本测试版："+s3+"\n近五个版本正式版："+s33))
		c.SendGroupMessage(msg.GroupCode, m)
	}
}



func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func urls(b *bot.Bot) {
	eexx := false
	eexxx := false
	betas := ""
	releases := ""
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			beta := ""
			resp, err := http.Get("https://feedback.minecraft.net/hc/en-us/categories/115000410252-Knowledge-Base")
			if err != nil {
				log.Println("http get error.")
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("http read error.")
				}
				src := string(body)
				docs := strings.NewReader(src)
				doc, err := html.Parse(docs)
				if err != nil {
					fmt.Fprint(os.Stderr, "findlinks: %v\n", err)
					os.Exit(1)
				}
			OuterLoop:
				for _, link := range visit(nil, doc) {
					rds := bufio.NewReader(strings.NewReader(link))
					for {
						lined, err := rds.ReadString('\n')
						res1 := strings.Contains(lined, "-Xbox-One-")
						if err == nil || io.EOF == err {
							if res1 == true {
								beta = strings.Replace(GetBetweenStr(lined, "1-1", "-Xbox"), "-", ".", 3)
								if eexx == false {
									betas = beta
									eexx = true
								}
								if beta != betas {
									betas = beta
									m := message.NewSendingMessage().Append(message.NewText("发现最新Beta版：" + beta))
									b.SendGroupMessage(823851812, m)
									break OuterLoop
								}
							}
						}
						break
					}
				}
			OuterLoops:
				for _, link := range visit(nil, doc) {
					rds := bufio.NewReader(strings.NewReader(link))
					for {
						lined, err := rds.ReadString('\n')
						res2 := strings.Contains(lined, "-Bedrock-")
						if err == nil || io.EOF == err {
							if res2 == true {
								release := strings.Replace(GetBetweenStr(lined, "1-1", "-Bed"), "-", ".", 2)
								if eexxx == false {
									releases = release
									eexxx = true
								}
								if release != releases {
									releases = release
									m := message.NewSendingMessage().Append(message.NewText("发现最新正式版：" + beta))
									b.SendGroupMessage(823851812, m)
									break OuterLoops
								}
							}
							break
						}
					}
				}
				t1.Reset(time.Second * 60)
			}
		}
	}
}