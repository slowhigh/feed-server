package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/team-nerd-planet/feed-server/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type FeedURL struct {
	Name        string
	CompanySize entity.CompanySizeType
	URL         string
}

var FeedURLArr = []FeedURL{
	{Name: "카카오", CompanySize: entity.LARGE, URL: `https://tech.kakao.com/blog/feed`},
	{Name: "쿠팡", CompanySize: entity.LARGE, URL: `https://medium.com/feed/coupang-engineering`},
	{Name: "왓챠", CompanySize: entity.SMALL, URL: `https://medium.com/feed/watcha`},
	{Name: "컬리", CompanySize: entity.SMALL, URL: `https://helloworld.kurly.com/feed`},
	{Name: "우아한형제들", CompanySize: entity.STARTUP, URL: `https://techblog.woowahan.com/feed`},
	{Name: "뱅크샐러드", CompanySize: entity.STARTUP, URL: `https://blog.banksalad.com/rss.xml`},
	{Name: "NHN", CompanySize: entity.LARGE, URL: `https://meetup.nhncloud.com/rss`},
	{Name: "하이퍼커넥트", CompanySize: entity.LARGE, URL: `https://hyperconnect.github.io/feed`},
	{Name: "당근마켓", CompanySize: entity.STARTUP, URL: `https://medium.com/feed/daangn`},
	{Name: "강남언니", CompanySize: entity.SMALL, URL: `https://blog.gangnamunni.com/blog`},
	{Name: "요기요", CompanySize: entity.STARTUP, URL: `https://techblog.yogiyo.co.kr/feed`},
	{Name: "이스트소프트", CompanySize: entity.MEDIUM, URL: `https://blog.est.ai/feed`},
	{Name: "플랫팜", CompanySize: entity.STARTUP, URL: `https://medium.com/feed/platfarm`},
	{Name: "직방", CompanySize: entity.MEDIUM, URL: `https://medium.com/feed/zigbang`},
	{Name: "스포카", CompanySize: entity.STARTUP, URL: `https://spoqa.github.io/rss`},
	{Name: "네이버플레이스", CompanySize: entity.LARGE, URL: `https://medium.com/feed/naver-place-dev`},
	{Name: "라인", CompanySize: entity.LARGE, URL: `https://engineering.linecorp.com/ko/feed/index.html`},
	{Name: "쏘카", CompanySize: entity.STARTUP, URL: `https://tech.socarcorp.kr/feed`},
	{Name: "리디", CompanySize: entity.SMALL, URL: `https://www.ridicorp.com/feed`},
	{Name: "네이버", CompanySize: entity.LARGE, URL: `https://d2.naver.com/d2.atom`},
	{Name: "데보션", CompanySize: entity.LARGE, URL: `https://devocean.sk.com/blog/rss.do`},
	{Name: "구글코리아", CompanySize: entity.LARGE, URL: `https://feeds.feedburner.com/GoogleDevelopersKorea`},
	{Name: "AWS코리아", CompanySize: entity.LARGE, URL: `https://aws.amazon.com/ko/blogs/tech/feed`},
	{Name: "무신사", CompanySize: entity.MEDIUM, URL: `https://medium.com/feed/musinsa-tech`},
	{Name: "데이블", CompanySize: entity.SMALL, URL: `https://teamdable.github.io/techblog/feed`},
	{Name: "토스", CompanySize: entity.STARTUP, URL: `https://toss.tech/rss.xml`},
	{Name: "스마일게이트", CompanySize: entity.MEDIUM, URL: `https://smilegate.ai/recent/feed`},
	{Name: "롯데온", CompanySize: entity.LARGE, URL: `https://techblog.lotteon.com/feed`},
	{Name: "카카오엔터프라이즈", CompanySize: entity.LARGE, URL: `https://tech.kakaoenterprise.com/feed`},
	{Name: "메가존클라우드", CompanySize: entity.MEDIUM, URL: `https://www.megazone.com/blog/feed`},
	{Name: "SKC&C", CompanySize: entity.LARGE, URL: `https://engineering-skcc.github.io/feed.xml`},
	{Name: "여기어때", CompanySize: entity.STARTUP, URL: `https://techblog.gccompany.co.kr/feed`},
	{Name: "원티드", CompanySize: entity.STARTUP, URL: `https://medium.com/feed/wantedjobs`},
	{Name: "29CM", CompanySize: entity.MEDIUM, URL: `https://medium.com/feed/29cm`},
	{Name: "비브로스", CompanySize: entity.STARTUP, URL: `https://boostbrothers.github.io/rss`},
	{Name: "포스타입", CompanySize: entity.STARTUP, URL: `https://team.postype.com/rss`},
	{Name: "지마켓", CompanySize: entity.LARGE, URL: `https://dev.gmarket.com/feed`},
	{Name: "SK플래닛", CompanySize: entity.LARGE, URL: `https://techtopic.skplanet.com/rss`},
	{Name: "AB180", CompanySize: entity.MEDIUM, URL: `https://raw.githubusercontent.com/ab180/engineering-blog-rss-scheduler/main/rss.xml`},
	{Name: "데브시스터즈", CompanySize: entity.SMALL, URL: `https://tech.devsisters.com/rss.xml`},
	{Name: "테이블링", CompanySize: entity.SMALL, URL: `https://techblog.tabling.co.kr/feed`},
	{Name: "넷마블", CompanySize: entity.LARGE, URL: `https://netmarble.engineering/feed`},
	{Name: "마키나락스", CompanySize: entity.SMALL, URL: `https://www.makinarocks.ai/blog/feed`},
	{Name: "드라마앤컴패니", CompanySize: entity.STARTUP, URL: `https://blog.dramancompany.com/feed`},
	{Name: "티몬", CompanySize: entity.MEDIUM, URL: `https://rss.blog.naver.com/tmondev.xml`},
	{Name: "루닛", CompanySize: entity.SMALL, URL: `https://medium.com/feed/lunit`},
	{Name: "야놀자", CompanySize: entity.MEDIUM, URL: `https://medium.com/feed/yanoljacloud-tech`},
	{Name: "인프런", CompanySize: entity.SMALL, URL: `https://tech.inflab.com/rss.xml`},
}

