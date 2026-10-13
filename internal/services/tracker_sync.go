package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
)

type TrackerSyncService struct {
	malService      *MALService
	aniListService  *AniListTrackerService
	mangaDexService *MangaDexTrackerService
}

func NewTrackerSyncService() *TrackerSyncService {
	return &TrackerSyncService{
		malService:      NewMALService(),
		aniListService:  NewAniListTrackerService(),
		mangaDexService: NewMangaDexTrackerService(),
	}
}

// UserIntegrations represents the stored tokens in Sanity auth document
type UserIntegrations struct {
	AniList struct {
		UserID      int    `json:"userId"`
		Username    string `json:"username"`
		AccessToken string `json:"accessToken"`
	} `json:"anilist"`
	MAL struct {
		UserID       int    `json:"userId"`
		Username     string `json:"username"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresAt    string `json:"expiresAt"`
	} `json:"mal"`
	MangaDex struct {
		UserID       string `json:"userId"`
		Username     string `json:"username"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresAt    string `json:"expiresAt"`
	} `json:"mangadex"`
}

func getUserAuthData(ctx context.Context, logtoID string) (authID string, roles []string, integrations UserIntegrations, err error) {
	sanity, err := newSanityClient()
	if err != nil {
		return "", nil, integrations, err
	}

	const q = `*[_type == "auth" && logtoId == $logtoId][0]{
		_id,
		roles,
		integrations
	}`

	var res struct {
		ID           string           `json:"_id"`
		Roles        []string         `json:"roles"`
		Integrations UserIntegrations `json:"integrations"`
	}

	if err := sanity.query(ctx, q, map[string]any{"logtoId": logtoID}, &res); err != nil {
		return "", nil, integrations, err
	}
	if res.ID == "" {
		return "", nil, integrations, fmt.Errorf("kullanici bulunamadi")
	}
	return res.ID, res.Roles, res.Integrations, nil
}

func (s *TrackerSyncService) GetUserTrackersDoc(ctx context.Context, logtoID string) (*models.UserTrackersDocument, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	const q = `*[_type == "userTrackers" && user->logtoId == $logtoId][0]{
		_id,
		_type,
		updatedAt,
		trackers
	}`

	var doc models.UserTrackersDocument
	if err := sanity.query(ctx, q, map[string]any{"logtoId": logtoID}, &doc); err != nil {
		return nil, err
	}
	if doc.Trackers == nil {
		doc.Trackers = []models.SeriesTrackerItem{}
	}
	return &doc, nil
}

// GetSeriesMangaInfo retrieves series global IDs from Sanity manga document
func getSeriesMangaInfo(ctx context.Context, seriesID string) (malID int, anilistID int, mangadexID string, title string) {
	sanity, err := newSanityClient()
	if err != nil {
		return 0, 0, "", ""
	}

	const q = `*[_type == "manga" && (_id == $id || slug.current == $id)][0]{
		title,
		myAnimeListId,
		anilistId,
		mangaDexId
	}`

	var res struct {
		Title         string `json:"title"`
		MyAnimeListID int    `json:"myAnimeListId"`
		AniListID     int    `json:"anilistId"`
		MangaDexID    string `json:"mangaDexId"`
	}

	if err := sanity.query(ctx, q, map[string]any{"id": seriesID}, &res); err == nil {
		return res.MyAnimeListID, res.AniListID, res.MangaDexID, res.Title
	}
	return 0, 0, "", ""
}

