package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/internal/models"
	"github.com/22Fariz22/musiclab/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type lyricsUseCase struct {
	cfg         *config.Config
	lyricsRepo  lyrics.Repository
	redisClient *redis.Client
	logger      logger.Logger
	httpClient  *http.Client
}

func NewLyricsUseCase(cfg *config.Config, lyricsRepo lyrics.Repository, redisClient *redis.Client, logger logger.Logger) lyrics.UseCase {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &lyricsUseCase{
		cfg:         cfg,
		lyricsRepo:  lyricsRepo,
		redisClient: redisClient,
		logger:      logger,
		httpClient:  httpClient,
	}
}

// Ping check
func (u lyricsUseCase) Ping() error {
	u.logger.Debug("Call UseCase Ping()\n")

	err := u.lyricsRepo.Ping()
	if err != nil {
		u.logger.Debug("error in lyricsUseCase Ping()")
		return err
	}

	u.logger.Debug("Pong in lyricsUseCase Ping()")
	return nil
}

func (u lyricsUseCase) DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error {
	u.logger.Debugf("in usecase DeleteSongByGroupAndTrack. Deleting song. Group: %s, Track: %s\n", groupName, trackName)
	return u.lyricsRepo.DeleteSongByGroupAndTrack(ctx, groupName, trackName)
}

func (u lyricsUseCase) UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error {
	u.logger.Debugf("in usecase UpdateTrackByID() ID:%d", updateData.ID)
	return u.lyricsRepo.UpdateTrackByID(ctx, updateData)
}

func (u lyricsUseCase) CreateTrack(ctx context.Context, songRequest models.SongRequest) (models.SongDetail, error) {
	u.logger.Debug("in usecase CreateTrack()\n")
	ctx, cancel := context.WithTimeout(ctx, u.cfg.API.APICtxTimeout)
	defer cancel()

	fullURL, err := u.BuildAPIURL(songRequest.Group, songRequest.Song)
	if err != nil {
		u.logger.Errorf("failed to build API URL: %v", err)
		return models.SongDetail{}, err
	}

	maxRetries := u.cfg.API.MaxRetries // Максимальное количество попыток
	retryDelay := u.cfg.API.RetryDelay // Задержка между попытками

	var songDetails models.SongDetail
	var lastErr error

	// Логика повторных попыток
	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			u.logger.Warnf("context cancelled during fetch: %v", ctx.Err())
			return models.SongDetail{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			songDetails, err = u.FetchAPI(ctx, fullURL)
			if err == nil {
				break
			}
			lastErr = err
			u.logger.Warnf("attempt %d/%d failed to fetch lyrics: %v", attempt, maxRetries, err)
			time.Sleep(retryDelay)
		}
	}

	if lastErr != nil {
		u.logger.Errorf("all attempts to fetch lyrics failed: %v", lastErr)
		return models.SongDetail{}, fmt.Errorf("failed to fetch lyrics: %w", lastErr)
	}

	if err := u.lyricsRepo.CreateTrack(ctx, songRequest, songDetails); err != nil {
		u.logger.Errorf("failed to save track: %v", err)
		return songDetails, fmt.Errorf("saving track: %w", err)
	}

	u.logger.Infof("track created successfully: %+v", songDetails)
	return songDetails, nil
}

func (u lyricsUseCase) BuildAPIURL(group, song string) (string, error) {
	APIAddr := fmt.Sprintf("%s:%s%s", u.cfg.Server.BaseUrl, u.cfg.Server.Port, u.cfg.API.APIPath)

	parsedURL, err := url.Parse(APIAddr)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}

	q := parsedURL.Query()
	q.Set("group", group)
	q.Set("song", song)
	parsedURL.RawQuery = q.Encode()

	u.logger.Debugf("APIAddr:%s", APIAddr)
	u.logger.Debugf("parsedURL:%s", parsedURL.String())

	return parsedURL.String(), nil
}

