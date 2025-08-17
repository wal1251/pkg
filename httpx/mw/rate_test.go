package mw_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/wal1251/pkg/httpx/mw"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	limit := 1
	per := time.Second

	rateLimiter := mw.RateLimiter(limit, per)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	// Первый запрос
	rateLimiter(mockHandler).ServeHTTP(rr, req)
	// Проверка на 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Еше раз, чтобы проверить.
	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req)

	// Проверка на 429 Too Many Requests
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimiter_AfterDuration(t *testing.T) {
	limit := 1
	per := time.Second

	rateLimiter := mw.RateLimiter(limit, per)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	rateLimiter(mockHandler).ServeHTTP(rr, req)
	// Проверка на 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Sleep пока истечет срок ограничения.
	time.Sleep(per)

	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req)

	// Проверка на 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limit := 1
	per := time.Second

	rateLimiter := mw.RateLimiter(limit, per)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Создаем запрос для первого IP
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "192.168.1.1"

	// Создаем запрос для второго IP
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "192.168.1.2"

	rr := httptest.NewRecorder()

	rateLimiter(mockHandler).ServeHTTP(rr, req1)

	// Проверка на 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req2)

	// Проверка на 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Снова отправляем запрос для первого IP
	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req1)

	// Проверка на 429 Too Many Requests
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Снова отправляем запрос для второго IP
	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req2)

	// Проверка на 429 Too Many Requests
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestGlobalRateLimiter(t *testing.T) {
	// Имитация HTTP-запросов с RateLimiter
	handler := mw.RateLimiter(2, 1*time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Создаем запросы
	req1 := httptest.NewRequest("GET", "/test", nil)
	req2 := httptest.NewRequest("GET", "/test", nil)
	req3 := httptest.NewRequest("GET", "/test", nil)

	// Мокируем ResponseWriter
	recorder1 := httptest.NewRecorder()
	recorder2 := httptest.NewRecorder()
	recorder3 := httptest.NewRecorder()

	// Первая попытка должна быть успешной
	handler.ServeHTTP(recorder1, req1)
	assert.Equal(t, http.StatusOK, recorder1.Code)

	// Вторая попытка тоже должна быть успешной
	handler.ServeHTTP(recorder2, req2)
	assert.Equal(t, http.StatusOK, recorder2.Code)

	// Третья попытка должна вернуть ошибку (Too Many Requests)
	handler.ServeHTTP(recorder3, req3)
	assert.Equal(t, http.StatusTooManyRequests, recorder3.Code)
}

func TestRateLimiter_DifferentPaths(t *testing.T) {
	limit := 1
	per := time.Second

	rateLimiter := mw.RateLimiter(limit, per)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req1 := httptest.NewRequest("GET", "/path1", nil)
	req1.RemoteAddr = "192.168.1.1"

	req2 := httptest.NewRequest("GET", "/path2", nil)
	req2.RemoteAddr = "192.168.1.1"

	rr := httptest.NewRecorder()

	rateLimiter(mockHandler).ServeHTTP(rr, req1)

	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req2)

	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req1)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	rr = httptest.NewRecorder()
	rateLimiter(mockHandler).ServeHTTP(rr, req2)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimiter_Concurrently(t *testing.T) {
	rateLimiter := mw.NewIPRateLimiter(rate.Every(1*time.Second), 10)

	ip := "192.168.1.1"
	rateLimiter.AddIP(ip)
	var (
		successCount, failureCount int
		wg                         sync.WaitGroup
		mutex                      sync.Mutex
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter := rateLimiter.GetLimiter(ip)
			if limiter.Allow() {
				mutex.Lock()
				defer mutex.Unlock()
				successCount++
			} else {
				mutex.Lock()
				defer mutex.Unlock()
				failureCount++
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, 10, successCount)
	assert.Equal(t, 90, failureCount)
}
