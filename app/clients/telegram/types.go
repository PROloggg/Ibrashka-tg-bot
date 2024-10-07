package telegram

type Update struct {
    UpdateID int     `json:"update_id"`
    Message  *Message `json:"message"`
}

type Message struct {
    MessageID int    `json:"message_id"`
    From      User   `json:"from"`
    Chat      Chat   `json:"chat"`
    Date      int64  `json:"date"`
    Text      string `json:"text"`
    Caption string    `json:"caption"`
    Photo     []Photo `json:"photo"`
}

type Photo struct {
    FileID string `json:"file_id"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

type User struct {
    ID            int    `json:"id"`
    IsBot         bool   `json:"is_bot"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    Username      string `json:"username"`
    LanguageCode  string `json:"language_code"`
}

type Chat struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Username  string `json:"username"`
    Type      string `json:"type"`
}

type UpdatesResponse struct {
    Ok     bool     `json:"ok"`
    Result []Update `json:"result"`
}