// GetSeriesLiveStatus aggregates live status across MAL, AniList, and MangaDex
func (s *TrackerSyncService) GetSeriesLiveStatus(ctx context.Context, logtoID, seriesID string) (*models.SeriesAllTrackersStatus, error) {
	_, _, integrations, err := getUserAuthData(ctx, logtoID)
	if err != nil {
		return nil, err
	}

	doc, _ := s.GetUserTrackersDoc(ctx, logtoID)
	var userItem *models.SeriesTrackerItem
	if doc != nil {
		for i := range doc.Trackers {
			if doc.Trackers[i].SeriesID == seriesID {
				userItem = &doc.Trackers[i]
				break
			}
		}
	}

	globalMAL, globalAniList, globalMD, _ := getSeriesMangaInfo(ctx, seriesID)

	targetMAL := globalMAL
	targetAniList := globalAniList
	targetMD := globalMD
	autoSync := true

	if userItem != nil {
		if userItem.MALID > 0 {
			targetMAL = userItem.MALID
		}
		if userItem.AniListID > 0 {
			targetAniList = userItem.AniListID
		}
		if userItem.MangaDexID != "" {
			targetMD = userItem.MangaDexID
		}
		autoSync = userItem.AutoSync
	}

	result := &models.SeriesAllTrackersStatus{
		SeriesID: seriesID,
		AutoSync: autoSync,
		Trackers: make(map[string]*models.TrackerStatus),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 1. MyAnimeList
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := &models.TrackerStatus{
			Provider:    "mal",
			IsConnected: integrations.MAL.AccessToken != "",
		}
		if targetMAL > 0 {
			st.MatchedID = strconv.Itoa(targetMAL)
			if integrations.MAL.AccessToken != "" {
				if live, err := s.malService.GetUserMangaStatus(ctx, integrations.MAL.AccessToken, targetMAL); err == nil {
					st = live
				}
			}
		}
		mu.Lock()
		result.Trackers["mal"] = st
		mu.Unlock()
	}()

	// 2. AniList
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := &models.TrackerStatus{
			Provider:    "anilist",
			IsConnected: integrations.AniList.AccessToken != "",
		}
		if targetAniList > 0 || targetMAL > 0 {
			if live, err := s.aniListService.GetUserMangaStatus(ctx, integrations.AniList.AccessToken, targetAniList, targetMAL); err == nil {
				st = live
			}
		}
		mu.Lock()
		result.Trackers["anilist"] = st
		mu.Unlock()
	}()

	// 3. MangaDex
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := &models.TrackerStatus{
			Provider:    "mangadex",
			IsConnected: integrations.MangaDex.AccessToken != "",
		}
		if targetMD != "" {
			st.MatchedID = targetMD
			if live, err := s.mangaDexService.GetUserMangaStatus(ctx, integrations.MangaDex.AccessToken, targetMD); err == nil {
				st = live
			}
		}
		mu.Lock()
		result.Trackers["mangadex"] = st
		mu.Unlock()
	}()

	wg.Wait()
	return result, nil
}

// BindTracker binds a tracker for a user, and optionally sets it globally on the manga document if admin/mod
func (s *TrackerSyncService) BindTracker(ctx context.Context, req models.BindTrackerRequest) error {
	authID, roles, _, err := getUserAuthData(ctx, req.LogtoID)
	if err != nil {
		return err
	}

	doc, err := s.GetUserTrackersDoc(ctx, req.LogtoID)
	if err != nil {
		return err
	}

	isAdmin := false
	for _, r := range roles {
		if r == "001" || r == "002" {
			isAdmin = true
			break
		}
	}

	// Update userTrackers document
	found := false
	for i := range doc.Trackers {
		if doc.Trackers[i].SeriesID == req.SeriesID {
			found = true
			switch req.Provider {
			case "mal":
				parsed, _ := strconv.Atoi(req.MatchedID)
				doc.Trackers[i].MALID = parsed
			case "anilist":
				parsed, _ := strconv.Atoi(req.MatchedID)
				doc.Trackers[i].AniListID = parsed
			case "mangadex":
				doc.Trackers[i].MangaDexID = req.MatchedID
			}
			doc.Trackers[i].LastSyncedAt = time.Now().UTC().Format(time.RFC3339)
			break
		}
	}

	if !found {
		newItem := models.SeriesTrackerItem{
			Key:          generateKey(),
			SeriesID:     req.SeriesID,
			SeriesTitle:  req.SeriesTitle,
			AutoSync:     true,
			LastSyncedAt: time.Now().UTC().Format(time.RFC3339),
		}
		switch req.Provider {
		case "mal":
			parsed, _ := strconv.Atoi(req.MatchedID)
			newItem.MALID = parsed
		case "anilist":
			parsed, _ := strconv.Atoi(req.MatchedID)
			newItem.AniListID = parsed
		case "mangadex":
			newItem.MangaDexID = req.MatchedID
		}
		doc.Trackers = append(doc.Trackers, newItem)
	}

	// Mutate userTrackers in Sanity
	mutations := []map[string]any{}
	if doc.ID == "" {
		mutations = append(mutations, map[string]any{
			"create": map[string]any{
				"_type":     "userTrackers",
				"updatedAt": time.Now().UTC().Format(time.RFC3339),
				"user": map[string]any{
					"_type": "reference",
					"_ref":  authID,
				},
				"trackers": doc.Trackers,
			},
		})
	} else {
		mutations = append(mutations, map[string]any{
			"patch": map[string]any{
				"id": doc.ID,
				"set": map[string]any{
					"updatedAt": time.Now().UTC().Format(time.RFC3339),
					"trackers":  doc.Trackers,
				},
			},
		})
	}

	// If admin/mod selected globalBind, update manga document directly
	if isAdmin && req.GlobalBind {
		patchSet := map[string]any{}
		switch req.Provider {
		case "mal":
			parsed, _ := strconv.Atoi(req.MatchedID)
			patchSet["myAnimeListId"] = parsed
		case "anilist":
			parsed, _ := strconv.Atoi(req.MatchedID)
			patchSet["anilistId"] = parsed
		case "mangadex":
			patchSet["mangaDexId"] = req.MatchedID
		}
		if len(patchSet) > 0 {
			mutations = append(mutations, map[string]any{
				"patch": map[string]any{
					"id":  req.SeriesID,
					"set": patchSet,
				},
			})
		}
	}

	_, err = MutateSanity(ctx, mutations)
	return err
}

