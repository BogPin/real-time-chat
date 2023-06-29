package services_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/BogPin/real-time-chat/backend/api/models"
	models_mocks "github.com/BogPin/real-time-chat/backend/api/models/mocks"
	"github.com/BogPin/real-time-chat/backend/api/services"
	services_mocks "github.com/BogPin/real-time-chat/backend/api/services/mocks"
	"github.com/BogPin/real-time-chat/backend/api/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatSuccess(t *testing.T) {
	//Arrange
	title := "test-chat"
	creatorId := 1
	chatFromRequest := models.ChatFromRequest{Title: title}
	expectedDTO := models.ChatDTO{
		Title:     title,
		CreatorId: creatorId,
	}
	expectedChat := models.Chat{
		Id:        1,
		Title:     title,
		CreatorId: creatorId,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: creatorId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Begin().
		Return(db.Begin())
	mockChatStorer.
		EXPECT().
		CreateInTx(gomock.AssignableToTypeOf(&sql.Tx{}), expectedDTO).
		Return(&expectedChat, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		CreateInTx(gomock.AssignableToTypeOf(&sql.Tx{}), expectedParticipant).
		Return(&expectedParticipant, nil)

	chatsService := services.ChatService{
		ChatStorer:        mockChatStorer,
		ParticipantStorer: mockParticipantStorer,
	}

	//Act
	actualChat, httpErr := chatsService.Create(creatorId, chatFromRequest)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, httpErr)
}

func TestCreateChatInternalError(t *testing.T) {
	//Arrange
	title := "test-chat"
	creatorId := 1
	chatFromRequest := models.ChatFromRequest{Title: title}
	expectedDTO := models.ChatDTO{
		Title:     title,
		CreatorId: creatorId,
	}
	expectedChat := models.Chat{
		Id:        1,
		Title:     title,
		CreatorId: creatorId,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: creatorId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusInternalServerError)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occured while opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Begin().
		Return(db.Begin())
	mockChatStorer.
		EXPECT().
		CreateInTx(gomock.AssignableToTypeOf(&sql.Tx{}), expectedDTO).
		Return(&expectedChat, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		CreateInTx(gomock.AssignableToTypeOf(&sql.Tx{}), expectedParticipant).
		Return(nil, expectedError)

	chatsService := services.ChatService{
		ChatStorer:        mockChatStorer,
		ParticipantStorer: mockParticipantStorer,
	}

	//Act
	actualChat, httpErr := chatsService.Create(creatorId, chatFromRequest)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestGetOneChatSuccess(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		GetOne(expectedChat.Id).
		Return(&expectedChat, nil)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.GetOne(userId, expectedChat.Id)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, httpErr)
}

func TestGetOneChatUserNotInChatError(t *testing.T) {
	//Arrange
	userId := 1
	chatId := 1
	expectedChat := models.Chat{Id: chatId, Title: "hejfe", CreatorId: 2, CreatedAt: "2023-06-10"}
	expectedError := fmt.Errorf("user %d doesn't participate in chat %d", userId, chatId)
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusForbidden)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.EXPECT().GetOne(chatId).Return(&expectedChat, nil)

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, chatId).
		Return(false, nil)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.GetOne(userId, chatId)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestGetOneChatInternalError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusInternalServerError)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		GetOne(expectedChat.Id).
		Return(nil, expectedError)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.GetOne(userId, expectedChat.Id)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestGetUserChatsSuccess(t *testing.T) {
	//Arrange
	userId := 1
	expectedChats := []models.Chat{
		{
			Id:        1,
			Title:     "test-chat",
			CreatorId: 2,
			CreatedAt: "2023-06-27",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		GetUserChats(userId).
		Return(expectedChats, nil)

	chatsService := services.ChatService{
		ChatStorer: mockChatStorer,
	}

	//Act
	actualChats, httpErr := chatsService.GetUserChats(userId)

	//Assert
	assert.Equal(t, expectedChats, actualChats)
	assert.Nil(t, httpErr)
}

func TestGetUserChatsInternalError(t *testing.T) {
	//Arrange
	userId := 1
	expectedError := errors.New("SOME SQL CONNECTION ERROR")
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusInternalServerError)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		GetUserChats(userId).
		Return(nil, expectedError)

	chatsService := services.ChatService{
		ChatStorer: mockChatStorer,
	}

	//Act
	actualChats, httpErr := chatsService.GetUserChats(userId)

	//Assert
	assert.Nil(t, actualChats)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestUpdateChatSuccess(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Update(expectedChat).
		Return(&expectedChat, nil)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Update(userId, expectedChat)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, httpErr)
}

func TestUpdateChatUserNotInChatError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedError := fmt.Errorf("user %d doesn't participate in chat %d", userId, expectedChat.Id)
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusForbidden)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(false, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Update(userId, expectedChat)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestUpdateChatUserNotAdminError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "member",
	}
	expectedError := fmt.Errorf("user %d doesn't have permission to update chat %d", userId, expectedChat.Id)
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusForbidden)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Update(userId, expectedChat)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestUpdateChatInternalError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusInternalServerError)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Update(expectedChat).
		Return(nil, expectedError)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Update(userId, expectedChat)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestDeleteChatSuccess(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockMessagesStorer := models_mocks.NewMockIMessageStorer(ctrl)
	mockMessagesStorer.
		EXPECT().
		DeleteAll(expectedChat.Id).
		Times(1)

	mockParticipantStorer.
		EXPECT().
		DeleteAll(expectedChat.Id).
		Times(1)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Delete(expectedChat.Id).
		Return(&expectedChat, nil)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		MessageStorer:      mockMessagesStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Delete(userId, expectedChat.Id)

	//Assert
	assert.Equal(t, expectedChat, *actualChat)
	assert.Nil(t, httpErr)
}

