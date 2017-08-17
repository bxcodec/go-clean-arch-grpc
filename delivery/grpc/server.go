package grpc

import (
	articleHandler "github.com/bxcodec/go-clean-arch-grpc/delivery/grpc/article"
	articleHandlerGrpc "github.com/bxcodec/go-clean-arch-grpc/delivery/grpc/article/article_grpc"
	"github.com/bxcodec/go-clean-arch-grpc/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewArticleServerGrpc(server *grpc.Server, articleUcase usecase.ArticleUsecase) {
	articleHandlerGrpc.RegisterArticleHandlerServer(server, articleHandler.NewArticleServer(articleUcase))
	reflection.Register(server)
}
