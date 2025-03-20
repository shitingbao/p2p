package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/pion/stun"
)

func WithFlagHost(flagHost string) Option {
	return func(o *option) {
		o.FlagHost = flagHost
	}
}

func WithStunRaw(raw string) Option {
	return func(o *option) {
		o.StunRaw = raw
	}
}

func NewClient(opts ...Option) *Client {
	o := &option{}

	for _, opt := range opts {
		opt(o)
	}

	if o.StunRaw == "" {
		o.StunRaw = DefaultUri
	}

	cli := &Client{
		StunRaw: o.StunRaw,
	}
	return cli
}

func (c *Client) GetIP() (string, error) {
	uri, err := stun.ParseURI(c.StunRaw)
	if err != nil {
		return "", err
	}

	cli, err := stun.DialURI(uri, &stun.DialConfig{})
	if err != nil {
		return "", err
	}

	internetIP := ""
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	if err := cli.Do(message, func(res stun.Event) {
		if res.Error != nil {
			panic(res.Error)
		}
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			panic(err)
		}

		internetIP = xorAddr.String()
		// log.Println("your IP is", xorAddr.String())
	}); err != nil {
		return "", err
	}

	return internetIP, nil
}

func (c *Client) SetClientId(id string) {
	c.clientId = id
}

// 传入基本路由和参数，反馈结果
func (c *Client) sendPost(rou string, v any) ([]byte, error) {
	jsonBody, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", path.Join(c.FlagHost, rou), bytes.NewBuffer(jsonBody))
	if err != nil {
		return []byte{}, err
	}

	// req.Header.Add("Authorization", "APPCODE "+appcode)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
