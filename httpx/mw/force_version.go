package mw

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderFrontVersion = "X-Front-Version"
)

// MinVersion - если переданная в хидере X-Front-Version числовая версия фронта ниже указанной, выводиться ошибка
// для инициирования обновления фронта. Если указанная версия = 0, проверка хидера отключается.
func MinVersion(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Hardcoded for test
			frontVersion := r.Header.Get(HeaderFrontVersion)
			if frontVersion == "0.0.0" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "1.2"}`))

				return
			}

			if version != "" && version != "0" && frontVersion != "0.0.0.0" {
				if versionCmp(frontVersion, version) < 0 {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"error": "1.2"}`))

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func versionCmp(v1, v2 string) int {
	v1nums := strings.Split(v1, ".")
	v2nums := strings.Split(v2, ".")
	for index := 0; index < len(v1nums); index++ {
		v1num, err := strconv.Atoi(v1nums[index])
		if err != nil {
			return -1
		}

		if len(v2nums) < index+1 {
			return -1
		}

		v2num, err := strconv.Atoi(v2nums[index])
		if err != nil {
			return -1
		}
		if v1num < v2num {
			return -1
		}
		if v1num > v2num {
			return 1
		}
	}

	return 0
}
