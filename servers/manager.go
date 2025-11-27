package servers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"
	
	"github.com/kfilin/watchtower-masterbot/internal/api"
)

type ServerManager struct {
	users map[int64]*User
	mu    sync.RWMutex
	key   []byte
}

func NewManager(encryptionKey string) *ServerManager {
	key := deriveKey(encryptionKey)
	return &ServerManager{
		users: make(map[int64]*User),
		key:   key,
	}
}

func (sm *ServerManager) AddServer(userID int64, nickname, watchtowerURL, token string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	user, exists := sm.users[userID]
	if !exists {
		user = &User{
			TelegramID: userID,
			Servers:    make(map[string]*ServerConfig),
			CreatedAt:  time.Now(),
		}
		sm.users[userID] = user
	}

	if _, exists := user.Servers[nickname]; exists {
		return errors.New("server with this nickname already exists")
	}

	encryptedToken, err := sm.encryptToken(token)
	if err != nil {
		return err
	}

	user.Servers[nickname] = &ServerConfig{
		Nickname:      nickname,
		WatchtowerURL: watchtowerURL,
		Token:         encryptedToken,
		CreatedAt:     time.Now(),
		IsActive:      true,
	}

	if user.CurrentServer == "" {
		user.CurrentServer = nickname
	}

	return nil
}

func (sm *ServerManager) GetCurrentServer(userID int64) (*ServerConfig, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	user, exists := sm.users[userID]
	if !exists || user.CurrentServer == "" {
		return nil, errors.New("no servers configured")
	}

	server, exists := user.Servers[user.CurrentServer]
	if !exists {
		return nil, errors.New("current server not found")
	}

	decryptedToken, err := sm.decryptToken(server.Token)
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		Nickname:      server.Nickname,
		WatchtowerURL: server.WatchtowerURL,
		Token:         decryptedToken,
		CreatedAt:     server.CreatedAt,
		IsActive:      server.IsActive,
	}, nil
}

func (sm *ServerManager) SwitchServer(userID int64, nickname string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	user, exists := sm.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	if _, exists := user.Servers[nickname]; !exists {
		return errors.New("server not found")
	}

	user.CurrentServer = nickname
	return nil
}

func (sm *ServerManager) ListServers(userID int64) ([]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	user, exists := sm.users[userID]
	if !exists {
		return nil, errors.New("no servers configured")
	}

	servers := make([]string, 0, len(user.Servers))
	for nickname := range user.Servers {
		servers = append(servers, nickname)
	}

	return servers, nil
}

// GetAPIClient returns a Watchtower API client for the user's current server
func (sm *ServerManager) GetAPIClient(userID int64) (*api.WatchtowerClient, error) {
	server, err := sm.GetCurrentServer(userID)
	if err != nil {
		return nil, err
	}
	
	return api.NewWatchtowerClient(server.WatchtowerURL, server.Token), nil
}

// FIXED: Using modern CFB encryption without deprecated functions
func (sm *ServerManager) encryptToken(plaintext string) (string, error) {
	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Modern CFB encryption - create new buffer for encrypted data
	stream := cipher.NewCFBEncrypter(block, iv)
	plaintextBytes := []byte(plaintext)
	encryptedData := make([]byte, len(plaintextBytes))
	stream.XORKeyStream(encryptedData, plaintextBytes)
	
	// Combine IV and encrypted data
	copy(ciphertext[aes.BlockSize:], encryptedData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// FIXED: Using modern CFB decryption without deprecated functions
func (sm *ServerManager) decryptToken(cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	encryptedData := ciphertext[aes.BlockSize:]

	// Modern CFB decryption - create new buffer for decrypted data
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(encryptedData))
	stream.XORKeyStream(plaintext, encryptedData)

	return string(plaintext), nil
}

func deriveKey(passphrase string) []byte {
	key := make([]byte, 32)
	copy(key, passphrase)
	return key
}
