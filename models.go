package main

import (
	"time"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}

func dbUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func dbFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		Name:      feed.Name,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func dbFeedsToFeeds(dbFeeds []database.Feed) (feed []Feed) {
	for _, dbFeed := range dbFeeds {
		feed = append(feed, dbFeedToFeed(dbFeed))
	}

	return feed
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func dbFeedFollowToFeedFollow(feed database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feed.ID,
		UserID:    feed.UserID,
		FeedID:    feed.FeedID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	}
}

func dbFeedsFollowsToFeedsFollows(dbFeeds []database.FeedFollow) (feeds []FeedFollow) {
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, dbFeedFollowToFeedFollow(dbFeed))
	}

	return feeds
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	FeedID      uuid.UUID  `json:"feed_id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Description *string    `json:"description"`
	PublishDate *time.Time `json:"publish_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func dbPostToPost(dbPost database.Post) Post {
	var description *string
	var pubDate *time.Time
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	if dbPost.PublishDate.Valid {
		pubDate = &dbPost.PublishDate.Time
	}

	return Post{
		ID:          dbPost.ID,
		FeedID:      dbPost.FeedID,
		Title:       dbPost.Title,
		Link:        dbPost.Link,
		Description: description,
		PublishDate: pubDate,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
	}
}

func dbPostsToPosts(dbPosts []database.Post) (posts []Post) {
	for _, v := range dbPosts {
		posts = append(posts, dbPostToPost(v))
	}

	return posts
}
