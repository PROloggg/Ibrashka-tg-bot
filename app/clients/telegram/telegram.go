package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"app/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	sendVoice         = "sendVoice"
	getFile           = "getFile"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)

	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendVoice(chatId int, voiceFilePath string) error {
	file, err := os.Open(voiceFilePath)
	if err != nil {
		return e.Wrap("can't open voice file", err)
	}
	defer func() { _ = file.Close() }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, sendVoice),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	// Создаем multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем voice файл в multipart
	part, err := writer.CreateFormFile("voice", path.Base(voiceFilePath))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	err = writer.WriteField("chat_id", strconv.Itoa(chatId))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Устанавливаем тело запроса и заголовок
	req.Body = io.NopCloser(body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send voice, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) GetFileURL(fileID string) (string, error) {
	q := url.Values{}
	q.Add("file_id", fileID)

	data, err := c.doRequest(getFile, q)
	if err != nil {
		return "", err
	}

	var res struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	if !res.Ok {
		return "", fmt.Errorf("failed to get file URL")
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("/file/", c.basePath, res.Result.FilePath),
	}

	return u.String(), nil
}
