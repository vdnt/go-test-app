package main

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLGreetingRoundTrip(t *testing.T) {
	db := openTestMySQL(t)
	resetGreetings(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	message := Greet("mysql")
	result, err := db.ExecContext(ctx, `INSERT INTO greetings (name, message) VALUES (?, ?)`, "mysql", message)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	var gotName, gotMessage string
	if err := db.QueryRowContext(ctx, `SELECT name, message FROM greetings WHERE id = ?`, id).Scan(&gotName, &gotMessage); err != nil {
		t.Fatal(err)
	}
	if gotName != "mysql" || gotMessage != "hello, mysql!" {
		t.Fatalf("stored greeting = (%q, %q)", gotName, gotMessage)
	}
}

func TestMySQLTransactionRollback(t *testing.T) {
	db := openTestMySQL(t)
	resetGreetings(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO greetings (name, message) VALUES (?, ?)`, "rollback", Greet("rollback")); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM greetings WHERE name = ?`, "rollback").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rolled-back rows = %d, want 0", count)
	}
}

func openTestMySQL(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		t.Skip("MYSQL_DSN is not set")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	deadline := time.Now().Add(20 * time.Second)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			return db
		}
		if time.Now().After(deadline) {
			t.Fatalf("mysql did not become ready: %v", err)
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func resetGreetings(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS greetings (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		message VARCHAR(255) NOT NULL
	)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `TRUNCATE TABLE greetings`); err != nil {
		t.Fatal(err)
	}
}
