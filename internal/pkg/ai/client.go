package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Client AI服务客户端
type Client struct {
    baseURL     string
    callbackURL string
    httpClient  *http.Client
}

// NewClient 创建AI服务客户端
func NewClient(baseURL, callbackURL string) *Client {
    return &Client{
        baseURL:     baseURL,
        callbackURL: callbackURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// PredictRequest AI预测请求
type PredictRequest struct {
    ImageURL    string `json:"image_url"`
    TaskID      string `json:"task_id"`
    CallbackURL string `json:"callback_url"`
}

// RequestPredict 请求AI预测
func (c *Client) RequestPredict(imageURL, taskID string) error {
    reqBody := PredictRequest{
        ImageURL:    imageURL,
        TaskID:      taskID,
        CallbackURL: c.callbackURL,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return fmt.Errorf("marshal request body error: %v", err)
    }

    req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/predict", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("create request error: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("do request error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
} 