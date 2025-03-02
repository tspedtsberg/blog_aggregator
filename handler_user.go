package main

import (
	"fmt"
	"time"
	"context"
	"Aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerfeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Usage: %s <nameOfFeed> <url>", cmd.name)
	}

	nameOfFeed := cmd.Args[0]
	url := cmd.Args[1]
	//fmt.Println(nameOfFeed)
	//fmt.Println(url)
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
		
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: nameOfFeed,
		Url: url,
		UserID: user.ID,

	})
	if err != nil {
		 return fmt.Errorf("error creating the feed: %w", err)
	}
	
	fmt.Printf("Feed: %+v\n", feed)

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}


func listUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error receiving users: %w", err)
	}
	currentuser := s.cfg.CurrentUserName

	for _, user := range users {
		if user == currentuser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		} 
	}
	return nil
}

func reset(s *state, cmd command) error {
	err := s.db.CleanTable(context.Background())
	if err != nil {
		return fmt.Errorf("error cleaning the tabel: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("error creating a user: %w", err)
	}
	
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	fmt.Printf("current user: %s", s.cfg.CurrentUserName)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}