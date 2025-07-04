package hmeili

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hcfg"
	"github.com/meilisearch/meilisearch-go"
	"strings"
)

type Meili struct {
	Code  string
	meili meilisearch.ServiceManager
}

func New(code string) *Meili {
	m := &Meili{Code: code}
	m.doInit()
	return m
}

func (m *Meili) Helix() helix.Helix {
	helixCode := fmt.Sprintf("hmeili_%s", strings.ToLower(m.Code))
	return helix.BuildHelix(helixCode, func() (context.Context, error) {
		return nil, nil
	}, func(ctx context.Context) {

	})
}

func (m *Meili) Meili() meilisearch.ServiceManager {
	return m.meili
}

func (m *Meili) doInit() {
	host := hcfg.GetString(m.cfg("host"), "http://127.0.0.1:7700")
	var options []meilisearch.Option
	options = append(options, meilisearch.WithAPIKey(hcfg.GetString(m.cfg("api.key"))))
	m.meili = meilisearch.New(host, options...)
}

func (m *Meili) cfg(prefix string) string {
	return fmt.Sprintf("hmeili.%s.%s", strings.ToLower(m.Code), prefix)
}
