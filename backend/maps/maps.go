package maps

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	gmaps "googlemaps.github.io/maps"
)

// Client interacts with map services
type Client struct {
	l  *zap.SugaredLogger
	gm *gmaps.Client

	cache *cache
}

// NewClient instantiates a maps client
func NewClient(l *zap.SugaredLogger, key string) (*Client, error) {
	gm, err := gmaps.NewClient(gmaps.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	l.Info("attempting to connect to gmaps")
	if _, _, err := gm.Directions(context.Background(), &maps.DirectionsRequest{
		Origin:      "Vancouver",
		Destination: "Surrey",
	}); err != nil {
		l.Errorw("failed to connect to google maps", "error", err)
		return nil, fmt.Errorf("failed to connect to google maps: %s", err.Error())
	}
	l.Info("successfully made query to gmaps")
	return &Client{
		l:     l,
		gm:    gm,
		cache: newCache(5*time.Minute, 5*time.Minute),
	}, nil
}

func (c *Client) PointsOfInterest(ctx context.Context) {
	c.gm.TextSearch(ctx, &gmaps.TextSearchRequest{})
}

// Close stops background jobs
func (c *Client) Close() { c.cache.stop <- true }
