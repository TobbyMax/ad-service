package tests

import (
	"context"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/user"
	"log"
	"testing"
)

var (
	BenchSink int64 // make sure compiler cannot optimize away benchmarks
	U         = user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	Ad        = ads.Ad{Title: "Missed Calls", Text: "Blue Slide Park"}
)

func BenchmarkRepoCreateUser(b *testing.B) {
	ctx := context.Background()
	repo := adrepo.New()
	for i := 0; i < b.N; i++ {
		id, err := repo.AddUser(ctx, U)
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink += id
	}
}

func BenchmarkRepoAddAd(b *testing.B) {
	ctx := context.Background()
	repo := adrepo.New()
	uid, err := repo.AddUser(ctx, U)
	if err != nil {
		log.Fatalf("Function returned error: %v", err)
	}
	Ad.AuthorID = uid
	for i := 0; i < b.N; i++ {
		_, err := repo.AddAd(ctx, Ad)
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
}

func BenchmarkRepoListAds10(b *testing.B) {
	n := 10
	ctx := context.Background()
	repo := adrepo.New()
	uid, err := repo.AddUser(ctx, U)
	if err != nil {
		log.Fatalf("Function returned error: %v", err)
	}
	Ad.AuthorID = uid
	for i := 0; i < n; i++ {
		_, err := repo.AddAd(ctx, Ad)
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAdList(ctx, app.ListAdsParams{})
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
}

func BenchmarkRepoListAds100(b *testing.B) {
	n := 100
	ctx := context.Background()
	repo := adrepo.New()
	uid, err := repo.AddUser(ctx, U)
	if err != nil {
		log.Fatalf("Function returned error: %v", err)
	}
	Ad.AuthorID = uid
	for i := 0; i < n; i++ {
		_, err := repo.AddAd(ctx, Ad)
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAdList(ctx, app.ListAdsParams{})
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
}

func BenchmarkRepoListAds1000(b *testing.B) {
	n := 1000
	ctx := context.Background()
	repo := adrepo.New()
	uid, err := repo.AddUser(ctx, U)
	if err != nil {
		log.Fatalf("Function returned error: %v", err)
	}
	Ad.AuthorID = uid
	for i := 0; i < n; i++ {
		_, err := repo.AddAd(ctx, Ad)
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAdList(ctx, app.ListAdsParams{})
		if err != nil {
			log.Fatalf("Function returned error: %v", err)
		}
		BenchSink++
	}
}
