package article

import (
	"testing"

	"github.com/bxcodec/faker"
	"github.com/bxcodec/go-clean-arch-grpc/delivery/grpc/article/article_grpc"
	"github.com/bxcodec/go-clean-arch-grpc/models"
	"github.com/bxcodec/go-clean-arch-grpc/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"
)

func TestGet(t *testing.T) {
	mockUsecase := &mocks.ArticleUsecase{}

	a := &models.Article{}
	faker.FakeData(a)
	id := int64(2)
	mockUsecase.On("GetByID", mock.AnythingOfType("int64")).Return(a, nil)

	handler := &server{
		usecase: mockUsecase,
	}

	single := &article_grpc.SingleRequest{Id: id}

	res, err := handler.GetArticle(context.Background(), single)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockUsecase.AssertCalled(t, "GetByID", mock.AnythingOfType("int64"))
}

func TestStore(t *testing.T) {
	mockUsecase := &mocks.ArticleUsecase{}

	a := &models.Article{}
	faker.FakeData(a)

	handler := &server{
		usecase: mockUsecase,
	}
	articleFromClient := handler.transformArticleRPC(a)
	mockUsecase.On("Store", a).Return(a, nil)

	res, err := handler.Store(context.Background(), articleFromClient)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockUsecase.AssertCalled(t, "Store", a)
}
func TestDelete(t *testing.T) {
	mockUsecase := &mocks.ArticleUsecase{}

	a := &article_grpc.DeleteResponse{}
	faker.FakeData(a)
	id := int64(2)
	mockUsecase.On("Delete", mock.AnythingOfType("int64")).Return(true, nil)

	handler := &server{
		usecase: mockUsecase,
	}

	single := &article_grpc.SingleRequest{Id: id}

	res, err := handler.Delete(context.Background(), single)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockUsecase.AssertCalled(t, "Delete", mock.AnythingOfType("int64"))

}

func TestUpdate(t *testing.T) {
	mockUsecase := &mocks.ArticleUsecase{}

	a := &models.Article{}
	faker.FakeData(a)

	handler := &server{
		usecase: mockUsecase,
	}
	articleFromClient := handler.transformArticleRPC(a)
	mockUsecase.On("Update", a).Return(a, nil)

	res, err := handler.UpdateArticle(context.Background(), articleFromClient)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	mockUsecase.AssertCalled(t, "Update", a)
}

func TestGetListArticle(t *testing.T) {

	mockUsecase := &mocks.ArticleUsecase{}

	a := &models.Article{}
	faker.FakeData(a)
	list := []*models.Article{
		a,
	}

	mockUsecase.On("Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(list, "next_cursor", nil)

	handler := &server{
		usecase: mockUsecase,
	}

	fetc := &article_grpc.FetchRequest{
		Num:    10,
		Cursor: "sample Cursor"}

	res, err := handler.GetListArticle(context.Background(), fetc)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockUsecase.AssertCalled(t, "Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64"))

}
