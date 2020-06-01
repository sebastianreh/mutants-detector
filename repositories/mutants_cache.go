package repositories

import (
	"github.com/patrickmn/go-cache"
	"github.com/sebastianreh/mutants-detector/models"
	log "github.com/sirupsen/logrus"
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

// Crea el cache

func NewMutantCacheRepository() MutantCache {
	return MutantCache{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Guarda los stats en el cache

func (c MutantCache) SaveStatsInCache(stats *models.MutantsStats) {
	c.cache.Set(cacheStats, stats, cache.NoExpiration)
}

// Retorna los stats del cache

func (c MutantCache) GetStatsFromCache() *models.MutantsStats {
	cacheStats, ok := c.cache.Get(cacheStats)
	if !ok {
		log.Info("repositories.GetStatsFromCache | No stats cache in memory")
		return nil
	}
	stats := cacheStats.(*models.MutantsStats)
	return stats
}
