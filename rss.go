package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RssFeed struct {
	Channel struct {
		Items []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PublishDate string `xml:"pubDate"`
}

func urlToFeed(feedUrl string) (rss RssFeed, err error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(feedUrl)
	if err != nil {
		return rss, fmt.Errorf("can't make request: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return rss, fmt.Errorf("can't read body: %v", err)
	}

	xml.Unmarshal(data, &rss)

	return rss, nil
}
