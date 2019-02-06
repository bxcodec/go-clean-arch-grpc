package repository

import (
	"database/sql"
	"errors"
	"log"

	models "github.com/bxcodec/go-clean-arch-grpc/domain"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

func NewMysqlArticleRepository(Conn *sql.DB) models.ArticleRepository {

	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) fetch(query string, args ...interface{}) ([]models.Article, error) {

	rows, err := m.Conn.Query(query, args...)

	if err != nil {
		log.Fatal(err)
		return nil, models.INTERNAL_SERVER_ERROR
	}
	defer rows.Close()
	result := make([]models.Article, 0)
	for rows.Next() {
		t := new(models.Article)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
			return nil, models.INTERNAL_SERVER_ERROR
		}
		result = append(result, *t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(cursor string, num int64) ([]models.Article, error) {

	query := `SELECT id,title,content,updated_at, created_at
  						FROM article WHERE ID > ? LIMIT ?`

	return m.fetch(query, cursor, num)

}
func (m *mysqlArticleRepository) GetByID(id int64) (models.Article, error) {
	query := `SELECT id,title,content,updated_at, created_at
  						FROM article WHERE ID = ?`
	var res models.Article
	list, err := m.fetch(query, id)
	if err != nil {

		return res, err
	}
	if len(list) > 0 {
		res = list[0]
	} else {

		return res, models.NOT_FOUND_ERROR
	}
	return res, nil
}

func (m *mysqlArticleRepository) GetByTitle(title string) (models.Article, error) {
	query := `SELECT id,title,content,updated_at, created_at
  						FROM article WHERE title = ?`
	var res models.Article
	list, err := m.fetch(query, title)
	if err != nil {
		return res, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, models.NOT_FOUND_ERROR
	}
	return res, nil
}

func (m *mysqlArticleRepository) Store(a *models.Article) error {
	query := `INSERT  article SET title=? , content=? , updated_at=? , created_at=?`
	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return models.INTERNAL_SERVER_ERROR
	}
	res, err := stmt.Exec(a.Title, a.Content, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return models.INTERNAL_SERVER_ERROR
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (m *mysqlArticleRepository) Delete(id int64) error {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return models.INTERNAL_SERVER_ERROR
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
		return models.INTERNAL_SERVER_ERROR
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return models.INTERNAL_SERVER_ERROR
	}
	if rowsAfected <= 0 {
		return models.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (m *mysqlArticleRepository) Update(ar *models.Article) error {
	query := `UPDATE article set title=?, content=?, updated_at=? WHERE ID = ?`
	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	excRes, err := stmt.Exec(ar.Title, ar.Content, ar.UpdatedAt, ar.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	affect, err := excRes.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if affect < 1 {
		return errors.New("Nothing Affected. Make sure your article is exist in DB")
	}

	return nil
}
