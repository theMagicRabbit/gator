package feed

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/theMagicRabbit/gator/internal/database"
	"github.com/theMagicRabbit/gator/internal/state"
)

type RSSFeed struct {
	Channel struct {
		Title		string		`xml:"title"`
		Link		string  	`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem 	`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	feed := new(RSSFeed)
	err = xml.Unmarshal(body, feed)
	if err != nil {
		return nil, err
	}
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i, item := range feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i] = item
	}
	return feed, nil
}

func ScrapeFeeds(s *state.State) error {
	next, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feed, err := FetchFeed(context.Background(), next.Url)
	if err != nil {
		return err
	}
	utcTimestamp := time.Now().UTC()
	params := database.MarkFeedFetchedParams{
		ID: next.ID,
		UpdatedAt: utcTimestamp,
	}
	_, err = s.Db.MarkFeedFetched(context.Background(), params)
	if err != nil { 
		return err
	}
	for _, item := range feed.Channel.Item {
		err = savePost(s, next, item)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				// If the URl already exists, we can safely skip the error
				if pqErr.Code == "23505" {
					continue
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func savePost(s *state.State, feed database.Feed, item RSSItem) error {
	utcNow := time.Now().UTC()
	itemTitle := sql.NullString{}
	itemUrl := sql.NullString{}
	itemDescription := sql.NullString{}
	itemPubDate := sql.NullTime{}
	if item.Title != "" {
		itemTitle.String = item.Title
		itemTitle.Valid = true
	}
	if item.Link != "" {
		itemUrl.String = item.Link
		itemUrl.Valid = true
	}
	if item.Description != "" {
		itemDescription.String = item.Description
		itemDescription.Valid = true
	}
	pubDate, err := time.Parse("Mon, 03 Jan 2006 13:04:05 +0000", item.PubDate)
	if err == nil {
		if !pubDate.IsZero() {
			itemPubDate.Time = pubDate
			itemPubDate.Valid = true
		}
	}

	params := database.CreatePostParams{
		ID: uuid.New(),
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
		Title: itemTitle,
		Url: itemUrl,
		Description: itemDescription,
		PublishedAt: itemPubDate,
		FeedID: feed.ID,
	}
	post, err := s.Db.CreatePost(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(post)
	return nil
}
