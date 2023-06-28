package models_test

import (
	"errors"
	"testing"

	"github.com/BogPin/real-time-chat/backend/api/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	dto := models.ChatDTO{Title: expectedChat.Title, CreatorId: expectedChat.CreatorId}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	row := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("INSERT INTO chats").WillReturnRows(row)

	//Act
	actualChat, err := chatStorer.Create(dto)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, err)
}

func TestCreateChatInternalError(t *testing.T) {
	//Arrange
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	dto := models.ChatDTO{Title: expectedChat.Title, CreatorId: expectedChat.CreatorId}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectQuery("INSERT INTO chats").WillReturnError(expectedError)

	//Act
	actualChat, err := chatStorer.Create(dto)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedError, err)
}

func TestCreateChatInTxSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	dto := models.ChatDTO{Title: expectedChat.Title, CreatorId: expectedChat.CreatorId}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectBegin()
	row := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("INSERT INTO chats").WillReturnRows(row)

	//Act
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' occured while begining transaction", err)
	}
	actualChat, err := chatStorer.CreateInTx(tx, dto)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, err)
}

func TestCreateChatInTxInternalError(t *testing.T) {
	//Arrange
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	dto := models.ChatDTO{Title: expectedChat.Title, CreatorId: expectedChat.CreatorId}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO chats").WillReturnError(expectedError)

	//Act
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' occured while begining transaction", err)
	}
	actualChat, err := chatStorer.CreateInTx(tx, dto)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedError, err)
}

func TestGetOneChatSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	row := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("SELECT (.+) FROM chats").WillReturnRows(row)

	//Act
	actualChat, err := chatStorer.GetOne(1)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, err)
}

func TestGetOneChatInternalError(t *testing.T) {
	//Arrange
	chatId := 1
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectQuery("SELECT (.+) FROM chats").WillReturnError(expectedError)

	//Act
	actualChat, err := chatStorer.GetOne(chatId)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedError, err)
}

func TestGetUserChatsSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	expectedChats := []models.Chat{expectedChat, expectedChat, expectedChat}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	rows := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("SELECT (.+) FROM chats").WillReturnRows(rows)

	//Act
	actualChats, err := chatStorer.GetUserChats(1)

	//Assert
	assert.Equal(t, expectedChats, actualChats)
	assert.Nil(t, err)
}

func TestGetUserChatsInternalError(t *testing.T) {
	//Arrange
	userId := 1
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM chats").WillReturnError(expectedError)

	chatStorer := models.NewChatStorer(db)

	//Act
	actualChats, err := chatStorer.GetUserChats(userId)

	//Assert
	assert.Nil(t, actualChats)
	assert.Equal(t, expectedError, err)
}

func TestUpdateChatSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	row := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("UPDATE chats").WillReturnRows(row)

	//Act
	actualChat, err := chatStorer.Update(expectedChat)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, err)
}

func TestUpdateChatInternalError(t *testing.T) {
	//Arrange
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectQuery("UPDATE chats").WillReturnError(expectedError)

	//Act
	actualChat, err := chatStorer.Update(expectedChat)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedError, err)
}

func TestDeleteChatSuccess(t *testing.T) {
	//Arrange
	chatIdStr := "1"
	creatorIdStr := "1"
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 1,
		CreatedAt: "2023-06-27",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	row := sqlmock.NewRows([]string{"id", "title", "creator_id", "created_at"}).
		AddRow(chatIdStr, expectedChat.Title, creatorIdStr, expectedChat.CreatedAt)
	mock.ExpectQuery("DELETE FROM chats").WillReturnRows(row)

	//Act
	actualChat, err := chatStorer.Delete(1)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, err)
}

func TestDeleteChatInternalError(t *testing.T) {
	//Arrange
	chatId := 1
	expectedError := errors.New("SOME SQL CONNECTION ERROR")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	chatStorer := models.NewChatStorer(db)

	mock.ExpectQuery("DELETE FROM chats").WillReturnError(expectedError)

	//Act
	actualChat, err := chatStorer.Delete(chatId)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedError, err)
}
