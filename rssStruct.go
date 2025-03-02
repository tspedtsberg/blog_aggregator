package main

import (
	"net/http"
	"io"
	"encoding/xml"
	"context"
	"fmt"
	"html"
)


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
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
		 return nil, fmt.Errorf("erroring creating the request: %w", err)
	}
	client := http.Client{}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erroring making the request: %w", err)
    }
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("erroring reading the data: %w", err)
    }

	var rssfeed RSSFeed
	if err := xml.Unmarshal(data, &rssfeed); err != nil{
		return nil, fmt.Errorf("erroring unmarshalling the data: %w", err)
    }
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString(rssfeed.Channel.Description)
	for i, item := range rssfeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssfeed.Channel.Item[i] = item
	}

	return &rssfeed, nil
}