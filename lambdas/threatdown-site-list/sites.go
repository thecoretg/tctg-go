package sites

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/thecoretg/tctg-go/threatdown"
)

type FullSite struct {
	threatdown.Site
	Subs   []threatdown.SiteSubscription `json:"subs"`
	AddOns []threatdown.AddOn            `json:"add_ons"`
}

func FetchSites(ctx context.Context) ([]FullSite, error) {
	c, err := createClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating threatdown client: %w", err)
	}

	return fetchAllSites(ctx, c)
}

func createClient(ctx context.Context) (*threatdown.Client, error) {
	cfg := threatdown.Config{
		ClientID:     os.Getenv("THREATDOWN_CLIENT_ID"),
		ClientSecret: os.Getenv("THREATDOWN_CLIENT_SECRET"),
	}

	return threatdown.NewClient(ctx, cfg)
}

func fetchAllSites(ctx context.Context, c *threatdown.Client) ([]FullSite, error) {
	raw, err := c.ListSites(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("fetching sites: %w", err)
	}

	var wg sync.WaitGroup
	fullSites := make(chan FullSite, len(raw))

	for _, s := range raw {
		wg.Go(func() {
			fullSites <- FullSite{
				Site:   s,
				Subs:   fetchSubs(ctx, c, s),
				AddOns: fetchAddOns(ctx, c, s),
			}
		})
	}

	wg.Wait()
	close(fullSites)

	results := make([]FullSite, 0, len(raw))
	for fs := range fullSites {
		results = append(results, fs)
	}

	return results, nil
}

func fetchSubs(ctx context.Context, c *threatdown.Client, site threatdown.Site) []threatdown.SiteSubscription {
	subs, err := c.GetSiteSubscriptions(ctx, site.ID)
	if err != nil {
		slog.Error("getting site subs; returning without subs", "site_id", site.ID, "site_name", site.CompanyName, "error", err)
	}

	return subs
}

func fetchAddOns(ctx context.Context, c *threatdown.Client, site threatdown.Site) []threatdown.AddOn {
	addOns, err := c.ListAddOns(ctx, site.ID)
	if err != nil {
		slog.Error("getting site addons; returning without addons", "site_id", site.ID, "site_name", site.CompanyName, "error", err)
	}

	return addOns
}
