package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	DB_NAME            = "go_api.db"
	DB_TYPE            = "sqlite3"
	DB_CONNECT_ERROR   = "failed to connect database"
	SERVER_PORT_NUMBER = ":8080"
)

type User struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

type Item struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

type RecommendUser struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func NewUser(name string, itemSlice []string) User {
	return User{
		Name:  name,
		Items: createItems(itemSlice),
	}
}

func createUsers(userNames []string, itemSlice []string) []User {
	users := make([]User, len(userNames))

	for i, name := range userNames {
		users[i] = NewUser(name, itemSlice)
	}

	return users
}

func NewItem(name string) Item {
	rand.Seed(time.Now().UnixNano())

	return Item{
		Name:  name,
		Score: rand.Intn(10),
	}
}

func createItems(s []string) []Item {
	items := make([]Item, len(s))

	for i, v := range s {
		items[i] = NewItem(v)
	}

	return items
}

func NewRecommendUser(user User, score float64) RecommendUser {
	return RecommendUser{
		ID:    user.ID,
		Name:  user.Name,
		Score: score,
	}
}

// sort
type ByScore []RecommendUser

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

// Euclid Distance Score
func getDistanceScore(user *User, otherUser *User) float64 {
	var matchItems map[string]float64 = make(map[string]float64)

	for _, item := range user.Items {
		for _, otherItem := range otherUser.Items {
			if item.Name == otherItem.Name {
				matchItems[item.Name] = math.Pow(float64(item.Score-otherItem.Score), 2)
			}
		}
	}

	if len(matchItems) == 0 {
		return 0
	}

	var sumOfSquares float64
	for _, matchItem := range matchItems {
		sumOfSquares += matchItem
	}

	return (1 / (1 + math.Sqrt(sumOfSquares)))
}

func init() {
	db, err := gorm.Open(DB_TYPE, DB_NAME)
	if err != nil {
		panic(DB_CONNECT_ERROR)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Item{})

	itemNames := []string{
		"マリオブラザース",
		"スーパーマリオブラザース",
		"ゼルダの伝説",
		"アイスクライマー",
		"エキサイトバイク",
		"パックマン",
		"魔界村",
		"ドクターマリオ",
	}

	userNames := []string{
		"ユーザー1",
		"ユーザー2",
		"ユーザー3",
		"ユーザー4",
		"ユーザー5",
		"ユーザー6",
		"ユーザー7",
		"ユーザー8",
		"ユーザー9",
		"ユーザー10",
	}

	users := createUsers(userNames, itemNames)

	count := 0
	db.Table("users").Count(&count)

	if count == 0 {
		for _, user := range users {
			db.Create(&user)
		}
	}
}

func setUserItem(user User, items []Item, db *gorm.DB) User {
	db.Model(&user).Related(&items)
	user.Items = items

	return user
}

func main() {
	db, err := gorm.Open(DB_TYPE, DB_NAME)
	if err != nil {
		panic(DB_CONNECT_ERROR)
	}
	defer db.Close()

	router := gin.Default()

	// Hello world！
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello world!")
	})

	// userList 取得
	router.GET("/users", func(c *gin.Context) {
		var (
			users []User
			items []Item
		)

		db.Find(&users)

		for i, user := range users {
			users[i] = setUserItem(user, items, db)
		}

		c.JSON(200, users)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var (
			user          User
			otherUsers    []User
			recommendUser []RecommendUser
			items         []Item
		)

		db.First(&user, c.Param("id"))
		user = setUserItem(user, items, db)

		db.Not("id", c.Param("id")).Find(&otherUsers)
		for i, otherUser := range otherUsers {
			otherUsers[i] = setUserItem(otherUser, items, db)
		}

		for _, otherUser := range otherUsers {
			recommendUser = append(recommendUser, NewRecommendUser(otherUser, getDistanceScore(&user, &otherUser)))
		}

		sort.Sort(sort.Reverse(ByScore(recommendUser)))

		c.JSON(200, user)
		c.JSON(200, recommendUser)
	})

	router.Run(SERVER_PORT_NUMBER)
}