// UnbindTracker unbinds a specific provider for a series
func (s *TrackerSyncService) UnbindTracker(ctx context.Context, logtoID, seriesID, provider string) error {
	doc, err := s.GetUserTrackersDoc(ctx, logtoID)
	if err != nil || doc.ID == "" {
		return err
	}

	for i := range doc.Trackers {
		if doc.Trackers[i].SeriesID == seriesID {
			switch provider {
			case "mal":
				doc.Trackers[i].MALID = 0
			case "anilist":
				doc.Trackers[i].AniListID = 0
			case "mangadex":
				doc.Trackers[i].MangaDexID = ""
			}
			doc.Trackers[i].LastSyncedAt = time.Now().UTC().Format(time.RFC3339)
			break
		}
	}

	mutations := []map[string]any{
		{
			"patch": map[string]any{
				"id": doc.ID,
				"set": map[string]any{
					"updatedAt": time.Now().UTC().Format(time.RFC3339),
					"trackers":  doc.Trackers,
				},
			},
		},
	}
	_, err = MutateSanity(ctx, mutations)
	return err
}

// UpdateSingleTracker updates user's status/chapter/score for a single provider
func (s *TrackerSyncService) UpdateSingleTracker(ctx context.Context, req models.UpdateTrackerStatusRequest) error {
	_, _, integrations, err := getUserAuthData(ctx, req.LogtoID)
	if err != nil {
		return err
	}

	doc, _ := s.GetUserTrackersDoc(ctx, req.LogtoID)
	globalMAL, globalAniList, globalMD, _ := getSeriesMangaInfo(ctx, req.SeriesID)

	targetMAL := globalMAL
	targetAniList := globalAniList
	targetMD := globalMD

	if doc != nil {
		for _, t := range doc.Trackers {
			if t.SeriesID == req.SeriesID {
				if t.MALID > 0 {
					targetMAL = t.MALID
				}
				if t.AniListID > 0 {
					targetAniList = t.AniListID
				}
				if t.MangaDexID != "" {
					targetMD = t.MangaDexID
				}
				break
			}
		}
	}

	switch req.Provider {
	case "mal":
		if integrations.MAL.AccessToken == "" || targetMAL == 0 {
			return fmt.Errorf("MAL hesabi bagli degil veya seri eslesmemis")
		}
		return s.malService.UpdateUserMangaStatus(ctx, integrations.MAL.AccessToken, targetMAL, req.Status, req.Chapter, req.Score)
	case "anilist":
		if integrations.AniList.AccessToken == "" {
			return fmt.Errorf("AniList hesabi bagli degil")
		}
		if targetAniList == 0 && targetMAL > 0 {
			if resolvedID, err := s.aniListService.ResolveMediaIDByMALID(ctx, targetMAL); err == nil && resolvedID > 0 {
				targetAniList = resolvedID
			}
		}
		if targetAniList == 0 {
			return fmt.Errorf("AniList hesabi bagli degil veya seri eslesmemis")
		}
		return s.aniListService.UpdateUserMangaStatus(ctx, integrations.AniList.AccessToken, targetAniList, req.Status, req.Chapter, req.Score)
	case "mangadex":
		if integrations.MangaDex.AccessToken == "" || targetMD == "" {
			return fmt.Errorf("MangaDex hesabi bagli degil veya seri eslesmemis")
		}
		return s.mangaDexService.UpdateUserMangaStatus(ctx, integrations.MangaDex.AccessToken, targetMD, req.Status, req.Score)
	}

	return fmt.Errorf("bilinmeyen saglayici: %s", req.Provider)
}

