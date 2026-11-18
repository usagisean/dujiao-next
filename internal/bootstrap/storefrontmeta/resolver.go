package storefrontmeta

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/dujiao-next/internal/app/container"
	"github.com/dujiao-next/internal/constants"
	productcontract "github.com/dujiao-next/internal/modules/catalog/product/contract"
	contentcontract "github.com/dujiao-next/internal/modules/content/contract"
	"github.com/dujiao-next/internal/shared/jsonmap"

	"github.com/dujiao-next/internal/web"
)

func New(c *container.Container) web.UserPageResolver {
	return func(r *http.Request, route, slug string) web.PageMetadata {
		switch route {
		case "/", "/products", "/products/:id", "/categories/:id", "/blog", "/blog/:id",
			"/notice", "/about", "/terms", "/privacy", "/refund":
		default:
			return web.PageMetadata{NoIndex: true}
		}

		cfg, err := c.SettingService.GetByKey(constants.SettingKeySiteConfig)
		if err != nil {
			return web.PageMetadata{}
		}
		brand, _ := cfg["brand"].(map[string]any)
		seo, _ := cfg["seo"].(map[string]any)
		origin, err := url.Parse(settingText(brand["site_url"]))
		// Never expose main-site metadata or unpublished tenant products on another host.
		if err != nil || (origin.Scheme != "https" && origin.Scheme != "http") ||
			!strings.EqualFold(origin.Host, r.Host) || origin.User != nil {
			return web.PageMetadata{}
		}
		base := origin.Scheme + "://" + origin.Host
		name := firstText(settingText(brand["site_name"]), "ZX NEXUS")
		description := firstText(englishText(seo["description"]), englishText(brand["site_description"]),
			"Browse digital products, track orders, and get help from ZX NEXUS.")
		meta := web.PageMetadata{SiteName: name, Description: description,
			Image: base + "/zx-logo.png", Home: route == "/"}
		title := map[string]string{
			"/products": "Products", "/blog": "Blog", "/notice": "Announcements",
			"/about": "About", "/terms": "Terms of Service", "/privacy": "Privacy Policy",
			"/refund": "Support & Refunds",
		}[route]
		canonicalPath := route
		switch route {
		case "/products/:id":
			product, err := c.ProductReadService.GetPublicBySlug(slug)
			if err != nil {
				return web.PageMetadata{Missing: errors.Is(err, productcontract.ErrNotFound)}
			}
			canonicalPath = "/products/" + product.Slug
			title = firstText(englishText(product.TitleJSON), "Product")
			if len(product.Images) > 0 {
				meta.Image = shareImage(base, product.Images[0], meta.Image)
			}
			meta.Description = firstText(englishText(product.DescriptionJSON), description)
		case "/categories/:id":
			category, err := c.CategoryRepo.GetBySlug(slug)
			if err != nil {
				return web.PageMetadata{}
			}
			if category == nil || !category.IsActive {
				return web.PageMetadata{Missing: true}
			}
			canonicalPath = "/categories/" + category.Slug
			title = firstText(englishText(category.NameJSON), "Products")
		case "/blog/:id":
			post, err := c.ContentPostService.GetPublicBySlug(r.Context(), slug)
			if err != nil {
				return web.PageMetadata{Missing: errors.Is(err, contentcontract.ErrNotFound)}
			}
			canonicalPath = "/blog/" + post.Slug
			title = firstText(englishText(post.TitleJSON), "Updates")
			meta.Description = firstText(englishText(post.SummaryJSON), description)
			meta.Image = shareImage(base, post.Thumbnail, meta.Image)
		}
		if meta.Home {
			meta.Title = firstText(englishText(seo["title"]), name)
		} else {
			meta.Title = title + " | " + name
		}
		meta.Canonical = web.StorefrontURL(base, canonicalPath)
		return meta
	}
}

func settingText(value any) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}

// The first response uses the same default as the storefront's first visit.
func englishText(value any) string {
	switch fields := value.(type) {
	case map[string]any:
		return settingText(fields["en-US"])
	case jsonmap.JSON:
		return settingText(fields["en-US"])
	default:
		return settingText(value)
	}
}

func firstText(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func shareImage(origin, value, fallback string) string {
	u, err := url.Parse(strings.TrimSpace(value))
	if err != nil || value == "" || u.User != nil {
		return fallback
	}
	if (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" {
		return u.String()
	}
	if u.Scheme == "" && u.Host == "" && strings.HasPrefix(u.Path, "/") {
		return origin + u.EscapedPath()
	}
	return fallback
}
