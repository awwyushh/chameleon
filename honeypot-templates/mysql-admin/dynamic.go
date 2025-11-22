package main

import (
	"fmt"
	"math/rand"
	"time"
)

type FakeDB struct {
	Name   string
	Tables []string
}

func GenerateFakeDB() FakeDB {
	rand.Seed(time.Now().UnixNano())
	db := FakeDB{
		Name: fmt.Sprintf("appdb_%d", rand.Intn(9999)),
		Tables: []string{
			"users",
			"transactions",
			"auth_logs",
			"config",
		},
	}
	return db
}

func GenerateSQLError(payload string) string {
	payloadExcerpt := payload
	if len(payloadExcerpt) > 30 {
		payloadExcerpt = payloadExcerpt[:30] + "..."
	}

	mysqlCodes := []string{
		"1064", "1146", "1054", "1049",
	}

	code := mysqlCodes[rand.Intn(len(mysqlCodes))]

	return fmt.Sprintf(
		"Error %s: You have an error in your SQL syntax near '%s' at line 1",
		code, payloadExcerpt,
	)
}

func SlowDown() {
	ms := rand.Intn(800) + 250
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
