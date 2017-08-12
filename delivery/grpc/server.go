package grpc

import (
	articleHandler "github.com/bxcodec/go-clean-arch-grpc/delivery/grpc/article"
	"github.com/bxcodec/go-clean-arch-grpc/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewArticleServerGrpc(server *grpc.Server, articleUcase usecase.ArticleUsecase) {
	articleHandler.RegisterArticleHandlerServer(server, articleHandler.NewArticleServer(articleUcase))
	reflection.Register(server)
}
