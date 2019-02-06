package usecase

import (
	"strconv"
	"time"

	models "github.com/bxcodec/go-clean-arch-grpc/domain"
)

type articleUsecase struct {
	articleRepos models.ArticleRepository
}

func NewArticleUsecase(a models.ArticleRepository) models.ArticleUsecase {
	return &articleUsecase{a}
}

func (a *articleUsecase) Fetch(cursor string, num int64) ([]models.Article, string, error) {
	if num == 0 {
		num = 10
	}

	listArticle, err := a.articleRepos.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""

	if size := len(listArticle); size == int(num) {
		lastId := listArticle[num-1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	return listArticle, nextCursor, nil
}

func (a *articleUsecase) GetByID(id int64) (models.Article, error) {

	return a.articleRepos.GetByID(id)
}

func (a *articleUsecase) Update(ar *models.Article) error {
	_, err := a.articleRepos.GetByID(ar.ID)
	if err != nil {
		return err
	}

	ar.UpdatedAt = time.Now()
	return a.articleRepos.Update(ar)
}

func (a *articleUsecase) GetByTitle(title string) (models.Article, error) {

	return a.articleRepos.GetByTitle(title)
}

func (a *articleUsecase) Store(m *models.Article) error {

	_, err := a.GetByTitle(m.Title)
	if err != nil {
		if err == models.NOT_FOUND_ERROR {
			return models.CONFLIT_ERROR
		}
		return err
	}

	return a.articleRepos.Store(m)

}

func (a *articleUsecase) Delete(id int64) error {
	_, err := a.GetByID(id)
	if err != nil {
		return err
	}
	return a.articleRepos.Delete(id)
}
