package storefrontmeta

import (
	"net/http/httptest"
	"testing"

	"github.com/dujiao-next/internal/app/container"
	"github.com/dujiao-next/internal/constants"
	settingsapp "github.com/dujiao-next/internal/modules/settings/application"
	"github.com/dujiao-next/internal/shared/jsonmap"
	"github.com/dujiao-next/internal/testkit/memorysettings"
)

func TestPageMetadataUsesSimplifiedChineseByDefault(t *testing.T) {
	settings := memorysettings.New()
	settings.Values[constants.SettingKeySiteConfig] = jsonmap.JSON{
		"brand": map[string]any{"site_url": "https://shop.example", "site_name": "ZX NEXUS"},
		"seo": map[string]any{
			"title":       map[string]any{"zh-CN": "ZX NEXUS 数字商店", "en-US": "ZX NEXUS Store"},
			"description": map[string]any{"zh-CN": "浏览商品", "en-US": "Browse products"},
		},
	}
	resolve := New(&container.Container{SettingService: settingsapp.NewService(settings)})
	for path, want := range map[string]string{"/": "ZX NEXUS 数字商店", "/about": "关于我们 | ZX NEXUS", "/refund": "售后与退款 | ZX NEXUS"} {
		meta := resolve(httptest.NewRequest("GET", "https://shop.example"+path, nil), path, "")
		if meta.Title != want || meta.Description != "浏览商品" {
			t.Errorf("%s title/description = %q / %q", path, meta.Title, meta.Description)
		}
	}
}
