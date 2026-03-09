package main

import (
	"context"
	"fmt"
	"log"
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

	// 4. Print the titles
	fmt.Printf("Fetching feed: %s\n", feed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf(" - %s\n", item.Title)
	}
}
