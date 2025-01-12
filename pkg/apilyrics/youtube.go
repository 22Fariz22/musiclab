package apilyrics

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func GetYoutubeLink(URL, query string) string {
	// Формируем URL для поискового запроса
	searchURL := fmt.Sprintf(URL, url.QueryEscape(query))

	// Выполняем HTTP-запрос
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalf("Failed to fetch YouTube search results: %v", err)
	}
	defer resp.Body.Close()

	// Читаем HTML содержимое страницы
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Регулярное выражение для поиска ссылки на видео
	re := regexp.MustCompile(`\/watch\?v=[\w-]{11}`)
	matches := re.FindAllString(string(body), -1)

	var videoURL string
	// Проверяем, найдены ли ссылки
	if len(matches) > 0 {
		// Берём первую уникальную ссылку
		videoID := matches[0]
		videoURL = fmt.Sprintf("https://www.youtube.com%s", videoID)
		fmt.Printf("First video: %s\n", videoURL)
	} else {
		fmt.Println("No video found")
		return videoURL
	}
	return videoURL
}
