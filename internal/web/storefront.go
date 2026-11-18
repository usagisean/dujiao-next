package web

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// UserPageResolver supplies metadata for recognized storefront routes.
type UserPageResolver func(r *http.Request, route, slug string) PageMetadata

type PageMetadata struct {
	Title, Description, Canonical, Image, SiteName string
	Home, Missing, NoIndex                         bool
}

// userRoute matches the history-mode route shapes in frontend/user/src/router/index.ts.
func userRoute(raw string) (route, slug string) {
	p := strings.ToLower(strings.TrimSuffix(raw, "/"))
	switch p {
	case "", "/index.html":
		return "/", ""
	case "/products", "/about", "/terms", "/privacy", "/refund", "/blog", "/notice",
		"/cart", "/checkout", "/pay", "/me", "/me/profile", "/me/security", "/me/orders", "/me/wallet",
		"/me/gift-cards", "/me/api", "/me/affiliate", "/me/reseller", "/guest/orders",
		"/reseller", "/reseller/apply", "/reseller/domains", "/reseller/site", "/reseller/products",
		"/reseller/orders", "/reseller/finance", "/reseller/ledger", "/reseller/withdraws",
		"/auth/login", "/auth/register", "/auth/forgot", "/auth/google/callback", "/auth/telegram/callback":
		return p, ""
	}
	dir, value, ok := strings.Cut(strings.TrimSuffix(raw, "/")[1:], "/")
	if !ok || value == "" {
		return "", ""
	}
	dir = "/" + strings.ToLower(dir)
	switch dir {
	case "/products", "/categories", "/blog", "/orders", "/recharge-orders":
		if !strings.Contains(value, "/") {
			return dir + "/:id", value
		}
	case "/guest", "/reseller":
		parent, id, nested := strings.Cut(value, "/")
		if nested && strings.EqualFold(parent, "orders") && id != "" && !strings.Contains(id, "/") {
			return dir + "/orders/:id", id
		}
	}
	return "", ""
}

// StorefrontURL builds a canonical URL without reflecting tracking or credential query strings.
func StorefrontURL(origin, route string) string {
	return origin + (&url.URL{Path: route}).EscapedPath()
}

var pageHead = template.Must(template.New("page-head").Parse(`<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}" data-zx-page-head>
<link rel="canonical" href="{{.Canonical}}" data-zx-page-head>
<meta property="og:type" content="website" data-zx-page-head>
<meta property="og:site_name" content="{{.SiteName}}" data-zx-page-head>
<meta property="og:title" content="{{.Title}}" data-zx-page-head>
<meta property="og:description" content="{{.Description}}" data-zx-page-head>
<meta property="og:url" content="{{.Canonical}}" data-zx-page-head>
<meta property="og:image" content="{{.Image}}" data-zx-page-head>
<meta name="twitter:card" content="summary_large_image" data-zx-page-head>
<meta name="twitter:title" content="{{.Title}}" data-zx-page-head>
<meta name="twitter:description" content="{{.Description}}" data-zx-page-head>
<meta name="twitter:image" content="{{.Image}}" data-zx-page-head>
{{if .Home}}<script type="application/ld+json" data-zx-page-head>
{"@context":"https://schema.org","@graph":[{"@type":"Organization","@id":{{printf "%s#organization" .Canonical}},"name":{{.SiteName}},"url":{{.Canonical}},"logo":{{.Image}}},{"@type":"WebSite","name":{{.SiteName}},"url":{{.Canonical}},"publisher":{"@id":{{printf "%s#organization" .Canonical}}}}]}
</script>{{end}}`))

func renderPageHead(index []byte, meta PageMetadata) []byte {
	var output bytes.Buffer
	if err := pageHead.Execute(&output, meta); err != nil {
		return index
	}
	return bytes.Replace(index, []byte("<title>ZX NEXUS</title>"), output.Bytes(), 1)
}
