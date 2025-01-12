package http_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	lyricsHTTP "github.com/22Fariz22/musiclab/internal/lyrics/delivery/http"
	mockUC "github.com/22Fariz22/musiclab/internal/lyrics/mock"
	"github.com/22Fariz22/musiclab/internal/models"
	mockLogger "github.com/22Fariz22/musiclab/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	mockTestify "github.com/stretchr/testify/mock"
)

func TestPingHandler(t *testing.T) {
	e := echo.New()

	mockUseCase := new(mockUC.MockUseCase)
	mockUseCase.On("Ping").Return(nil)

	appLogger := mockLogger.CreateTestLogger()

	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Ping()(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `"pong"`, strings.TrimSpace(rec.Body.String())) // Исправлено здесь
	mockUseCase.AssertExpectations(t)
}

func TestDeleteSongByGroupAndTrack(t *testing.T) {
	e := echo.New()

	// Создаем мок UseCase
	mockUseCase := new(mockUC.MockUseCase)

	appLogger := mockLogger.CreateTestLogger()

	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	t.Run("Successful deletion", func(t *testing.T) {
		mockUseCase.On("DeleteSongByGroupAndTrack", mockTestify.Anything, "Muse", "Uprising").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/delete?group=Muse&track=Uprising", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.DeleteSongByGroupAndTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `"Track is deleted"`, strings.TrimSpace(rec.Body.String()))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Track not found", func(t *testing.T) {
		mockUseCase.On("DeleteSongByGroupAndTrack", mockTestify.Anything, "Muse", "UnknownSong").Return(sql.ErrNoRows).Once()

		req := httptest.NewRequest(http.MethodDelete, "/delete?group=Muse&track=UnknownSong", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.DeleteSongByGroupAndTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, `"Track not found"`, strings.TrimSpace(rec.Body.String()))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockUseCase.On("DeleteSongByGroupAndTrack", mockTestify.Anything, "Muse", "Uprising").Return(errors.New("unexpected error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/delete?group=Muse&track=Uprising", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.DeleteSongByGroupAndTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, `"Failed to delete song"`, strings.TrimSpace(rec.Body.String()))
		mockUseCase.AssertExpectations(t)
	})
}

func TestUpdateTrackByID(t *testing.T) {
	e := echo.New()

	mockUseCase := new(mockUC.MockUseCase)

	appLogger := mockLogger.CreateTestLogger()

	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	t.Run("Successful update", func(t *testing.T) {
		updateData := models.UpdateTrackRequest{
			ID:        1,
			GroupName: "Muse",
			SongName:  "Uprising",
			Text:      ptr("Updated lyrics text"),
		}

		mockUseCase.On("UpdateTrackByID", mockTestify.Anything, updateData).Return(nil)

		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.UpdateTrackByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message":"track updated successfully"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(`invalid-json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.UpdateTrackByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error":"invalid JSON format"}`, rec.Body.String())
	})

	t.Run("Song not found", func(t *testing.T) {
		updateData := models.UpdateTrackRequest{
			ID:        1,
			GroupName: "Muse",
			SongName:  "Unknown Song",
		}

		mockUseCase.On("UpdateTrackByID", mockTestify.Anything, updateData).Return(sql.ErrNoRows)

		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.UpdateTrackByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error":"song not found"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		updateData := models.UpdateTrackRequest{
			ID:        1,
			GroupName: "Muse",
			SongName:  "Uprising",
		}

		mockUseCase.On("UpdateTrackByID", mockTestify.Anything, updateData).Return(errors.New("unexpected error"))

		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.UpdateTrackByID()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"failed to update song"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})
}

func ptr(s string) *string {
	return &s
}

func TestCreateTrack(t *testing.T) {
	e := echo.New()

	mockUseCase := new(mockUC.MockUseCase)

	appLogger := mockLogger.CreateTestLogger()

	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	t.Run("Successful creation", func(t *testing.T) {
		songRequest := models.SongRequest{
			Group: "Muse",
			Song:  "Uprising",
		}
		requestBody, _ := json.Marshal(songRequest)

		mockUseCase.On("CreateTrack", mockTestify.Anything, songRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message": "Track created successfully"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader([]byte(`{invalid-json}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "Invalid JSON for SongRequest"}`, rec.Body.String())
	})

	t.Run("Validation error", func(t *testing.T) {
		songRequest := models.SongRequest{
			Group: "",
			Song:  "",
		}
		requestBody, _ := json.Marshal(songRequest)

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateTrack()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid JSON fields")
	})
}

func TestGetSongVerseByPage(t *testing.T) {
	e := echo.New()

	// Создаем мок UseCase
	mockUseCase := new(mockUC.MockUseCase)

	appLogger := mockLogger.CreateTestLogger()

	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	t.Run("Successful fetch", func(t *testing.T) {
		id := 1
		page := 1
		expectedVerse := "This is a test verse."

		mockUseCase.On("GetSongVerseByPage", mockTestify.Anything, uint(id), page).Return(expectedVerse, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/verses/%d?page=%d", id, page), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", id))
		c.QueryParams().Set("page", fmt.Sprintf("%d", page))

		err := h.GetSongVerseByPage()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"page": %d, "verse": "%s"}`, page, expectedVerse), rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid song ID", func(t *testing.T) {
		// Сброс предыдущих ожиданий и вызовов
		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		req := httptest.NewRequest(http.MethodGet, "/verses/abc?page=1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("abc")

		err := h.GetSongVerseByPage()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "Invalid song ID"}`, rec.Body.String())
	})

	t.Run("Invalid page number", func(t *testing.T) {
		// Сброс предыдущих ожиданий и вызовов
		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		req := httptest.NewRequest(http.MethodGet, "/verses/1?page=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.GetSongVerseByPage()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "Invalid page number"}`, rec.Body.String())
	})

	t.Run("Verse not found for page=2", func(t *testing.T) {
		// Сброс предыдущих ожиданий и вызовов
		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		id := 1
		page := 2

		// Настройка мока для возврата ошибки для page=2
		mockUseCase.On("GetSongVerseByPage", mockTestify.Anything, uint(id), page).Return("", errors.New("verse not found"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/verses/%d?page=%d", id, page), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", id))

		err := h.GetSongVerseByPage()(c)

		t.Logf("Response code: %d, body: %s", rec.Code, rec.Body.String()) // Логирование для отладки

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error": "verse not found"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Empty text for page=1", func(t *testing.T) {
		// Сброс предыдущих ожиданий и вызовов
		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		id := 1
		page := 1

		// Настройка мока для возврата пустого текста для page=1
		mockUseCase.On("GetSongVerseByPage", mockTestify.Anything, uint(id), page).Return("", nil).Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/verses/%d?page=%d", id, page), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", id))

		err := h.GetSongVerseByPage()(c)

		t.Logf("Response code: %d, body: %s", rec.Code, rec.Body.String()) // Логирование для отладки

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"page": 1, "verse": ""}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

}

