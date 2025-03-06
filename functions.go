package main

import (
	"fmt"
	"time"
	"log"
	"context"
	"Aggregator/internal/database"
	"github.com/google/uuid"
	"strings"
	"strconv"
	"database/sql"
)

func browse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		specificlimit, err := strconv.Atoi(cmd.Args[0])
		if err == nil {
			limit = specificlimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("Author: %s\n", post.FeedName)
		fmt.Printf("---%s--- \n", post.Title)
		fmt.Printf("   %v\n", post.Description.String)
		fmt.Printf("===============================\n")
	}

	return nil

}


func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url to unfollow>", cmd.name)
	}

	urlOfFeed := cmd.Args[0]
	err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		Url: urlOfFeed,
	})
	if err != nil {
		fmt.Printf("error unfollowing: %s", err)
	}

	return nil
}


func listfeedsfollow(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feedfollows: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil
}


func follow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.name)
	}
	
	urlOfFeed := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), urlOfFeed)
	if err != nil {
		return fmt.Errorf("erroring fetching feed by url: %w", err)
	}

	feedfollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating the feedfollow: %w", err)
   }
    
    //fmt.Printf("test: %+v\n", feedfollow)
	fmt.Printf("Name of User: %s\n", feedfollow.UserName)
	fmt.Printf("Name of feed: %s\n", feedfollow.FeedName)


	return nil
}


func listfeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("erroring fecthing feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("erroring fetching decoding UserId: %w", err)
		}
	
		printFeed(feed, user)
	}
	
	return nil
}


func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Usage: %s <nameOfFeed> <url>", cmd.name)
	}

	nameOfFeed := cmd.Args[0]
	url := cmd.Args[1]
	//fmt.Println(nameOfFeed)
	//fmt.Println(url)
			
	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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
	
	followcmd := command{
		name: "follow",
		Args: []string{url},
	}

	err = follow(s, followcmd, user)
	if err != nil {
		fmt.Printf("Warning: feed was added but could not be followed: %v\n", err)
	}

	printFeed(newFeed, user)

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("*ID: %s\n", 		feed.ID)
	fmt.Printf("*CreatedAt: %s\n", 	feed.CreatedAt)
	fmt.Printf("*UpdatedAt: %s\n", 	feed.UpdatedAt)
	fmt.Printf("*Name: %s\n", 		feed.Name)
	fmt.Printf("*URL: %s\n", 		feed.Url)
	fmt.Printf("*UserId: %s\n", 	feed.UserID)
	fmt.Printf("*UserID decoded: %s\n",user.Name )
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Couldn't get next feed: ", err)
		return
	}

	log.Printf("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark the feed %s fetched: %v", feed.Name, err)
		return
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch the feed: %v", err)
		return
	}
	for _, item := range feedData.Channel.Item {
		publishedat := sql.NullTime{}
		t, err := time.Parse("Mon, 2 jan 2006 15:05:05", item.PubDate)
		if err == nil {
			publishedat = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}
		
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: 			uuid.New(),
			CreatedAt: 		time.Now().UTC(),
			UpdatedAt: 		time.Now().UTC(),
			Title: 			item.Title,
			Url: 			feed.Url,
			Description:	sql.NullString{
				String:	item.Description,
				Valid: true,
			},
			PublishedAt:	publishedat,
			FeedID:			feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
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