func main() {
	r := regexp.MustCompile("<[^>]*>")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul",
		"localhost",
		5432,
		"nerd",
		"planet1!",
		"nerd_planet",
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  false,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if db == nil {
		slog.Error("db is null")
		return
	}

	if err := db.AutoMigrate(&entity.Rss{}, &entity.Feed{}, &entity.Item{}, &entity.JobTag{}, &entity.SkillTag{}); err != nil {
		slog.Error(err.Error())
		return
	}

	newJobTags := []entity.JobTag{
		{Name: "FE", Keyword: []string{"fe", "frontend", "javascript", "typescript", "html", "css"}},
		{Name: "BE", Keyword: []string{"be", "backend", "java", "go", "rust", "python", "spring boot", "gin", "fastapi"}},
		{Name: "DEVOPS", Keyword: []string{"devops", "docker", "kubernetes", "k8s"}},
		{Name: "SECURITY", Keyword: []string{"security", "sdl", "취약점", "보안", "해킹", "해커"}},
		{Name: "DATA", Keyword: []string{"data", "전처리", "streamlit", "snowflake", "sql", "mysql", "mssql", "postgresql"}},
		{Name: "AI", Keyword: []string{"ai", "인공지능", "tensorflow", "nlp", "자연어"}},
		{Name: "LLM", Keyword: []string{"llm", "rnn", "lstm", "chatgpt", "gemini", "gpt", "openai", "프롬프트"}},
	}

	if err := db.Create(&newJobTags).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	newSkillTags := []entity.SkillTag{
		{Name: "TYPESCRIPT", Keyword: []string{"typescript"}},
		{Name: "PYTHON", Keyword: []string{"python"}},
		{Name: "KOTLIN", Keyword: []string{"kotlin"}},
		{Name: "GO", Keyword: []string{"go", "golang"}},
		{Name: "RUBY", Keyword: []string{"ruby"}},
		{Name: "C++", Keyword: []string{"c++"}},
		{Name: "C", Keyword: []string{"c언어", "c 언어"}},
		{Name: "JAVA", Keyword: []string{"java"}},
		{Name: "C#", Keyword: []string{"c#"}},
		{Name: "PHP", Keyword: []string{"php"}},
	}

	if err := db.Create(&newSkillTags).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("create feed, item entity")

	fp := gofeed.NewParser()
	for _, feedURL := range FeedURLArr {
		newRss := entity.Rss{
			Name:    feedURL.Name,
			Link:    feedURL.URL,
			Updated: time.Now(),
		}

		feed, err := fp.ParseURL(feedURL.URL)
		if err != nil {
			slog.Error(err.Error())
			slog.Error("get feed error", "url", feedURL)
			newRss.Ok = false
			newRss.Error = err.Error()
			continue
		}

		newRss.Ok = true

		if err := db.Create(&newRss).Error; err != nil {
			slog.Error(err.Error())
			return
		}

		var feedUpdated time.Time

		if feed.UpdatedParsed != nil {
			feedUpdated = *feed.UpdatedParsed
		} else if feed.PublishedParsed != nil {
			feedUpdated = *feed.PublishedParsed
		} else {
			feedUpdated = time.Now()
		}

		newItems := make([]entity.Item, len(feed.Items))

		for i, item := range feed.Items {
			var itemPublished time.Time

			if item.PublishedParsed != nil {
				itemPublished = *item.PublishedParsed
			} else if item.UpdatedParsed != nil {
				itemPublished = *item.UpdatedParsed
			} else {
				itemPublished = time.Now()
			}

			relatedJobTags := make([]entity.JobTag, 0)
			for _, jobTag := range newJobTags {
				for _, keyword := range jobTag.Keyword {
					if strings.Contains(strings.ToLower(r.ReplaceAllString(item.Title, "")), keyword) {
						relatedJobTags = append(relatedJobTags, jobTag)
						break
					} else if strings.Contains(strings.ToLower(r.ReplaceAllString(item.Description, "")), keyword) {
						relatedJobTags = append(relatedJobTags, jobTag)
						break
					}
				}
			}

			relatedSkillTags := make([]entity.SkillTag, 0)
			for _, skillTag := range newSkillTags {
				for _, keyword := range skillTag.Keyword {
					if strings.Contains(strings.ToLower(r.ReplaceAllString(item.Title, "")), keyword) {
						relatedSkillTags = append(relatedSkillTags, skillTag)
						break
					} else if strings.Contains(strings.ToLower(r.ReplaceAllString(item.Description, "")), keyword) {
						relatedSkillTags = append(relatedSkillTags, skillTag)
						break
					}
				}
			}

			var thumbnail *string
			if item.Image != nil {
				if !strings.Contains(item.Image.URL, "http://") && !strings.Contains(item.Image.URL, "https://") {
					url := feed.Link + item.Image.URL
					thumbnail = &url
				} else {
					thumbnail = &item.Image.URL
				}
			} else {
				thumbnail = nil
			}

			newItems[i] = entity.Item{
				Title:       item.Title,
				Description: item.Description,
				Link:        item.Link,
				Thumbnail:   thumbnail,
				Published:   itemPublished,
				GUID:        item.GUID,
				JobTags:     relatedJobTags,
				SkillTags:   relatedSkillTags,
			}
		}

		newFeed := entity.Feed{
			Name:        feedURL.Name,
			Title:       feed.Title,
			Description: feed.Description,
			Link:        feed.Link,
			Updated:     feedUpdated,
			Copyright:   feed.Copyright,
			CompanySize: feedURL.CompanySize,
			Items:       newItems,
			RssID:       newRss.ID,
		}

		if err := db.Create(&newFeed).Error; err != nil {
			slog.Error(err.Error())
		}
		slog.Info("create feed", "url", feed.Link)
	}

	query := db.
		Table("items i").
		Select(`
			i."id" as item_id, 
			i.title as item_title, 
			i.description as item_description, 
			i."link" as item_link,
			i.thumbnail as item_thumbnail,
	    	i.published as item_published,
			f."id" as feed_id,
			f."name" as feed_name, 
			f.title as feed_title, 
			f."link" as feed_link,
	    	f.company_size as company_size, 
			job_tags.id_arr as job_tags_id_arr,
			skill_tags.id_arr as skill_tags_id_arr
			`).
		Joins(`LEFT JOIN feeds f ON f."id" = i.feed_id`).
		Joins(`LEFT JOIN (?) as job_tags ON job_tags.item_id = i."id"`,
			db.Table("item_job_tags").
				Select("item_id, array_agg(job_tag_id) as id_arr").
				Group("item_id")).
		Joins(`LEFT JOIN (?) as skill_tags ON skill_tags.item_id = i."id"`,
			db.Table("item_skill_tags").
				Select("item_id, array_agg(skill_tag_id) as id_arr").
				Group("item_id")).
		Order(`i.published desc, i."id" desc`)

	db.Migrator().CreateView("vw_items", gorm.ViewOption{Query: query, Replace: true})

	var itemView []entity.ItemView

	if err := db.Find(&itemView).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	// for _, item := range itemView {
	// 	fmt.Printf("%+v\n", item)
	// }
}
