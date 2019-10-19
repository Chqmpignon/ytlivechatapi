package ytlivechatapi

import (
  youtube "google.golang.org/api/youtube/v3"
  "golang.org/x/oauth2/google"
  "golang.org/x/oauth2"
  
  "bytes"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strings"
  "time"
  "log"
  "io"
)

type Client struct {
  http			*http.Client
  clientId		string
  clientSecret	string
  apiKey		string
}

func NewClient(clientId, clientSecret, apiKey string) *Client {
  key := ""
  
  go func() {
    handler := func(w http.ResponseWriter, r *http.Request) {
      queries := strings.Split(r.URL.RawQuery, "&")
      for _, element := range queries {
        betweenEquals := strings.Split(element, "=")
        if betweenEquals[0] == "code" {
          fmt.Fprintf(w, betweenEquals[1])
          key = betweenEquals[1]
        }
      }
    }
    
    http.HandleFunc("/oauth", handler)
    http.ListenAndServe(":80", nil)
  }()
  
  conf := &oauth2.Config{
    ClientID:     clientId,
    ClientSecret: clientSecret,
    RedirectURL:  "http://localhost/oauth",
    Scopes: []string{
		youtube.YoutubeScope,
    },
    Endpoint: google.Endpoint,
  }
  
  url := conf.AuthCodeURL("state")
  fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
  
  for {
    time.Sleep(1 * time.Second)
    if key != "" {
      break
	}
  }
  
  tok, err := conf.Exchange(oauth2.NoContext, key)
  if err != nil {
    log.Fatal(err)
  }
  client := conf.Client(oauth2.NoContext, tok)
  
  return &Client{client, clientId, clientSecret, apiKey}
}

func (c *Client) makeRequest(typeRequest, url string) (resp *http.Response, err error) {
  req, err := http.NewRequest(typeRequest, url + "&key=" + c.apiKey, nil)
  if err != nil {
    return nil, err
  }
  
  req.Header.Add("Authorization", "Bearer " + c.clientSecret)
  return c.http.Do(req)
}

func (c *Client) post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
  req, err := http.NewRequest("POST", url + "&key=" + c.apiKey, body)
  if err != nil {
    return nil, err
  }
  
  req.Header.Set("Content-Type", contentType)
  req.Header.Add("Authorization", "Bearer " + c.clientSecret)
  return c.http.Do(req)
}



func (c *Client) delete(url string) (resp *http.Response, err error) {
  return c.makeRequest("DELETE", url)
}

func (c *Client) ListLiveBroadcasts(params string) (*LiveBroadcastListResponse, error) {
  resp, err := c.makeRequest("GET", "https://www.googleapis.com/youtube/v3/liveBroadcasts?part=id,snippet,status,contentDetails&" + params)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  liveBroadcastListResponse := &LiveBroadcastListResponse{}
  err = json.Unmarshal(body, liveBroadcastListResponse)
  if err != nil {
    return nil, err
  }

  if liveBroadcastListResponse.Error != nil {
    return nil, liveBroadcastListResponse.Error.NewError("getting broadcasts")
  }

  return liveBroadcastListResponse, nil
}

func (c *Client) ListLiveChatMessages(liveChatId string, pageToken string) (*LiveChatMessageListResponse, error) {
  pageTokenString := ""
  if pageToken != "" {
    pageTokenString = "&pageToken=" + pageToken
  }

  resp, err := c.makeRequest("GET", "https://www.googleapis.com/youtube/v3/liveChat/messages?maxResults=2000&part=id,snippet,authorDetails&liveChatId=" + liveChatId + pageTokenString)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  liveChatMessageListResponse := &LiveChatMessageListResponse{}
  err = json.Unmarshal(body, liveChatMessageListResponse)
  if err != nil {
    return nil, err
  }

  return liveChatMessageListResponse, nil
}

func (c *Client) InsertLiveChatMessage(liveChatMessage *LiveChatMessage) error {
  jsonString, err := json.Marshal(liveChatMessage)
  if err != nil {
    return err
  }

  resp, err := c.post("https://www.googleapis.com/youtube/v3/liveChat/messages?part=snippet", "application/json", bytes.NewBuffer(jsonString))
  if err != nil {
    return err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  liveChatMessage = &LiveChatMessage{}
  err = json.Unmarshal(body, liveChatMessage)
  if err != nil {
    return err
  }

  if liveChatMessage.Error != nil {
    return liveChatMessage.Error.NewError("inserting LiveChatMessage")
  }

  return nil
}

func (c *Client) DeleteLiveChatMessage(liveChatMessage *LiveChatMessage) error {
  resp, err := c.delete("https://www.googleapis.com/youtube/v3/liveChat/messages?id=" + liveChatMessage.Id)
  if err != nil {
    return err
  }
  return resp.Body.Close()
}

func (c *Client) InsertLiveChatBan(liveChatBan *LiveChatBan) error {
  jsonString, err := json.Marshal(liveChatBan)
  if err != nil {
    return err
  }

  resp, err := c.post("https://www.googleapis.com/youtube/v3/liveChatBans?part=snippet", "application/json", bytes.NewBuffer(jsonString))
  if err != nil {
    return err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  liveChatBan = &LiveChatBan{}
  err = json.Unmarshal(body, liveChatBan)
  if err != nil {
    return err
  }

  if liveChatBan.Error != nil {
    return liveChatBan.Error.NewError("inserting LiveChatBan")
  }

  return nil
}

func (c *Client) DeleteLiveChatBan(liveChatBan *LiveChatBan) error {
  resp, err := c.delete("https://www.googleapis.com/youtube/v3/liveChatBans?id=" + liveChatBan.Id)
  if err != nil {
    return err
  }
  return resp.Body.Close()
}

func (c *Client) ListLiveChatModerators(liveChatId string, pageToken string) (*LiveChatModeratorListResponse, error) {
  pageTokenString := ""
  if pageToken != "" {
    pageTokenString = "&pageToken=" + pageToken
  }

  resp, err := c.makeRequest("GET", "https://www.googleapis.com/youtube/v3/liveChatModerators?maxResults=50&part=id,snippet&liveChatId=" + liveChatId + pageTokenString)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  liveChatModeratorListResponse := &LiveChatModeratorListResponse{}
  err = json.Unmarshal(body, liveChatModeratorListResponse)
  if err != nil {
    return nil, err
  }

  return liveChatModeratorListResponse, nil
}

func (c *Client) InsertLiveChatModerator(liveChatModerator *LiveChatModerator) error {
  jsonString, err := json.Marshal(liveChatModerator)
  if err != nil {
    return err
  }

  resp, err := c.post("https://www.googleapis.com/youtube/v3/liveChatModerators?part=snippet", "application/json", bytes.NewBuffer(jsonString))
  if err != nil {
    return err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  liveChatModerator = &LiveChatModerator{}
  err = json.Unmarshal(body, liveChatModerator)
  if err != nil {
    return err
  }

  if liveChatModerator.Error != nil {
    return liveChatModerator.Error.NewError("inserting LiveChatModerator")
  }

  return nil
}

func (c *Client) DeleteLiveChatModerator(liveChatModerator *LiveChatModerator) error {
  resp, err := c.delete("https://www.googleapis.com/youtube/v3/liveChatModerators?id=" + liveChatModerator.Id)
  if err != nil {
    return err
  }
  return resp.Body.Close()
}