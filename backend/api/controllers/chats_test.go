package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/BogPin/real-time-chat/backend/api/services"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func getChatService(db *sql.DB) services.ChatService {
	messageStorer := models.NewMessageStorer(db)
	participantStorer := models.NewParticipantStorer(db)
	participantService := services.NewParticipantService(participantStorer)
	chatStorer := models.NewChatStorer(db)
	return services.NewChatService(chatStorer, participantStorer, messageStorer, participantService)
}

func TestCreateChat1(t *testing.T) {
	// Створення тіла запиту
	chat := models.ChatFromRequest{
		Title: "Test Chat",
	}
	body, _ := json.Marshal(chat)

	// Створення POST-запиту
	createChatEndpoint := fmt.Sprintf("%s/api/chats", apiServer.URL)
	req, err := http.NewRequest("POST", createChatEndpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create createChat request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))
	req.Header.Set("Content-Type", "application/json")

	// Виконання запиту
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send createChat request: %v", err)
	}
	defer resp.Body.Close()

	// Перевірка статусу відповіді
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Декодування результату
	var createdChat models.Chat
	err = json.NewDecoder(resp.Body).Decode(&createdChat)
	if err != nil {
		t.Fatalf("Failed to decode createChat response: %v", err)
	}
	// Перевірка результатів
	assert.Equal(t, chat.Title, createdChat.Title)
}

func TestGetChatByID(t *testing.T) {
	// Отримання ідентифікатора створеного чату
	chatID := 1
	// Створення GET-запиту для отримання чату за його ідентифікатором
	getChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chatID)
	getReq, err := http.NewRequest("GET", getChatEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create getChat request: %v", err)
	}
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))

	// Виконання GET-запиту для отримання чату
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Failed to send getChat request: %v", err)
	}
	defer getResp.Body.Close()

	// Перевірка статусу відповіді на GET-запит
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	// Декодування результату GET-запиту
	var retrievedChat models.Chat
	err = json.NewDecoder(getResp.Body).Decode(&retrievedChat)
	if err != nil {
		t.Fatalf("Failed to decode getChat response: %v", err)
	}

	// Перевірка результатів
	assert.Equal(t, chatID, retrievedChat.Id)
}

func TestGetAllChats(t *testing.T) {

	// Створення GET-запиту для отримання всіх чатів
	getChatsEndpoint := fmt.Sprintf("%s/api/chats", apiServer.URL)
	getReq, err := http.NewRequest("GET", getChatsEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create getChats request: %v", err)
	}
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenN))

	// Виконання GET-запиту для отримання всіх чатів
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Failed to send getChats request: %v", err)
	}
	defer getResp.Body.Close()

	// Перевірка статусу відповіді на GET-запит
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	// Декодування результату GET-запиту
	var chats []models.Chat
	err = json.NewDecoder(getResp.Body).Decode(&chats)
	if err != nil {
		t.Fatalf("Failed to decode getChats response: %v", err)
	}

	// Перевірка результатів
	assert.Equal(t, 2, len(chats))
}

func TestUpdateChatUserNotAdmin(t *testing.T) {
	// Створення тіла запиту для створення чату
	chat := models.Chat{
		Id:        2,
		Title:     "Updated",
		CreatorId: 1,
		CreatedAt: "2023-05-16",
	}

	updatedBody, _ := json.Marshal(chat)

	// Створення PATCH-запиту для оновлення чату
	updateChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chat.Id)
	updateReq, err := http.NewRequest("PATCH", updateChatEndpoint, bytes.NewReader(updatedBody))
	if err != nil {
		t.Fatalf("Failed to create updateChat request: %v", err)
	}
	updateReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))
	updateReq.Header.Set("Content-Type", "application/json")

	// Виконання PATCH-запиту для оновлення чату
	updateResp, err := http.DefaultClient.Do(updateReq)
	if err != nil {
		t.Fatalf("Failed to send updateChat request: %v", err)
	}
	defer updateResp.Body.Close()

	// Перевірка статусу відповіді на PATCH-запит
	assert.Equal(t, http.StatusForbidden, updateResp.StatusCode)
}

func TestUpdateChatUserIsAdmin(t *testing.T) {
	// Створення тіла запиту для створення чату
	chat := models.Chat{
		Id:        1,
		Title:     "Updated",
		CreatorId: 1,
		CreatedAt: "2023-05-16",
	}

	updatedBody, _ := json.Marshal(chat)

	// Створення PATCH-запиту для оновлення чату
	updateChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chat.Id)
	updateReq, err := http.NewRequest("PATCH", updateChatEndpoint, bytes.NewReader(updatedBody))
	if err != nil {
		t.Fatalf("Failed to create updateChat request: %v", err)
	}
	updateReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))
	updateReq.Header.Set("Content-Type", "application/json")

	// Виконання PATCH-запиту для оновлення чату
	updateResp, err := http.DefaultClient.Do(updateReq)
	if err != nil {
		t.Fatalf("Failed to send updateChat request: %v", err)
	}
	defer updateResp.Body.Close()

	// Перевірка статусу відповіді на PATCH-запит
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	// Виконання GET-запиту для отримання оновленого чату
	getChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chat.Id)
	getReq, err := http.NewRequest("GET", getChatEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create getChat request: %v", err)
	}
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))

	// Виконання GET-запиту для отримання оновленого чату
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Failed to send getChat request: %v", err)
	}
	defer getResp.Body.Close()

	// Перевірка статусу відповіді на GET-запит
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	// Декодування результату GET-запиту
	var updatedChat models.Chat
	err = json.NewDecoder(getResp.Body).Decode(&updatedChat)
	if err != nil {
		t.Fatalf("Failed to decode getChat response: %v", err)
	}

	// Перевірка, що заголовок чату був успішно оновлений
	assert.Equal(t, chat.Title, updatedChat.Title)
}

func TestDeleteChat(t *testing.T) {
	// Отримання ідентифікатора створеного чату
	chatID := 3
	// Створення DELETE-запиту для отримання чату за його ідентифікатором
	getChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chatID)
	deleteReq, err := http.NewRequest("DELETE", getChatEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create deleteChat request: %v", err)
	}
	deleteReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))

	// Виконання delete-запиту для отримання чату
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("Failed to send deleteChat request: %v", err)
	}
	defer deleteResp.Body.Close()

	// Перевірка статусу відповіді на DELETE-запит
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)

	getReq, err := http.NewRequest("GET", getChatEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create getChat request: %v", err)
	}
	getReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenB))

	// Виконання GET-запиту для отримання інформації про чат
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Failed to send getChat request: %v", err)
	}
	defer getResp.Body.Close()

	// Перевірка статусу відповіді на GET-запит
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestDeleteChatNotAdmin(t *testing.T) {
	chatID := 3
	// Створення DELETE-запиту для отримання чату за його ідентифікатором
	getChatEndpoint := fmt.Sprintf("%s/api/chats/%d", apiServer.URL, chatID)
	deleteReq, err := http.NewRequest("DELETE", getChatEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create deleteChat request: %v", err)
	}
	deleteReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokenN))

	// Виконання DELETE-запиту для отримання чату
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("Failed to send deleteChat request: %v", err)
	}
	defer deleteResp.Body.Close()

	// Перевірка статусу відповіді на DELETE-запит
	assert.Equal(t, http.StatusForbidden, deleteResp.StatusCode)
}
