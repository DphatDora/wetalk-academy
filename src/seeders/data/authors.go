package data

import "wetalk-academy/internal/domain/model"

var authors = []model.TopicAuthor{
	{UserID: 11001, Name: "Nguyễn Minh An", Avatar: "https://i.pravatar.cc/300?img=12"},
	{UserID: 11002, Name: "Trần Hoàng Nam", Avatar: "https://i.pravatar.cc/300?img=14"},
	{UserID: 11003, Name: "Lê Thu Hà", Avatar: "https://i.pravatar.cc/300?img=47"},
	{UserID: 11004, Name: "Phạm Quốc Bảo", Avatar: "https://i.pravatar.cc/300?img=52"},
	{UserID: 11005, Name: "Vũ Ngọc Linh", Avatar: "https://i.pravatar.cc/300?img=32"},
	{UserID: 11006, Name: "Đặng Tuấn Kiệt", Avatar: "https://i.pravatar.cc/300?img=11"},
	{UserID: 11007, Name: "Bùi Thanh Mai", Avatar: "https://i.pravatar.cc/300?img=45"},
}

// The assignment is randomized once and kept stable so repeated seeds do not
// unexpectedly change topic ownership.
var authorIndexByTopic = map[string]int{
	"go-programming":             3,
	"backend-go":                 5,
	"modern-js-ts":               1,
	"react-engineering":          4,
	"python-automation-data":     2,
	"devops-cicd":                0,
	"database-engineering":       6,
	"system-design":              1,
	"cloud-native-microservices": 3,
	"web-security":               5,
}

func authorFor(topicKey string) model.TopicAuthor {
	if index, exists := authorIndexByTopic[topicKey]; exists {
		return authors[index]
	}
	id := stableID("author", topicKey)
	return authors[int(id[0])%len(authors)]
}
