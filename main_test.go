package main

import (
	"testing"
	"sort"
)

func TestNewUser(t *testing.T) {
	user := NewUser("test user", []string{ "test item 1", "test item 2" })

	if len(user.Items) != 2 {
		t.Error("Item数が異なる")
	}

	if user.Items[0].Name != "test item 1" {
		t.Fatalf("ItemNameが異なる", user.Items[0])
	}
}

func BenchmarkGetDistanceScore(b *testing.B) {
	var matchingUser []MatchingUser

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

	user := NewUser("test user", itemNames)
	otherUsers := CreateUsers(userNames, itemNames)

	for i := 0; i < b.N; i++ {
		for _, otherUser := range otherUsers {
			matchingUser = append(matchingUser, NewMatchingUser(otherUser, GetDistanceScore(&user, &otherUser)))
		}
	}

	sort.Sort(sort.Reverse(ByScore(matchingUser)))
}