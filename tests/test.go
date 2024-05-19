package main

import (
	"database/sql"

	"testing"

	"github.com/damirbeybitov/todo_project/internal/user/repository"
	_ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
    // Connect to the test database
    db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/to_do")
    if err != nil {
        t.Fatalf("Error connecting to MySQL: %v", err)
    }
    return db
}

func TestCheckUserInDB(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Begin a transaction
    tx, err := db.Begin()
    if err != nil {
        t.Fatalf("Error beginning transaction: %v", err)
    }
    defer tx.Rollback()

    r := &repository.Repository{DB: db}

    // Test case 1: user does not exist
    err = r.CheckUserInDB(tx, "test_user", "test_email@example.com")
    if err != nil {
        t.Errorf("CheckUserInDB returned an error: %v", err)
    }

    // Test case 2: user already exists
    err = r.CheckUserInDB(tx, "existing_user", "existing_email@example.com")
    if err == nil {
        t.Error("Expected CheckUserInDB to return an error for existing user, but it didn't")
    }
}

func TestAddUserToDB(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        t.Fatalf("Error beginning transaction: %v", err)
    }
    defer tx.Rollback()

    r := &repository.Repository{DB: db}

    // Test case: Add user to DB
    id, err := r.AddUserToDB(tx, "new_user", "new_email@example.com", "password")
    if err != nil {
        t.Errorf("AddUserToDB returned an error: %v", err)
    }

    if id <= 0 {
        t.Error("Expected AddUserToDB to return a valid ID, but it returned <= 0")
    }
}

func TestCheckPassword(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    r := &repository.Repository{DB: db}

    // Add a user for testing
    tx, err := db.Begin()
    if err != nil {
        t.Fatalf("Error beginning transaction: %v", err)
    }
    defer tx.Rollback()

    _, err = r.AddUserToDB(tx, "existing_user", "existing_email@example.com", "correct_password")
    if err != nil {
        t.Fatalf("Error adding user to test database: %v", err)
    }

    // Test case 1: correct password
    err = r.CheckPassword("existing_user", "correct_password")
    if err != nil {
        t.Errorf("CheckPassword returned an unexpected error: %v", err)
    }

    // Test case 2: incorrect password
    err = r.CheckPassword("existing_user", "incorrect_password")
    if err == nil {
        t.Error("Expected CheckPassword to return an error for incorrect password, but it didn't")
    }
}

func TestDeleteuserFromDB(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    r := &repository.Repository{DB: db}

    // Add a user for testing
    tx, err := db.Begin()
    if err != nil {
        t.Fatalf("Error beginning transaction: %v", err)
    }
    defer tx.Rollback()

    _, err = r.AddUserToDB(tx, "user_to_delete", "email@example.com", "password")
    if err != nil {
        t.Fatalf("Error adding user to test database: %v", err)
    }

    // Test case: Delete user from DB
    err = r.DeleteuserFromDB("user_to_delete")
    if err != nil {
        t.Errorf("DeleteuserFromDB returned an error: %v", err)
    }
}