func TestDeleteChatUserNotInChatError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedError := fmt.Errorf("user %d doesn't participate in chat %d", userId, expectedChat.Id)
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusForbidden)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(false, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockMessagesStorer := models_mocks.NewMockIMessageStorer(ctrl)
	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		ParticipantStorer:  mockParticipantStorer,
		MessageStorer:      mockMessagesStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Delete(userId, expectedChat.Id)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestDeleteChatUserNotAdminError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "member",
	}
	expectedError := fmt.Errorf("user %d doesn't have permission to update chat %d", userId, expectedChat.Id)
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusForbidden)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockMessagesStorer := models_mocks.NewMockIMessageStorer(ctrl)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		MessageStorer:      mockMessagesStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Delete(userId, expectedChat.Id)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}

func TestDeleteChatInternalError(t *testing.T) {
	//Arrange
	userId := 1
	expectedChat := models.Chat{
		Id:        1,
		Title:     "test-chat",
		CreatorId: 2,
		CreatedAt: "2023-06-27",
	}
	expectedParticipant := models.Participant{
		UserId: userId,
		ChatId: expectedChat.Id,
		Role:   "admin",
	}
	expectedError := errors.New("SOME SQL CONNECTION ERROR")
	expectedHTTPError := utils.NewHttpError(expectedError, http.StatusInternalServerError)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParticipantService := services_mocks.NewMockIParticipantService(ctrl)
	mockParticipantService.
		EXPECT().
		UserInChat(userId, expectedChat.Id).
		Return(true, nil)

	mockParticipantStorer := models_mocks.NewMockIParticipantStorer(ctrl)
	mockParticipantStorer.
		EXPECT().
		GetOne(userId, expectedChat.Id).
		Return(&expectedParticipant, nil)

	mockMessagesStorer := models_mocks.NewMockIMessageStorer(ctrl)
	mockMessagesStorer.
		EXPECT().
		DeleteAll(expectedChat.Id).
		Times(1)

	mockParticipantStorer.
		EXPECT().
		DeleteAll(expectedChat.Id).
		Times(1)

	mockChatStorer := models_mocks.NewMockIChatStorer(ctrl)
	mockChatStorer.
		EXPECT().
		Delete(expectedChat.Id).
		Return(nil, expectedError)

	chatsService := services.ChatService{
		ChatStorer:         mockChatStorer,
		MessageStorer:      mockMessagesStorer,
		ParticipantStorer:  mockParticipantStorer,
		ParticipantService: mockParticipantService,
	}

	//Act
	actualChat, httpErr := chatsService.Delete(userId, expectedChat.Id)

	//Assert
	assert.Nil(t, actualChat)
	assert.Equal(t, expectedHTTPError, httpErr)
}
