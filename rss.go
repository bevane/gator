package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/bevane/gator/internal/database"
	"github.com/google/uuid"
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
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description *string `xml:"description"`
	PubDate     *string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := &http.Client{}
	rssFeed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, http.NoBody)
	if err != nil {
		return &rssFeed, err
	}
	req.Header.Add("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return &rssFeed, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rssFeed, err
	}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return &rssFeed, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		unescapedString := html.UnescapeString(*item.Description)
		rssFeed.Channel.Item[i].Description = &unescapedString
	}
	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Fetching feed: %s\n", feedToFetch.Name)
	err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}
	rssFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		pubDateSQL := sql.NullTime{}
		if item.PubDate != nil {
			pubDate, err := parsePubDate(item.PubDate)
			pubDateSQL.Time = pubDate
			pubDateSQL.Valid = true
			if err != nil {
				fmt.Println(err)
				pubDateSQL.Valid = false
			}

		}
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: *item.Description,
				Valid:  item.Description != nil,
			},
			PublishedAt: pubDateSQL,
			FeedID:      feedToFetch.ID,
		})
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}
			fmt.Println(err.Error())
		}
	}
	return nil
}

func parsePubDate(dateString *string) (time.Time, error) {
	var datetime time.Time
	var err error
	if _, err = time.Parse(time.ANSIC, *dateString); err == nil {
		datetime, _ = time.Parse(time.ANSIC, *dateString)

	} else if _, err = time.Parse(time.UnixDate, *dateString); err == nil {
		datetime, _ = time.Parse(time.UnixDate, *dateString)

	} else if _, err = time.Parse(time.RubyDate, *dateString); err == nil {
		datetime, _ = time.Parse(time.RubyDate, *dateString)

	} else if _, err = time.Parse(time.RFC822, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC822, *dateString)

	} else if _, err = time.Parse(time.RFC822Z, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC822Z, *dateString)

	} else if _, err = time.Parse(time.RFC850, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC850, *dateString)

	} else if _, err = time.Parse(time.RFC1123, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC1123, *dateString)

	} else if _, err = time.Parse(time.RFC1123Z, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC1123Z, *dateString)

	} else if _, err = time.Parse(time.RFC3339, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC3339, *dateString)

	} else if _, err = time.Parse(time.RFC3339Nano, *dateString); err == nil {
		datetime, _ = time.Parse(time.RFC3339Nano, *dateString)

	} else if _, err = time.Parse(time.Stamp, *dateString); err == nil {
		datetime, _ = time.Parse(time.Stamp, *dateString)

	} else if _, err = time.Parse(time.StampMilli, *dateString); err == nil {
		datetime, _ = time.Parse(time.StampMilli, *dateString)

	} else if _, err = time.Parse(time.StampMicro, *dateString); err == nil {
		datetime, _ = time.Parse(time.StampMicro, *dateString)

	} else if _, err = time.Parse(time.StampNano, *dateString); err == nil {
		datetime, _ = time.Parse(time.StampNano, *dateString)

	} else if _, err = time.Parse(time.DateTime, *dateString); err == nil {
		datetime, _ = time.Parse(time.DateTime, *dateString)

	} else if _, err = time.Parse(time.DateOnly, *dateString); err == nil {
		datetime, _ = time.Parse(time.DateOnly, *dateString)

	} else if _, err = time.Parse(time.TimeOnly, *dateString); err == nil {
		datetime, _ = time.Parse(time.TimeOnly, *dateString)
	} else {
		return time.Time{}, fmt.Errorf("Unable to parse datetime")
	}
	return datetime, nil
}
