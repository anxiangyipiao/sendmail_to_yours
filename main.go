package main

import (
	"os"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"gopkg.in/gomail.v2"
	"github.com/robfig/cron/v3"
)


// 模板字符串
const templateString = `
<!DOCTYPE html>
<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Responsive Content</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 20px;
		}
		.container {
			max-width: 100%;
			margin: 0 auto;
		}
		img {
			max-width: 100%;
			height: auto;
		}
		@media (min-width: 768px) {
			.container {
				max-width: 768px;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<p>{{.LoveSentence}}</p>
		<p>每日开心一下：{{.Xiaohua}}</p>
		<img src="{{.Image}}" alt="Responsive Image">
		
	</div>
</body>
</html>
`
// <p>{{.Dog}}</p>

// 生成 HTML 内容
func generateHTML(templateString string, data Custom) (string, error) {
	tmpl, err := template.New("htmlTemplate").Parse(templateString)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}


//数据模型
type Custom struct {

	LoveSentence string
	Dog string
	Image string
	Xiaohua string
	
}

// API 响应数据模型
type APIResponse struct {
	Data struct {
		Text string `json:"text"`
		Url string `json:"url"`
	} `json:"data"`
}


// 获得数据函数
func fetchDataFromAPI(apiURL string, target interface{}){
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		log.Fatal(err)
	}

}

//"https://www.dmoe.cc/random.php?return=json"
//"http://api.btstu.cn/yan/api.php?charset=utf-8&encode=json"
// 生成你的数据
func generateData() Custom {


	var api1 APIResponse	//情话
	fetchDataFromAPI("https://api.gumengya.com/Api/LoveSentence?format=json", &api1)

	var api2 APIResponse	//狗狗
	fetchDataFromAPI("https://api.gumengya.com/Api/Dog", &api2)

	var api3 APIResponse     //图片
	fetchDataFromAPI("https://api.gumengya.com/Api/FjImg", &api3)

	var api4 APIResponse	//笑话
	fetchDataFromAPI("https://api.gumengya.com/Api/Xiaohua", &api4)
	// fetchDataFromAPI("https://www.mxnzp.com/api/jokes/list?page=1&app_id="+APPID+"&app_secret="+APPSECRET, &custom)
	
	var custom Custom
	custom.LoveSentence = api1.Data.Text
	custom.Dog = api2.Data.Text
	custom.Image = api3.Data.Url
	custom.Xiaohua = api4.Data.Text


	return custom

}

// 发送邮件
func sendEmail(htmlContent, recipientEmail string, config Config) {


	sender := gomail.NewMessage()
	sender.SetHeader("From", config.SenderEmail)
	sender.SetHeader("To", recipientEmail)
	sender.SetHeader("Subject", "Daily Report")
	sender.SetBody("text/html", htmlContent)

	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort ,config.SenderEmail, config.SenderPassword)

	if err := dialer.DialAndSend(sender); err != nil {
		log.Println("Error sending email to", recipientEmail, ":", err)
	} else {
		log.Println("Email sent successfully to", recipientEmail)
	}
}

// 配置信息
type Config struct {
	SenderEmail     string   `json:"senderEmail"`
	SenderPassword  string   `json:"senderPassword"`
	SMTPHost        string   `json:"smtpHost"`
	SMTPPort        int      `json:"smtpPort"`
	RecipientEmails []string `json:"recipientEmails"`
	Times           string   `json:"times"`
}

// 读取配置文件
func readConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {

	// 读取配置文件
	config, err := readConfig("config.json")
	if err != nil {
			log.Println("Error reading config file:", err)
		return
	}


	c := cron.New()

	// 添加一个定时任务
	_,err = c.AddFunc(config.Times, func() {
		
		// 获取待发送的数据
		data := generateData()
		
		// 生成 HTML 内容
		htmlContent, err := generateHTML(templateString, data)
		if err != nil {
			log.Fatal(err)
		}

		// 发送邮件给每个收件人
		for _, recipientEmail := range config.RecipientEmails {
			sendEmail(htmlContent, recipientEmail, config)
		}
	})

	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	// 启动定时任务
	c.Start()

	// 保持主程序运行
	select {}

}








