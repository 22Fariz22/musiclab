package apilyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LyricsAPI struct {
	Lyrics string `json:"lyrics"`
	Verses string
}

// FetchLyrics идем в АПИ получать текст песни
func FetchLyrics(ctx context.Context, APIUrl string, artist, title string) (LyricsAPI, error) {
	url := fmt.Sprintf(APIUrl, artist, title)

	// Создаем HTTP-запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return LyricsAPI{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LyricsAPI{}, fmt.Errorf("failed to fetch lyrics: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return LyricsAPI{}, fmt.Errorf("lyrics not found, status code: %d", resp.StatusCode)
	}

	// Декодируем JSON-ответ
	var lyricsResponse LyricsAPI
	if err := json.NewDecoder(resp.Body).Decode(&lyricsResponse); err != nil {
		return LyricsAPI{}, fmt.Errorf("failed to decode response: %w", err)
	}

	//делим на куплеты
	verses := prepareLyrics(lyricsResponse.Lyrics)
	lyricsResponse.Verses = verses

	return lyricsResponse, nil
}

// prepareLyrics делит текст песни на куплеты
func prepareLyrics(lyrics string) string {
	// Разделяем текст на строки
	lines := strings.Split(lyrics, "\n")

	// Сохраняем только непустые строки
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	// Объединяем строки с использованием "\n"
	return strings.Join(cleanedLines, "\\n")
}