// SyncSeriesProgress propagates chapter progress to all bound trackers concurrently
func (s *TrackerSyncService) SyncSeriesProgress(ctx context.Context, req models.SyncProgressRequest) {
	// Execute in background goroutine to never block the caller
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_, _, integrations, err := getUserAuthData(bgCtx, req.LogtoID)
		if err != nil {
			return
		}

		doc, _ := s.GetUserTrackersDoc(bgCtx, req.LogtoID)
		globalMAL, globalAniList, globalMD, _ := getSeriesMangaInfo(bgCtx, req.SeriesID)

		targetMAL := globalMAL
		targetAniList := globalAniList
		targetMD := globalMD
		autoSync := true

		if doc != nil {
			for i := range doc.Trackers {
				if doc.Trackers[i].SeriesID == req.SeriesID {
					if doc.Trackers[i].MALID > 0 {
						targetMAL = doc.Trackers[i].MALID
					}
					if doc.Trackers[i].AniListID > 0 {
						targetAniList = doc.Trackers[i].AniListID
					}
					if doc.Trackers[i].MangaDexID != "" {
						targetMD = doc.Trackers[i].MangaDexID
					}
					autoSync = doc.Trackers[i].AutoSync
					break
				}
			}
		}

		if !autoSync {
			return
		}

		// Determine status: if chapter >= totalChapters, completed; else reading
		status := "reading"
		if req.TotalChapters > 0 && req.ChapterNumber >= float64(req.TotalChapters) {
			status = "completed"
		}

		var wg sync.WaitGroup

		// MAL Sync
		if integrations.MAL.AccessToken != "" && targetMAL > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := s.malService.UpdateUserMangaStatus(bgCtx, integrations.MAL.AccessToken, targetMAL, status, req.ChapterNumber, -1); err != nil {
					log.Printf("[Sync MAL Hatasi]: %v", err)
				}
			}()
		}

		// AniList Sync
		if integrations.AniList.AccessToken != "" && (targetAniList > 0 || targetMAL > 0) {
			effectiveAniList := targetAniList
			if effectiveAniList == 0 && targetMAL > 0 {
				if resolvedID, err := s.aniListService.ResolveMediaIDByMALID(bgCtx, targetMAL); err == nil && resolvedID > 0 {
					effectiveAniList = resolvedID
				}
			}
			if effectiveAniList > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := s.aniListService.UpdateUserMangaStatus(bgCtx, integrations.AniList.AccessToken, effectiveAniList, status, req.ChapterNumber, -1); err != nil {
						log.Printf("[Sync AniList Hatasi]: %v", err)
					}
				}()
			}
		}

		// MangaDex Sync
		if integrations.MangaDex.AccessToken != "" && targetMD != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := s.mangaDexService.UpdateUserMangaStatus(bgCtx, integrations.MangaDex.AccessToken, targetMD, status, -1); err != nil {
					log.Printf("[Sync MangaDex Hatasi]: %v", err)
				}
			}()
		}

		wg.Wait()

		// Update lastSyncedChapter in Sanity userTrackers
		if doc != nil && doc.ID != "" {
			for i := range doc.Trackers {
				if doc.Trackers[i].SeriesID == req.SeriesID {
					doc.Trackers[i].LastChapter = req.ChapterNumber
					doc.Trackers[i].LastStatus = status
					doc.Trackers[i].LastSyncedAt = time.Now().UTC().Format(time.RFC3339)
					break
				}
			}
			mutations := []map[string]any{
				{
					"patch": map[string]any{
						"id": doc.ID,
						"set": map[string]any{
							"updatedAt": time.Now().UTC().Format(time.RFC3339),
							"trackers":  doc.Trackers,
						},
					},
				},
			}
			_, _ = MutateSanity(bgCtx, mutations)
		}
	}()
}

// SaveUserIntegration saves provider tokens into Sanity auth document
func (s *TrackerSyncService) SaveUserIntegration(ctx context.Context, logtoID string, patchSet map[string]any) error {
	authID, _, _, err := getUserAuthData(ctx, logtoID)
	if err != nil {
		return err
	}
	mutations := []map[string]any{
		{
			"patch": map[string]any{
				"id":  authID,
				"set": patchSet,
			},
		},
	}
	_, err = MutateSanity(ctx, mutations)
	return err
}

// DisconnectUserIntegration unsets the provider tokens in Sanity auth document
func (s *TrackerSyncService) DisconnectUserIntegration(ctx context.Context, logtoID, provider string) error {
	authID, _, _, err := getUserAuthData(ctx, logtoID)
	if err != nil {
		return err
	}
	mutations := []map[string]any{
		{
			"patch": map[string]any{
				"id": authID,
				"unset": []string{
					fmt.Sprintf("integrations.%s", provider),
				},
			},
		},
	}
	_, err = MutateSanity(ctx, mutations)
	return err
}

