package entity

import "time"

type CompanySizeType int

const (
	STARTUP CompanySizeType = iota
	SMALL
	MEDIUM
	LARGE
	FOREIGN
)

type Feed struct {
	ID          uint            `gorm:"column:id;primarykey"`
	Name        string          `gorm:"column:name;type:varchar;not null"`
	Title       string          `gorm:"column:title;type:varchar;not null"`
	Description string          `gorm:"column:description;type:varchar;not null"`
	Link        string          `gorm:"column:link;type:varchar;not null;unique"`
	Updated     time.Time       `gorm:"column:updated;type:timestamp;not null"`
	Copyright   string          `gorm:"column:copyright;type:varchar;not null"`
	CompanySize CompanySizeType `gorm:"column:company_size;type:int2;not null"`
	RssID       uint            `gorm:"column:rss_id;type:int8;not null"`
	Items       []Item          `gorm:"foreignKey:FeedID"`
}
