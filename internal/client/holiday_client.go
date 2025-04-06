package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Holiday struct {
	Date        string `json:"date"`
	LocalName   string `json:"localName"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	Fixed       bool   `json:"fixed"`
	Global      bool   `json:"global"`
	Type        string `json:"type"`
}

type cachedHolidays struct {
	Holidays []Holiday
	CachedAt time.Time
}

// Cache structure
var (
	holidayCache   = make(map[string]map[int]cachedHolidays)
	holidayCacheMu sync.RWMutex
	cacheTTL       = 24 * time.Hour // You can make this configurable
)

func IsHoliday(countryCode string, date time.Time) (bool, error) {
	year := date.Year()

	// Try cache first
	holidayCacheMu.RLock()
	yearMap, ok := holidayCache[countryCode]
	cacheEntry, yearFound := yearMap[year]
	holidayCacheMu.RUnlock()

	cacheValid := ok && yearFound && time.Since(cacheEntry.CachedAt) < cacheTTL

	var holidays []Holiday
	if cacheValid {
		holidays = cacheEntry.Holidays
	} else {
		// Cache miss or expired
		fetchedHolidays, err := fetchHolidaysFromAPI(countryCode, year)
		if err != nil {
			return false, err
		}

		// Save to cache
		holidayCacheMu.Lock()
		if _, exists := holidayCache[countryCode]; !exists {
			holidayCache[countryCode] = make(map[int]cachedHolidays)
		}
		holidayCache[countryCode][year] = cachedHolidays{
			Holidays: fetchedHolidays,
			CachedAt: time.Now(),
		}
		holidayCacheMu.Unlock()

		holidays = fetchedHolidays
	}

	targetDate := date.Format("2006-01-02")
	for _, h := range holidays {
		if h.Date == targetDate {
			return true, nil
		}
	}

	return false, nil
}

func fetchHolidaysFromAPI(countryCode string, year int) ([]Holiday, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s", year, countryCode)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch holidays: %s", resp.Status)
	}

	var holidays []Holiday
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}