func TestGetLibrary(t *testing.T) {
	e := echo.New()
	mockUseCase := new(mockUC.MockUseCase)
	appLogger := mockLogger.CreateTestLogger()
	h := lyricsHTTP.NewLyricsHandler(nil, mockUseCase, appLogger)

	t.Run("Successful response", func(t *testing.T) {
		mockSongs := []models.Song{
			{
				ID:          1,
				GroupName:   "Muse",
				SongName:    "Uprising",
				Text:        "Some text",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				Link:        nil,
				ReleaseDate: time.Time{},
			},
			{
				ID:          2,
				GroupName:   "Muse",
				SongName:    "Madness",
				Text:        "Other text",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				Link:        nil,
				ReleaseDate: time.Time{},
			},
		}

		mockUseCase.On("GetLibrary", mockTestify.Anything, "Muse", "Uprising", "2024-11-26", 1, 10).Return(mockSongs, 2, nil)

		req := httptest.NewRequest(http.MethodGet, "/library?group=Muse&song=Uprising&release_date=2024-11-26&page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetLibrary()(c)

		t.Logf("Response code: %d, body: %s", rec.Code, rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
        "data": [
            {
                "ID": 1,
                "GroupName": "Muse",
                "SongName": "Uprising",
                "Text": "Some text",
                "CreatedAt": "0001-01-01T00:00:00Z",
                "UpdatedAt": "0001-01-01T00:00:00Z",
                "Link": null,
                "ReleaseDate": null
            },
            {
                "ID": 2,
                "GroupName": "Muse",
                "SongName": "Madness",
                "Text": "Other text",
                "CreatedAt": "0001-01-01T00:00:00Z",
                "UpdatedAt": "0001-01-01T00:00:00Z",
                "Link": null,
                "ReleaseDate": null
            }
        ],
        "limit": 10,
        "page": 1,
        "total": 2
    }`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid page parameter", func(t *testing.T) {
		mockUseCase.On("GetLibrary", mockTestify.Anything, "", "", "", 1, 10).
			Return([]models.Song{}, 0, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/library?page=-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetLibrary()(c)

		t.Logf("Response code: %d, body: %s", rec.Code, rec.Body.String())

		expectedResponse := map[string]interface{}{
			"page":  1,
			"limit": 10,
			"total": 0,
			"data":  []map[string]interface{}{},
		}

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, toJSON(expectedResponse), rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockUseCase.On("GetLibrary", mockTestify.Anything, "", "", "", 1, 10).
			Return([]models.Song{}, 0, errors.New("unexpected error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/library", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetLibrary()(c)

		t.Logf("Response code: %d, body: %s", rec.Code, rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "Failed to fetch library"}`, rec.Body.String())
		mockUseCase.AssertExpectations(t)
	})
}

func toJSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

