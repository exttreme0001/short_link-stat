package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"restapi/internal/auth"
	"restapi/internal/user"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //временно игнорировать миграции в первый раз а потом их добавить
	})
	if err != nil {
		panic(err)
	}
	return db // тесты нужно проводить на база на которых уже была проведена миграция
}

func removeData(db *gorm.DB) {
	// SOFT DELETE получается благодаря колонке deleted at , поэтому и не удаляется сразу
	//нужно сделать uncsope для софт делита чтобы его убрать(полное удаление безвозвратно)
	db.Unscoped().
		Where("email = ?", "sadtest123@d.ru").
		Delete(&user.User{})
}
func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "sadtest123@d.ru",
		Password: "$2a$10$Ijdf1/9t3f/OrAqq6U9/LeTD.NocZJigAPvyLrLbWBSsOf.lX3Mdi",
		Name:     "DANYA",
	})
}
func TestLoginSuccess(t *testing.T) {
	// prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "sadtest123@d.ru",
		Password: "123",
	})

	//query
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatal("Token empty")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	// prepare
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "sadtest123@d.ru",
		Password: "12",
	})

	//query
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}
	removeData(db)

}
