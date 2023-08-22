package main

import (

	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"gopkg.in/gomail.v2"
	"github.com/robfig/cron/v3"
)

// 配置信息
var APPID string = "rmjvslk7ktj1gllc"
var APPSECRET  =  ""
var your_sender_email = "@qq.com"
var your_sender_password = ""
var smtpHost = "smtp.qq.com"
var smtpPort = 465
var recipientEmails = []string{"qq.com"}
var times = "22 22 * * *"         //<分钟> <小时> <日期> <月份> <星期>


// 模板字符串
var templateString = `
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
		<p>{{.Text}}</p>	
		<img src="{{.Image}}" alt="Responsive Image">	
	</div>
</body>
</html>
`

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
	Text string `json:"text"`
	Image string `json:"imgurl"`
}

// 获得数据函数
func fetchDataFromAPI(apiURL string, target interface{}) {
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


// 生成你的数据
func generateData() Custom {

	var custom Custom

	fetchDataFromAPI("https://www.dmoe.cc/random.php?return=json", &custom)

	fetchDataFromAPI("http://api.btstu.cn/yan/api.php?charset=utf-8&encode=json", &custom)

	// fetchDataFromAPI("https://www.mxnzp.com/api/jokes/list?page=1&app_id="+APPID+"&app_secret="+APPSECRET, &custom)
	
	return custom

}

// 发送邮件
func sendEmail(htmlContent, recipientEmail string) {

	senderEmail := your_sender_email
	senderPassword := your_sender_password

	sender := gomail.NewMessage()
	sender.SetHeader("From", senderEmail)
	sender.SetHeader("To", recipientEmail)
	sender.SetHeader("Subject", "Daily Report")
	sender.SetBody("text/html", htmlContent)

	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	if err := dialer.DialAndSend(sender); err != nil {
		log.Println("Error sending email to", recipientEmail, ":", err)
	} else {
		log.Println("Email sent successfully to", recipientEmail)
	}
}



func main() {

	c := cron.New()

	// 添加一个定时任务，每天中午8点发送邮件
	_, err := c.AddFunc(times, func() {
		// 获取待发送的数据
		data := generateData()

		// 生成 HTML 内容
		htmlContent, err := generateHTML(templateString, data)
		if err != nil {
			log.Fatal(err)
		}

	
		// 发送邮件给每个收件人
		for _, recipientEmail := range recipientEmails {
			sendEmail(htmlContent, recipientEmail)
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








