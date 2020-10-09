package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/thinkerou/favicon"
)

const (
	mediaParams = "fields=id,username,caption,media_type,media_url,permalink,thumbnail_url,timestamp"
)

// Client instagram connection representation.
type Client struct {
	accessToken string
	client      *http.Client
	baseURL     string
}

// Entry represents an instagram media post.
type Entry struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Caption      string `json:"caption"`
	MediaType    string `json:"media_type"`
	MediaURL     string `json:"media_url"`
	Permalink    string `json:"permalink"`
	ThumbnailURL string `json:"thumbnail_url"`
	Timestamp    string `json:"timestamp"`
}

// Paging represents endpoints for paging results if above limit (100)
// https://developers.facebook.com/docs/graph-api/using-graph-api?_fb_noscript=1
type Paging struct {
	Cursors struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}
	// If not included, this is the last page of data.
	// Stop paging when the next link no longer appears.
	Next string `json:"next,omitempty"`
	// If not included, this is the first page of data.
	Previous string `json:"previous,omitempty"`
}

// MediaResp is representing JSON received from graph
type MediaResp struct {
	Data   []Entry `json:"data"`
	Paging Paging  `json:"paging,omitempty"`
}

func main() {
}

func init() {
	// Starts a new Gin instance with no middle-ware
	r := gin.New()
	r.Use(favicon.New("./favicon.png"))
	// Define your handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "instagram basic API")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/json", getJSON)

	r.Run() // listen and serve on 0.0.0.0:8080
	// Handle all requests using net/http
	http.Handle("/", r)
}

func getJSON(c *gin.Context) {
	var (
		longToken = os.Getenv("IG_TOKEN")
		client    = NewClient(longToken)
		media     = []Entry{}
		next      = ""
		counter   = 0
		limit     = c.DefaultQuery("limit", "20")
		lim, _    = strconv.Atoi(limit)
	)
	for counter < lim {
		response, err := client.GetMedia(limit, next)
		if err != nil {
			log.Println(err.Error())
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": true, "message": err.Error()})
			return
		}
		paging := &response.Paging
		media = append(media, response.Data...)
		// https://developers.facebook.com/docs/graph-api/using-graph-api?_fb_noscript=1
		// Do not depend on the number of results being fewer than the limit value to indicate
		// that your query reached the end of the list of data, use the absence of next instead
		if next = paging.Next; next == "" {
			// Stop
			counter = lim
		} else {
			counter += len(response.Data)
		}
		log.Printf("Received # media: %d. Counter: %d", len(response.Data), counter)
	}
	c.JSON(http.StatusOK, media)
	return
}

// NewClient creates a new client with a given accessToken and clientSecret.
func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		baseURL:     "https://graph.instagram.com",
		client:      &http.Client{},
	}
}

func (e *Entry) String() string {
	return fmt.Sprintf(
		"%T{ID: %s,Username: %s,Caption: %s,MediaType: %s,MediaURL: %s,Permalink: %s,ThumbnailURL: %s,Timestamp:%s}",
		e,
		e.ID,
		e.Username,
		e.Caption,
		e.MediaType,
		e.MediaURL,
		e.Permalink,
		e.ThumbnailURL,
		e.Timestamp,
	)
}

// Tags returns all identified hashtags in caption.
func (e Entry) Tags() []string {
	var result []string
	arr := strings.Split(e.Caption, " ")
	for _, s := range arr {
		if strings.HasPrefix(s, "#") {
			for _, h := range strings.Split(s, "#") {
				if h == "" {
					continue
				}
				result = append(result, h)
			}
		}
	}
	return result
}

// GetMedia fetches media from a user configured within a Client, returns an array of Entries,
// error if something goes wrong with communication.
func (c *Client) GetMedia(limit string, next string) (*MediaResp, error) {
	var (
		resp MediaResp
		url  string
	)
	if next == "" {
		url = buildURL(c.baseURL, "/me/media", c.accessToken, limit, mediaParams)
	} else {
		url = next
	}
	log.Println(url)
	bytes, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		return nil, errors.Wrapf(err, "unable to Unmarshal json response: %s", string(bytes))
	}
	return &resp, nil
}

func (c *Client) fetch(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch profile")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	if err = resp.Body.Close(); err != nil {
		return nil, errors.Wrap(err, "unable to close body")
	}

	return body, nil
}

func buildURL(base, path, token string, limit string, extraParams ...string) string {
	var params string
	if len(extraParams) > 0 {
		params = fmt.Sprintf("&%s", strings.Join(extraParams, "&"))
	}
	// log.Println(params)
	url := fmt.Sprintf("%s%s?access_token=%s&limit=%s%s", base, path, token, limit, params)
	return url
}
