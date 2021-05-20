package data

import (
	"context"
	"log"
	"time"

	pb "github.com/hi20160616/fetcher_web/api/fetcher_web/v1"
	"github.com/hi20160616/fetcher_web/config"
	"google.golang.org/grpc"
)

func NewsSites() []config.Site {
	return config.Value.Sites
}

func List(ctx context.Context, in *pb.ListArticlesRequest) (*pb.Articles, error) {
	as := &pb.Articles{}
	switch in.Domain {
	case "www.bbc.com":
		conn, err := grpc.Dial("localhost:10001", grpc.WithInsecure(), grpc.WithBlock())
		defer conn.Close()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewNewsFetcherClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		resp, err := c.ListArticles(ctx, in)
		if err != nil {
			log.Printf("c.GetPosts err: %+v", err)
		}
		for _, item := range resp.Articles.Articles {
			as.Articles = append(as.Articles, item)
		}
	}
	return as, nil
}
