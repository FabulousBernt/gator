package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/FabulousBernt/gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) {
	// 1. Get the next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("couldn't get next feed to fetch:", err)
		return
	}

	// 2. Mark it as fetched
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("couldn't mark feed as fetched:", err)
		return
	}

	// 3. Fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("couldn't fetch feed:", err)
		return
	}

	// 4. Save each post
	for _, item := range rssFeed.Channel.Item {
		// Parse published_at - try common formats
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		} else if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		}

		// Handle nullable description
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Ignore duplicate URL errors, log everything else
			if err.Error() != "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				log.Printf("couldn't save post %s: %v", item.Title, err)
			}
		}
	}
	log.Printf("collected %d posts from %s", len(rssFeed.Channel.Item), feed.Name)
}
