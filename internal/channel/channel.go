package channel

import (
	"context"
	"errors"
)

type Channel struct {
	Publisher  chan string
	Subscriber *Client
}

type Client struct {
	Name    string
	Channel chan string
}

func (c *Channel) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.Publisher:
			if c.Subscriber != nil && c.Subscriber.Channel != nil {
				c.Subscriber.Channel <- data
			}
		}
	}
}

func (c *Channel) BindSubscriber(client *Client) (*Client, error) {
	if c.Subscriber == nil {
		c.Subscriber = client
		return c.Subscriber, nil
	}
	if c.Subscriber != nil && c.Subscriber.Name != client.Name {
		return nil, errors.New("频道已有其他绑定")
	}
	return c.Subscriber, nil
}

var channelMap = make(map[string]*Channel)

func Set(name string, channel *Channel) {
	channelMap[name] = channel
}

func Get(name string) *Channel {
	if channel, ok := channelMap[name]; ok {
		return channel
	}
	return nil
}
