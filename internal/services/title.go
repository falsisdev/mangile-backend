package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

type TitleDetailsWorker struct {
	jikanService   *JikanService
	aniListService *AniListService
	httpClient     *http.Client
}

func NewTitleDetailsWorker(jikan *JikanService, anilist *AniListService) *TitleDetailsWorker {
	return &TitleDetailsWorker{
		jikanService:   jikan,
		aniListService: anilist,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *TitleDetailsWorker) GetCompleteTitle(ctx context.Context, id string, malID int) (*models.TitleMedia, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var sanityResults []models.TitleMedia
	params := map[string]any{"id": id, "malId": malID}
	if err := sanity.query(ctx, queries.TitleQuery, params, &sanityResults); err != nil {
		return nil, err
	}

	var singleTitle models.TitleMedia
	if len(sanityResults) == 0 {
		singleTitle = models.TitleMedia{
			SanityDataMissing: true,
		}
		if malID != 0 {
			singleTitle.MyAnimeListID = malID
		}
	} else {
		singleTitle = sanityResults[0]
	}

	targetMalID := singleTitle.MyAnimeListID

	if targetMalID == 0 {
		return &singleTitle, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if malData, err := w.jikanService.FetchDetails(ctx, targetMalID); err == nil {
			singleTitle.ExternalMAL = malData
		}
	}()

	go func() {
		defer wg.Done()
		if aniListData, err := w.aniListService.FetchDetails(ctx, targetMalID); err == nil {
			singleTitle.ExternalAniList = aniListData
		}
	}()

	wg.Wait()

	return &singleTitle, nil
}
