package repositories

import (
	"ExamenMeLiMutante/models"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	MutantCache struct {
		IMutantCache
		cache *cache.Cache
	}

	IMutantCache interface {
		SaveStatsInCache(stats *models.MutantsStats)
		GetStatsFromCache() *models.MutantsStats
	}
)

const (
	cacheStats = "stats"
)

var (
	expirationTime  = 10 * time.Minute
	cleanUpInterval = 1 * time.Minute
)

// Crea el cache

func NewMutantCacheRepository() MutantCache {
	return MutantCache{
		cache: cache.New(expirationTime, cleanUpInterval),
	}
}

// Guarda los stats en el cache

func (c MutantCache) SaveStatsInCache(stats *models.MutantsStats) {
	c.cache.Set(cacheStats, stats, expirationTime)
}

// Retorna los stats del cache

func (c MutantCache) GetStatsFromCache() (*models.MutantsStats) {
	cacheStats, ok := c.cache.Get(cacheStats)
	if !ok {
		log.Info("repositories.GetStatsFromCache | No stats cache in memory")
	}
	stats := cacheStats.(*models.MutantsStats)
	return stats
}
