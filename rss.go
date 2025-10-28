package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}


	feedStruct := RSSFeed{}

	if err := xml.Unmarshal(data, &feedStruct); err != nil {
		return &RSSFeed{}, err
	}

	feedStruct.Channel.Title = html.UnescapeString(feedStruct.Channel.Title)
	feedStruct.Channel.Description = html.UnescapeString(feedStruct.Channel.Description)
	for i, item := range feedStruct.Channel.Items {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feedStruct.Channel.Items[i] = item
	}


	return &feedStruct, nil
}

func scrapeFeeds(s *State) error {
	feedData, err := s.db.GetNextFeedToFetch(context.Background())	
	if err != nil {
		return fmt.Errorf("unable to fetch next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feedData.ID)
	if err != nil {
		return fmt.Errorf("unable to mark feed fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), feedData.Url)
	if err != nil {
		return fmt.Errorf("unable to list users: %w", err)
	}

	fmt.Println()
	fmt.Println(feed.Channel.Title, "posts:")
	for _, item := range feed.Channel.Items {
		fmt.Println(item.Title)
	}
	fmt.Println()
	
	return nil
}
