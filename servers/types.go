package servers

import "time"

type ServerConfig struct {
	Nickname      string    `json:"nickname"`
	WatchtowerURL string    `json:"watchtower_url"`
	Token         string    `json:"token"`
	CreatedAt     time.Time `json:"created_at"`
	IsActive      bool      `json:"is_active"`
}

type User struct {
	TelegramID    int64                    `json:"telegram_id"`
	Servers       map[string]*ServerConfig `json:"servers"`
	CurrentServer string                   `json:"current_server"`
	CreatedAt     time.Time                `json:"created_at"`
}
