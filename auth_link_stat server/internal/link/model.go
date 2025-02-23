package link

import (
	"math/rand"
	"restapi/internal/stat"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//поставили каскадную связь между таблицами что не позволит просто так удалить ссылку, так как она может относиться ко множеству статистик
	//ограничения некритичны
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenereteHash()
	return link
}
func (link *Link) GenereteHash() {
	link.Hash = RandStringRunes(10)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
