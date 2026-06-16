package sites

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/time/rate"

	"github.com/thecoretg/tctg-go/threatdown"
)

const requestsPerMinute = 360

var limiter = rate.NewLimiter(rate.Limit(requestsPerMinute/60.0), 90)

type FetchOptions struct {
	SkipSubs   bool
	SkipAddOns bool
}

type FullSite struct {
	threatdown.Site
	Subs   []threatdown.SiteSubscription `json:"subs,omitempty"`
	AddOns []threatdown.AddOn            `json:"add_ons,omitempty"`
	Errors FetchErrors                   `json:"errors"`
}

// FetchErrors holds non-fatal errors encountered while fetching a site's
// subscriptions or add-ons. A nil field means that fetch succeeded.
type FetchErrors struct {
	SubFetchError   *string `json:"sub_fetch_error"`
	AddOnFetchError *string `json:"addon_fetch_error"`
}

func FetchSites(ctx context.Context, opts FetchOptions) ([]FullSite, error) {
	c, err := threatdown.NewClientFromEnv(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating threatdown client: %w", err)
	}

	return fetchAllSites(ctx, c, opts)
}

func fetchAllSites(ctx context.Context, c *threatdown.Client, opts FetchOptions) ([]FullSite, error) {
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("fetching sites: %w", err)
	}

	raw, err := c.ListSites(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("fetching sites: %w", err)
	}

	var wg sync.WaitGroup
	results := make(chan FullSite, len(raw))

	for _, s := range raw {
		wg.Go(func() {
			fs := FullSite{Site: s}

			if !opts.SkipSubs {
				subs, err := fetchSubs(ctx, c, s)
				if err != nil {
					msg := err.Error()
					fs.Errors.SubFetchError = &msg
				} else {
					fs.Subs = subs
				}
			}

			if !opts.SkipAddOns {
				addOns, err := fetchAddOns(ctx, c, s)
				if err != nil {
					msg := err.Error()
					fs.Errors.AddOnFetchError = &msg
				} else {
					fs.AddOns = addOns
				}
			}

			results <- fs
		})
	}

	wg.Wait()
	close(results)

	var fullSites []FullSite
	for fs := range results {
		fullSites = append(fullSites, fs)
	}

	return fullSites, nil
}

func fetchSubs(ctx context.Context, c *threatdown.Client, site threatdown.Site) ([]threatdown.SiteSubscription, error) {
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	subs, err := c.GetSiteSubscriptions(ctx, site.ID)
	if errors.Is(err, threatdown.ErrSubNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting site subs: %w", err)
	}

	return subs, nil
}

func fetchAddOns(ctx context.Context, c *threatdown.Client, site threatdown.Site) ([]threatdown.AddOn, error) {
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	addOns, err := c.ListAddOns(ctx, site.ID)
	if errors.Is(err, threatdown.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting site add-ons: %w", err)
	}

	return addOns, nil
}
