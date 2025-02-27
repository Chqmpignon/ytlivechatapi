package ytlivechatapi

/*
Go implementation of the YouTube Live Chat API

Expects an authorized *http.Client, golang.org/x/oauth2 is a good package to use.

Send "Hello world!" to the first default stream chat.
    c := ytlivechatapi.NewClient("CLIENTID", "CLIENTSECRET", "APIKEY")
    if response, err := c.ListLiveBroadcasts("default=true"); err == nil {
      c.InsertLiveChatMessage(ytlivechatapi.NewLiveChatMessage(response.Items[0].Snippet.LiveChatId, "Hello world!"))
    }

Polls on the first default stream chat and bans everyone that sends a message for 10 seconds.
    c := ytlivechatapi.NewClient("CLIENTID", "CLIENTSECRET", "APIKEY")
    if response, err := c.ListLiveBroadcasts("default=true"); err == nil {
      liveChatId := response.Items[0].Snippet.LiveChatId
      nextPageToken := ""
      for {
        if response, err := c.ListLiveChatMessages(liveChatId, nextPageToken); err == nil {
          nextPageToken = response.NextPageToken

          for _, message := range response.Items {
            c.InsertLiveChatBan(ytlivechatapi.NewLiveChatBan(liveChatId, message.AuthorDetails.ChannelId, 10))
          }

          time.Sleep(time.Duration(response.PollingIntervalMillis) * time.Millisecond)
        }
      }
    }
*/