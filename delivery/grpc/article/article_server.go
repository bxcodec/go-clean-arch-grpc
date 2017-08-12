package article

import (
	"io"
	"time"

	"github.com/bxcodec/go-clean-arch-grpc/models"
	"github.com/bxcodec/go-clean-arch-grpc/usecase"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
)

type server struct {
	usecase usecase.ArticleUsecase
}

func (s *server) transformArticleRPC(ar *models.Article) *Article {
	updated_at := &google_protobuf.Timestamp{

		Seconds: ar.UpdatedAt.Unix(),
	}
	created_at := &google_protobuf.Timestamp{
		Seconds: ar.CreatedAt.Unix(),
	}
	res := &Article{
		ID:        ar.ID,
		Title:     ar.Title,
		Content:   ar.Content,
		UpdatedAt: updated_at,
		CreatedAt: created_at,
	}
	return res
}

func (s *server) transformArticleData(ar *Article) *models.Article {
	updated_at := time.Unix(ar.GetUpdatedAt().GetSeconds(), 0)
	created_at := time.Unix(ar.GetCreatedAt().GetSeconds(), 0)
	res := &models.Article{
		ID:        ar.ID,
		Title:     ar.Title,
		Content:   ar.Content,
		UpdatedAt: updated_at,
		CreatedAt: created_at,
	}
	return res
}

func (s *server) GetArticle(ctx context.Context, in *SingleRequest) (*Article, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}
	ar, err := s.usecase.GetByID(id)
	if err != nil {
		return nil, err
	}

	res := s.transformArticleRPC(ar)
	return res, nil
}

func (s *server) FetchArticle(in *FetchRequest, stream ArticleHandler_FetchArticleServer) error {

	cursor := ""
	num := int64(0)
	if in != nil {
		cursor = in.Cursor
		num = in.Num
	}
	list, _, err := s.usecase.Fetch(cursor, num)
	if err != nil {
		return err
	}

	for _, a := range list {
		ar := s.transformArticleRPC(a)

		if err := stream.Send(ar); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) GetListArticle(ctx context.Context, in *FetchRequest) (*ListArticle, error) {
	cursor := ""
	num := int64(0)
	if in != nil {

		cursor = in.Cursor

		num = in.Num
	}
	list, nextCursor, err := s.usecase.Fetch(cursor, num)

	if err != nil {
		return nil, err
	}
	arrArticle := make([]*Article, len(list))
	for i, a := range list {
		ar := s.transformArticleRPC(a)
		arrArticle[i] = ar
	}
	result := &ListArticle{
		Artilces: arrArticle,
		Cursor:   nextCursor,
	}
	return result, nil
}

func (s *server) UpdateArticle(c context.Context, ar *Article) (*Article, error) {
	a := s.transformArticleData(ar)
	res, err := s.usecase.Update(a)
	if err != nil {
		return nil, nil
	}
	l := s.transformArticleRPC(res)
	return l, nil
}

func (s *server) Delete(c context.Context, in *SingleRequest) (*DeleteResponse, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}

	ok, err := s.usecase.Delete(id)
	if err != nil {
		return nil, err
	}
	resp := &DeleteResponse{
		Status: "Not Oke To Delete",
	}
	if ok {
		resp.Status = "Succesfull To Delete"
	}

	return resp, nil
}

func (s *server) Store(ctx context.Context, a *Article) (*Article, error) {
	ar := s.transformArticleData(a)
	data, err := s.usecase.Store(ar)
	if err != nil {
		return nil, err
	}
	res := s.transformArticleRPC(data)

	return res, nil
}

func (s *server) BatchInsert(stream ArticleHandler_BatchInsertServer) error {
	errs := make([]*ErrorMessage, 0)
	totalSukses := int64(0)
	for {
		article, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&BatchInsertResponse{
				Errors:       errs,
				TotalSuccess: totalSukses,
			})
		}
		if err != nil {
			return err
		}
		a := s.transformArticleData(article)
		res, err := s.usecase.Store(a)
		if err != nil {
			e := &ErrorMessage{
				Message: err.Error(),
			}
			errs = append(errs, e)
		}
		if res != nil {
			totalSukses++
		}
	}

}

func (s *server) BatchUpdate(stream ArticleHandler_BatchUpdateServer) error {
	for {
		ar, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		a := s.transformArticleData(ar)
		a, er := s.usecase.Update(a)
		if er != nil {
			return er
		}
		res := s.transformArticleRPC(a)
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

func NewArticleServer(u usecase.ArticleUsecase) ArticleHandlerServer {
	return &server{usecase: u}
}
