package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"hotNews/cache"
	mysql "hotNews/db"
	"hotNews/http/models"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type TopResult struct {
	Fresh_Text string    `json:"fresh_text"`
	Paging     Paging    `json:"paging"`
	Data       []TopData `json:"data"`
}

type Paging struct {
	Is_End   bool   `json:"is_end"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type TopData struct {
	Detail_Text string     `json:"detail_text"`
	Target      Target     `json:"target"`
	Children    []Children `json:"children"`
}

type Target struct {
	Title string `json:"title"`
	Id    int    `json:"id"`
}

type Children struct {
	Thumbnail string `json:"thumbnail"`
	Type      string `json:"type"`
}

type Health struct {
	Id    int
	Title string
	Cover string
	Hit   int
	Url   string
	Date  string
}

func ZhTop() {
	limit := 1
	key := "zhihu-top"
	redis_key := "limiter:zhihu-top" + key
	res, _ := cache.Get(redis_key)
	int_res, _ := strconv.Atoi(string(res))
	if int_res > 0 {
		fmt.Println(key + "数据抓取完成")
		return
	}
	db := mysql.DbEngin
	application := model.Application{}
	//获取app数据 增加计算时间
	appDb := db.Where("alias = ?", key)
	appDb.First(&application)
	if application.Id < 0 {
		return
	}
	number := cache.Limiter(key, limit, application.Polling)
	if !number {
		return
	}
	body, _ := Curl("GET", application.Url, "")
	var result TopResult
	_ = json.Unmarshal([]byte(string(body)), &result)

	article := model.Article{}
	for _, item := range result.Data {
		article.TargetId = strconv.Itoa(item.Target.Id)
		article.ApplicationId = 1
		tmpDb := db.Where("target_id = ?", article.TargetId).Where("application_id = ?", article.ApplicationId)
		tmpDb.First(&article)

		hit, _ := strconv.Atoi(string(regexp.MustCompile("\\d+").Find([]byte(item.Detail_Text))))
		article.Hit = hit * 10000
		article.Title = item.Target.Title
		jsonBytes, _ := json.Marshal(item)
		article.Json = string(jsonBytes)
		if len(item.Children) > 0 {
			article.Cover = item.Children[0].Thumbnail
		}
		if article.Id > 0 {
			tmpDb.Save(&article)
		} else {
			article.CreatedAt = time.Now()
			tmpDb.Create(&article)
		}
		article.Id = 0
	}
	//ReturnJson(c, http.StatusOK, "success", result.Data)
}

func QueryHtml(c *gin.Context) {

	for i := 0; i < 50; i++ {

		res, err := http.Get("https://www.cnys.com/article/list_2_" + strconv.Itoa(i) + ".html")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		db := mysql.DbEngin
		var health Health
		article := model.Article{}
		article.ApplicationId = 2
		domain := "https://www.cnys.com"
		doc.Find(".leftLists a").Each(func(i int, s *goquery.Selection) {

			href, _ := s.Attr("href")
			article.TargetId = string(regexp.MustCompile("\\d+").Find([]byte(href)))

			tmpDb := db.Where("target_id = ?", article.TargetId).Where("application_id = ?", article.ApplicationId)
			tmpDb.First(&article)

			article.Cover, _ = s.Find("img").Attr("data-original")
			article.Title = s.Find("h2").Text()

			hit, _ := strconv.Atoi(string(regexp.MustCompile("\\d+").Find([]byte(s.Find(".leftListTip span").Text()))))
			article.Hit = hit * 10000

			date := s.Find(".nationalListText span").Text()
			health.Date = string(regexp.MustCompile("\\d+\\-\\d+\\-\\d+").Find([]byte(date)))

			timeLayout := "2006-01-02 15:04:05"  //转化所需模板
			loc, _ := time.LoadLocation("Local") //重要：获取时区
			article.CreatedAt, _ = time.ParseInLocation(timeLayout, health.Date, loc)

			if article.Id > 0 {
				tmpDb.Save(&article)
			} else {
				tmpDb.Create(&article)
			}
			article.Id = 0
			health.Url = domain + href

			fmt.Println(health)
		})
	}
	ReturnJson(c, http.StatusOK, "success", "")
}

func Detail(c *gin.Context) {
	db := mysql.DbEngin

	rows, _ := db.Model(&model.Article{}).Where("application_id = 2").Rows()
	defer rows.Close()
	for rows.Next() {
		modelArticle := model.Article{}
		_ = db.ScanRows(rows, &modelArticle)

		url := "https://www.cnys.com/article/" + modelArticle.TargetId + ".html"

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		article := model.Article{}

		doc.Find(".readbox").Each(func(i int, s *goquery.Selection) {

			pp := s.Find(".digest p").Text()
			dd, _ := s.Find(".reads").Html()
			fmt.Println(pp)
			fmt.Println(dd)
			tmpDb := db.Where("id = ?", modelArticle.Id)
			tmpDb.First(&article)
			article.Json = string(dd)
			tmpDb.Save(&article)

		})
	}
	ReturnJson(c, http.StatusOK, "success", "")
}

/*
百度文库
*/
func Wenku(c *gin.Context) {
	url := "https://wenku.baidu.com/view/f618bfa2eef9aef8941ea76e58fafab069dc4482.html"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	d, _ := ioutil.ReadAll(res.Body)
	//reader := transform.NewReader(bytes.NewReader(d), simplifiedchinese.GBK.NewDecoder())

	//doc, err := goquery.NewDocumentFromReader(reader)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//doc.Find("#bd").Each(func(i int, s *goquery.Selection) {
	//	init_number := 3
	//
	//	for i := 3; i <= init_number; i++ {
	//		fmt.Println("循环：", i)
	//		target := ".reader-page-" + strconv.Itoa(i)
	//		fmt.Println(target)
	//		html, _ := s.Find(target).Html()
	//		fmt.Println(html)
	//	}
	//	more_btn, _ := s.Find(".moreBtn").Html()
	//	left_number, _ := strconv.Atoi(string(regexp.MustCompile("\\d+").Find([]byte(more_btn))))
	//
	//	if left_number > 0 {
	//
	//	}
	//
	//	fmt.Println("剩下页数", left_number)
	//})
	ReturnJson(c, http.StatusOK, "success", d)
}