// FetchAPI обращение к API
func (u *lyricsUseCase) FetchAPI(ctx context.Context, url string) (models.SongDetail, error) {
	u.logger.Debugf("UrlAPI: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		u.logger.Errorf("failed to create request: %v", err)
		return models.SongDetail{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		u.logger.Errorf("request failed: %v", err)
		return models.SongDetail{}, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.SongDetail{}, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		u.logger.Errorf("failed to decode response: %v", err)
		return models.SongDetail{}, fmt.Errorf("decoding response: %w", err)
	}

	return songDetail, nil

	// return models.SongDetail{
	// 	ReleaseDate: "01.03.89",
	// 	Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
	// 	Link:        "yputube.com/sd2ff6ggf",
	// }, nil
}

// GetSongVerseByPage
func (u lyricsUseCase) GetSongVerseByID(ctx context.Context, id uint, page int) (string, error) {
	u.logger.Debugf("in UC GetSongVerseByPage ID:%d, page:%d\n", id, page)

	cacheKey := fmt.Sprintf("song:%d", id)

	// Проверяем кэш
	cachedSong, err := u.redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		u.logger.Errorf("Error fetching from Redis: %v", err)
		cachedSong = ""
	}

	var songText string
	if cachedSong == "" {
		u.logger.Debugf("Cache miss for key: %s. Fetching from database.", cacheKey)

		// Если в кэше ничего нет, идём в базу данных
		song, err := u.lyricsRepo.GetSongByID(ctx, id)
		if err != nil {
			u.logger.Debugf("error in uc u.lyricsRepo.GetSongByID():%v", err)
			return "", fmt.Errorf("failed to get song from database: %w", err)
		}

		songText = song.Text

		// Сохраняем песню в кэше
		err = u.redisClient.Set(ctx, cacheKey, songText, u.cfg.Redis.SongTextCasheTTL).Err()
		if err != nil {
			u.logger.Errorf("Error caching song in Redis: %v", err)
		}
	} else {
		u.logger.Debugf("Cache hit for key: %s", cacheKey)
		songText = cachedSong
	}

	u.logger.Debugf("text song:", songText)

	// Разделяем текст на куплеты
	verses := prepareLyrics(songText)

	// Проверяем, существует ли куплет для указанной страницы
	if page <= 0 || page > len(verses) {
		return "", fmt.Errorf("no verse available for page %d", page)
	}

	// Возвращаем куплет по индексу (page - 1, так как индексация с 0)
	return verses[page-1], nil
}

// prepareLyrics делим песню на куплеты
func prepareLyrics(lyrics string) []string {
	// lines := strings.Split(lyrics, "\\\n")
	// // lines = strings.Split(lyrics, "\n")
	// return lines

	// Разделяем текст по строкам
	lines := strings.Split(lyrics, "\n")

	var verses []string
	var currentVerse []string

	for _, line := range lines {
		// Если строка пустая, завершаем текущий куплет
		if strings.TrimSpace(line) == "" {
			if len(currentVerse) > 0 {
				verses = append(verses, strings.Join(currentVerse, "\n"))
				currentVerse = []string{}
			}
		} else {
			// Добавляем строку к текущему куплету
			currentVerse = append(currentVerse, line)
		}
	}

	// Добавляем последний куплет, если он не пуст
	if len(currentVerse) > 0 {
		verses = append(verses, strings.Join(currentVerse, "\n"))
	}

	return verses
}

func (u lyricsUseCase) GetLibrary(ctx context.Context, group, song, text, releaseDate string, page, limit int) ([]models.Song, int, error) {
	u.logger.Debugf("Fetching library with filters: group=%s, song=%s,text=%s, releaseDate=%s, page=%d, limit=%d", group, song, releaseDate, page, limit)

	offset := (page - 1) * limit

	songs, total, err := u.lyricsRepo.GetLibrary(ctx, group, song, text, releaseDate, offset, limit)
	if err != nil {
		u.logger.Errorf("Error fetching library from repository: %v", err)
		return nil, 0, err
	}

	return songs, total, nil
}
