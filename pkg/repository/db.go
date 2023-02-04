package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tokenTable = "authorization_token"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func NewSQLiteDB() (*sql.DB, error) {
	conf := ReadConfig()
	db, err := sql.Open("sqlite3", conf.DBName)
	if err != nil {
		return nil, err
	}
	// проверка подкл. бд
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = CreateTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ReadConfig() *Config {
	cfg := new(Config)
	jsonfile, err := os.ReadFile("pkg/repository/config/config.json")
	if err != nil {
		log.Fatalf("error while reading josn file: %s", err)
	}
	err = json.Unmarshal(jsonfile, &cfg)
	if err != nil {
		log.Fatalf("error with unmarshal file: %s", err)
	}
	return cfg
}

func CreateTables(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		Email TEXT NOT NULL UNIQUE, 
		Username TEXT NOT NULL UNIQUE, 
		Password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS authorization_token (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		UserId INTEGER NOT NULL, 
		auth_token VARCHAR(255), 
		expires_at  DATETIME NOT NULL, 
		FOREIGN KEY(UserId) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS post (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		UserID INTEGER NOT NULL, 
		title TEXT NOT NULL, 
		text TEXT NOT NULL, 
		category TEXT NOT NULL, 
		createdAt TEXT NOT NULL, 
		author TEXT NOT NULL, 
		like_vote TEXT DEFAULT '0', 
		dislike TEXT DEFAULT '0',
		image TEXT,
		FOREIGN KEY(UserID) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS comment (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		UserId INTEGER NOT NULL, 
		PostId INTEGER NOT NULL, 
		comment TEXT NOT NULL, 
		createdAt TEXT NOT NULL, 
		author TEXT NOT NULL,
		like_vote TEXT NOT NULL, 
		dislike TEXT NOT NULL,
		FOREIGN KEY(UserId) REFERENCES users(id), 
		FOREIGN KEY(PostId) REFERENCES post(id)
	);
	
	CREATE TABLE IF NOT EXISTS evaluate (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		UserId INTEGER NOT NULL, 
		PostId INTEGER NOT NULL, 
		Vote INTEGER NOT NULL, 
		FOREIGN KEY(UserId) REFERENCES users(id),
		FOREIGN KEY(PostId) REFERENCES post(id) 
	);

	CREATE TABLE IF NOT EXISTS evaluateComment (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		userId INTEGER NOT NULL, 
		commentId INTEGER NOT NULL, 
		vote INTEGER NOT NULL, 
		FOREIGN KEY(UserId) REFERENCES users(id),
		FOREIGN KEY(commentId) REFERENCES comment(id) 
	);

	
	`

	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
