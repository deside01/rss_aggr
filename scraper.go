package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/google/uuid"
)

func startScraper(db *database.Queries, concurrency int, requestDelay time.Duration) {
	log.Printf("\nStarting scraping with %v gorutines and %v delay", concurrency, requestDelay)
	ticker := time.NewTicker(requestDelay)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeeds(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Can't get next feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, v := range feeds {
			wg.Add(1)

			go fetchFeed(wg, db, v)
		}
		wg.Wait()
	}

}

func fetchFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	rss, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Can't parse feed:", err)
		return
	}

	for _, v := range rss.Channel.Items {
		description := sql.NullString{}
		pubDate := sql.NullTime{}

		if v.Description != "" {
			description.String = v.Description
			description.Valid = true
		}
		if v.PublishDate != "" {
			pubDateTime, _ := time.Parse(time.RFC1123, v.PublishDate)
			pubDate.Time = pubDateTime
			pubDate.Valid = true
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       v.Title,
			Link:        v.Link,
			Description: description,
			PublishDate: pubDate,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})

		if err != nil {
			if strings.Contains(err.Error(), "posts_link_key") {
				continue
			}
			log.Println("failed to create post:", err)
		}
	}

	log.Printf("Total: %v\n", len(rss.Channel.Items))
}
