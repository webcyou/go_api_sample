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
	DB_NAME                  = "go_api.db"
	DB_TYPE                  = "sqlite3"
	DB_CONNECT_ERROR_MESSAGE = "failed to connect database"
	HELLO_WORLD_MESSAGE      = "Hello world!"
	SERVER_PORT_NUMBER       = ":8080"
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

type MatchingUser struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func NewUser(name string, itemSlice []string) User {
	return User{
		Name:  name,
		Items: CreateItems(itemSlice),
	}
}

func CreateUsers(userNames []string, itemSlice []string) []User {
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
		Score: rand.Intn(10) + 1,
	}
}

func CreateItems(s []string) []Item {
	items := make([]Item, len(s))

	for i, v := range s {
		items[i] = NewItem(v)
	}

	return items
}

func NewMatchingUser(user User, score float64) MatchingUser {
	return MatchingUser{
		ID:    user.ID,
		Name:  user.Name,
		Score: score,
	}
}

// sort
type ByScore []MatchingUser

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

// Euclid Distance Score
func GetDistanceScore(user *User, otherUser *User) float64 {
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
		panic(DB_CONNECT_ERROR_MESSAGE)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Item{})

	itemNames := []string{
		"ゲーム1",
		"ゲーム2",
		"ゲーム3",
		"ゲーム4",
		"ゲーム5",
		"ゲーム6",
		"ゲーム7",
		"ゲーム8",
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

	users := CreateUsers(userNames, itemNames)

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
		panic(DB_CONNECT_ERROR_MESSAGE)
	}
	defer db.Close()

	router := gin.Default()

	// Hello world！
	router.GET("/", func(c *gin.Context) {
		c.String(200, HELLO_WORLD_MESSAGE)
	})

	// userList 取得
	router.GET("/users", func(c *gin.Context) {
		var (
			users   []User
			items   []Item
			jsonMap map[string]interface{} = make(map[string]interface{})
		)

		db.Find(&users)

		for i, user := range users {
			users[i] = setUserItem(user, items, db)
		}

		jsonMap["users"] = users
		c.JSON(200, jsonMap)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var (
			user         User
			otherUsers   []User
			matchingUser []MatchingUser
			items        []Item
			jsonMap      map[string]interface{} = make(map[string]interface{})
		)

		db.First(&user, c.Param("id"))
		user = setUserItem(user, items, db)

		db.Not("id", c.Param("id")).Find(&otherUsers)
		for i, otherUser := range otherUsers {
			otherUsers[i] = setUserItem(otherUser, items, db)
		}

		for _, otherUser := range otherUsers {
			matchingUser = append(matchingUser, NewMatchingUser(otherUser, GetDistanceScore(&user, &otherUser)))
		}

		sort.Sort(sort.Reverse(ByScore(matchingUser)))

		jsonMap["user"] = user
		jsonMap["matching_users"] = matchingUser
		c.JSON(200, jsonMap)
	})

	router.Run(SERVER_PORT_NUMBER)
}
