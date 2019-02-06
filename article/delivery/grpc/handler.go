package grpc

import (
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"context"

	"github.com/bxcodec/go-clean-arch-grpc/article/delivery/grpc/article_grpc"

	models "github.com/bxcodec/go-clean-arch-grpc/domain"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
)

func NewArticleServerGrpc(gserver *grpc.Server, articleUcase models.ArticleUsecase) {

	articleServer := &server{
		usecase: articleUcase,
	}

	article_grpc.RegisterArticleHandlerServer(gserver, articleServer)
	reflection.Register(gserver)
}

type server struct {
	usecase models.ArticleUsecase
}

func (s *server) transformArticleRPC(ar *models.Article) *article_grpc.Article {

	if ar == nil {
		return nil
	}

	updated_at := &google_protobuf.Timestamp{

		Seconds: ar.UpdatedAt.Unix(),
	}
	created_at := &google_protobuf.Timestamp{
		Seconds: ar.CreatedAt.Unix(),
	}
	res := &article_grpc.Article{
		ID:        ar.ID,
		Title:     ar.Title,
		Content:   ar.Content,
		UpdatedAt: updated_at,
		CreatedAt: created_at,
	}
	return res
}

func (s *server) transformArticleData(ar *article_grpc.Article) *models.Article {
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

func (s *server) GetArticle(ctx context.Context, in *article_grpc.SingleRequest) (*article_grpc.Article, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}
	ar, err := s.usecase.GetByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	res := s.transformArticleRPC(&ar)
	return res, nil
}

func (s *server) FetchArticle(in *article_grpc.FetchRequest, stream article_grpc.ArticleHandler_FetchArticleServer) error {

	cursor := ""
	num := int64(0)
	if in != nil {
		cursor = in.Cursor
		num = in.Num
	}
	list, _, err := s.usecase.Fetch(cursor, num)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, a := range list {
		ar := s.transformArticleRPC(&a)

		if err := stream.Send(ar); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (s *server) GetListArticle(ctx context.Context, in *article_grpc.FetchRequest) (*article_grpc.ListArticle, error) {
	cursor := ""
	num := int64(0)
	if in != nil {

		cursor = in.Cursor

		num = in.Num
	}
	list, nextCursor, err := s.usecase.Fetch(cursor, num)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	arrArticle := make([]*article_grpc.Article, len(list))
	for i, a := range list {
		ar := s.transformArticleRPC(&a)
		arrArticle[i] = ar
	}
	result := &article_grpc.ListArticle{
		Artilces: arrArticle,
		Cursor:   nextCursor,
	}
	return result, nil
}

func (s *server) UpdateArticle(c context.Context, ar *article_grpc.Article) (*article_grpc.Article, error) {
	a := s.transformArticleData(ar)
	err := s.usecase.Update(a)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	l := s.transformArticleRPC(a)
	return l, nil
}

func (s *server) Delete(c context.Context, in *article_grpc.SingleRequest) (*article_grpc.DeleteResponse, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}

	err := s.usecase.Delete(id)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		return nil, err
	}
	resp := &article_grpc.DeleteResponse{
		Status: "Not Oke To Delete",
	}
	resp.Status = "Succesfull To Delete"
	return resp, nil
}

func (s *server) Store(ctx context.Context, a *article_grpc.Article) (*article_grpc.Article, error) {
	ar := s.transformArticleData(a)
	err := s.usecase.Store(ar)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := s.transformArticleRPC(ar)
	return res, nil
}

func (s *server) BatchInsert(stream article_grpc.ArticleHandler_BatchInsertServer) error {
	errs := make([]*article_grpc.ErrorMessage, 0)
	totalSukses := int64(0)
	for {
		article, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&article_grpc.BatchInsertResponse{
				Errors:       errs,
				TotalSuccess: totalSukses,
			})
		}
		if err != nil {
			log.Println(err.Error())
			return err
		}
		a := s.transformArticleData(article)
		err = s.usecase.Store(a)
		if err != nil {
			log.Println(err.Error())
			e := &article_grpc.ErrorMessage{
				Message: err.Error(),
			}
			errs = append(errs, e)
		}
		if err == nil {
			totalSukses++
		}
	}

}

func (s *server) BatchUpdate(stream article_grpc.ArticleHandler_BatchUpdateServer) error {
	for {
		ar, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Println(err.Error())
			return err
		}

		a := s.transformArticleData(ar)
		er := s.usecase.Update(a)

		if er != nil {
			log.Println("Something error when updating Article", er)
			return er
		}
		if a == nil {
			log.Println("Article Not Found")
			return models.NOT_FOUND_ERROR
		}
		res := s.transformArticleRPC(a)
		if err := stream.Send(res); err != nil {
			log.Println(err.Error())
			return err
		}
	}
}
