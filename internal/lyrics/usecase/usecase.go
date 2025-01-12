package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/internal/models"
	"github.com/22Fariz22/musiclab/pkg/apilyrics"
	"github.com/22Fariz22/musiclab/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type lyricsUseCase struct {
	cfg         *config.Config
	lyricsRepo  lyrics.Repository
	redisClient *redis.Client
	logger      logger.Logger
}

func NewLyricsUseCase(cfg *config.Config, lyricsRepo lyrics.Repository, redisClient *redis.Client, logger logger.Logger) lyrics.UseCase {
	return &lyricsUseCase{cfg: cfg, lyricsRepo: lyricsRepo, redisClient: redisClient, logger: logger}
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
	u.logger.Debugf("in usecase UpdateTrackByID() ID:%d", updateData)
	return u.lyricsRepo.UpdateTrackByID(ctx, updateData)
}

func (u lyricsUseCase) CreateTrack(ctx context.Context, song models.SongRequest) error {
	u.logger.Debug("in usecase CreateTrack()\n")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	maxRetries := u.cfg.API.MaxRetries // Максимальное количество попыток
	retryDelay := u.cfg.API.RetryDelay // Задержка между попытками

	var lyrics apilyrics.LyricsAPI
	var err error

	// Логика повторных попыток
	for attempt := 1; attempt <= maxRetries; attempt++ {
		//обращаемся к апи который выдает текст песни
		lyrics, err = apilyrics.FetchLyrics(ctx, u.cfg.API.APILyric, song.Group, song.Song)
		if err == nil {
			// Успешный запрос — выходим из цикла
			break
		}

		// Логируем ошибку и пытаемся повторить, если это временная ошибка
		u.logger.Warnf("Attempt %d: failed to fetch lyrics: %v", attempt, err)

		// Если это последняя попытка, возвращаем ошибку
		if attempt == maxRetries {
			u.logger.Errorf("All attempts to fetch lyrics failed")
			return fmt.Errorf("failed to fetch lyrics after %d attempts: %w", maxRetries, err)
		}

		// Задержка перед следующей попыткой
		time.Sleep(retryDelay)
	}

	//создаем данные для передачи обогащенной информации в репозиторий
	songDetails := models.SongDetail{}
	songDetails.Text = lyrics.Verses

	//добавляем ссылку из ютуба
	youtube := apilyrics.GetYoutubeLink(u.cfg.API.YoutubeMainURL, u.cfg.API.YoutubeURL, fmt.Sprintf("%s %s", song.Group, song.Song))
	fmt.Println("youtubeURL ", youtube)
	songDetails.Link = youtube
	songDetails.ReleaseDate = time.Now()

	// Вывод текста песни в дебаг лог
	u.logger.Debugf("Fetched lyrics successfully: %s", songDetails.Text)

	err = u.lyricsRepo.CreateTrack(ctx, song, songDetails)
	if err != nil {
		u.logger.Errorf("Failed to save track in repository: %v", err)
		return fmt.Errorf("failed to save track: %w", err)
	}

	u.logger.Debugf("In usecase in CreateTrack() successfully created: %+v", songDetails)
	return nil
}

// GetSongVerseByPage
func (u lyricsUseCase) GetSongVerseByPage(ctx context.Context, id uint, page int) (string, error) {
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
			u.logger.Debugf("error in uc u.lyricsRepo.GetSongByID():", err)
			return "", fmt.Errorf("failed to get song from database: %w", err)
		}

		songText = song.Text

		// Сохраняем песню в кэше
		err = u.redisClient.Set(ctx, cacheKey, songText, 6*time.Hour).Err()
		if err != nil {
			u.logger.Errorf("Error caching song in Redis: %v", err)
		}
	} else {
		u.logger.Debugf("Cache hit for key: %s", cacheKey)
		songText = cachedSong
	}
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
	lines := strings.Split(lyrics, "\\n")
	return lines
}

func (u lyricsUseCase) GetLibrary(ctx context.Context, group, song, releaseDate string, page, limit int) ([]models.Song, int, error) {
	u.logger.Debugf("Fetching library with filters: group=%s, song=%s, releaseDate=%s, page=%d, limit=%d", group, song, releaseDate, page, limit)

	offset := (page - 1) * limit

	songs, total, err := u.lyricsRepo.GetLibrary(ctx, group, song, releaseDate, offset, limit)
	if err != nil {
		u.logger.Errorf("Error fetching library from repository: %v", err)
		return nil, 0, err
	}

	return songs, total, nil
}

