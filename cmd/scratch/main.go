package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AfshinJalili/gonod/internal/config"
	"github.com/AfshinJalili/gonod/internal/domain"
	"github.com/AfshinJalili/gonod/internal/platform"
	"github.com/AfshinJalili/gonod/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	db, err := platform.SetupDatabase(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	plainPassword := "password123"
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	newUser := &domain.User{
		Email:    "test@example.com",
		Password: string(hashedBytes),
	}

	fmt.Println("Inserting user...")
	err = repo.CreateUser(ctx, newUser)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	fmt.Printf("Sucsess! Generated UUID: %s\n", newUser.ID)

	fmt.Println("Fetching user by email...")
	fetchedUser, err := repo.GetUserByEmail(ctx, newUser.Email)
	if err != nil {
		log.Fatal("Failed to fetch user:", err)
	}

	fmt.Printf("Found User: %s | Created At: %v\n", fetchedUser.Email, fetchedUser.CreatedAt)

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(plainPassword))
	if err == nil {
		fmt.Println("Password verified successfully")
	}
}
