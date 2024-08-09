package repository_test

import (
	"context"
	"errors"
	mocks2 "github.com/amitshekhariitbhu/go-backend-clean-architecture/database/mongo/mocks"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {

	var databaseHelper *mocks2.Database
	var collectionHelper *mocks2.Collection

	databaseHelper = &mocks2.Database{}
	collectionHelper = &mocks2.Collection{}

	collectionName := domain.CollectionUser

	mockUser := &domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	mockEmptyUser := &domain.User{}
	mockUserID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockUserID, nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}
