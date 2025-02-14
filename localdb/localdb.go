package localdb

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/vrypan/farma/config"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Instance *sql.DB
	dbPath   string
)

func init() {
	dbDir, err := config.ConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dbPath = filepath.Join(dbDir, "local.db")
}

func GetDbPath() string {
	return dbPath
}

func IsOpen() bool {
	return Instance != nil
}

func Open() error {
	if !IsOpen() {
		var err error
		Instance, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return fmt.Errorf("failed to open db: %v", err)
		}

		if _, err := Instance.Exec("PRAGMA journal_mode = WAL;"); err != nil {
			log.Fatalf("Failed to enable WAL mode: %v", err)
		}
	}
	log.Printf("DB is open. Path=%s", dbPath)
	return nil
}

func Close() error {
	if Instance != nil {
		err := Instance.Close()
		Instance = nil
		return err
	}
	return nil
}

func AssertOpen() {
	if Instance == nil {
		log.Fatalf("Database is not open.")
	}
}

func CreateTables() error {
	AssertOpen()
	sql := `
		CREATE TABLE IF NOT EXISTS frames (
			id 	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			desc TEXT DEFAULT "",
			domain TEXT DEFAULT "",
			endpoint TEXT UNIQUE NOT NULL
		);
		CREATE INDEX idx_frames_endpoint on frames (
			endpoint
		);
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS users_frames (
			id TEXT NOT NULL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			frame_id INTEGER NOT NULL,
			app_id INTEGER NOT NULL,
			status INTEGER NOT NULL CHECK (status IN (0, 1, 2)),
			url TEXT,
			token TEXT,
			ctime TEXT DEFAULT CURRENT_TIMESTAMP,
			mtime TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_users_frames_ids ON users_frames (
			user_id, frame_id
		);
		CREATE INDEX idx_users_frames_frame ON users_frames (
			frame_id
		);
		CREATE INDEX idx_users_frames_user ON users_frames (
			user_id
		);
		CREATE INDEX idx_users_frames_token ON users_frames (
			token
		);

		CREATE TABLE IF NOT EXISTS user_history (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			frame_id INTEGER NOT NULL,
			app_id INTEGER NOT NULL,
			event_type TEXT,
			event_details TEXT,
			ctime TEXT DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err := Instance.Exec(sql)
	return err
}

func LogUserHistory(
	UserFid int,
	FrameId int,
	AppFid int,
	EventType string,
	EventDetails string,
) error {
	AssertOpen()

	tx, txErr := Instance.Begin()
	if txErr != nil {
		return fmt.Errorf("Error starting transaction: %v", txErr)
	}

	stmt, prepareErr := tx.Prepare(`
			INSERT INTO user_history(user_id, frame_id, app_id, event_type, event_details)
			VALUES (?, ?, ?, ?, ?)
			`)

	if prepareErr != nil {
		return fmt.Errorf("Error preparing statement: %v", prepareErr)
	}

	_, execErr := stmt.Exec(UserFid, FrameId, AppFid, EventType, EventDetails)

	if execErr != nil {
		return fmt.Errorf("Error while updating user_history: %v", execErr)
	}

	err := tx.Commit()
	stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func UpdateFrameStatus(
	UserFid int,
	FrameId int,
	AppFid int,
	IsActive bool,
	AppUrl string,
	Token string,
) error {
	AssertOpen()

	tx, txErr := Instance.Begin()
	if txErr != nil {
		return fmt.Errorf("Error starting transaction: %v", txErr)
	}

	stmt, prepareErr := tx.Prepare(`
		INSERT INTO users_frames(id, user_id, frame_id, app_id, status, url, token)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE
			SET status = excluded.status,
			url = excluded.url,
			token = excluded.token,
			mtime = CURRENT_TIMESTAMP
	`)

	if prepareErr != nil {
		return fmt.Errorf("Error preparing statement: %v", prepareErr)
	}

	_, execErr := stmt.Exec(
		fmt.Sprintf("%012d-%012d-%012d", UserFid, FrameId, AppFid),
		UserFid, FrameId, AppFid, IsActive, AppUrl, Token)

	if execErr != nil {
		return fmt.Errorf("Error while updating users_frames: %v", execErr)
	}

	err := tx.Commit()
	stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func UpdateInvalidTokens(tokens []string) error {
	AssertOpen()
	params := strings.Join(strings.Split(strings.Repeat("?", len(tokens)), ""), ",")

	tx, txErr := Instance.Begin()
	if txErr != nil {
		return fmt.Errorf("Error starting transaction: %v", txErr)
	}

	stmt, prepareErr := tx.Prepare(
		fmt.Sprintf("UPDATE users_frames SET status=0, url='', token='' WHERE token IN (%s)", params),
	)

	if prepareErr != nil {
		return fmt.Errorf("UpdateInvalidTokens - Error preparing statement: %v", prepareErr)
	}

	args := make([]interface{}, len(tokens))
	for i, token := range tokens {
		args[i] = token
	}
	_, execErr := stmt.Exec(args...)

	if execErr != nil {
		return fmt.Errorf("UpdateInvalidTokens - Error while updating users_frames: %v", execErr)
	}

	err := tx.Commit()
	stmt.Close()
	if err != nil {
		return err
	}

	return nil
}
