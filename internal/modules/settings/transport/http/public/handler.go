package publicconfighttp

import (
	"context"
	"strings"
	"time"

	ginutil "github.com/dujiao-next/internal/platform/http/ginutil"

	"github.com/dujiao-next/internal/constants"
	reseller "github.com/dujiao-next/internal/modules/reseller/contract"
	"github.com/dujiao-next/internal/platform/http/response"
	"github.com/dujiao-next/internal/shared/jsonmap"
	"github.com/dujiao-next/internal/version"

	"github.com/gin-gonic/gin"
)

const publicConfigCacheTTL = 60 * time.Second

// ConfigCache 公开配置缓存端口。
type ConfigCache interface {
	CacheKey(resellerID *uint) string
	GetJSON(ctx context.Context, key string, dest *map[string]interface{}) (bool, error)
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// Settings 公开配置所需的设置端口。
type Settings interface {
	GetConfig(defaults map[string]interface{}) (map[string]interface{}, error)
	GetWalletRechargeChannelIDs() []uint
	GetWalletOnlyPayment() bool
	GetAffiliateSettingMap() (map[string]interface{}, error)
	GetSMTPEnabled() bool
	GetRegistrationEnabled(defaultValue bool) (bool, error)
	GetEmailVerificationEnabled(defaultValue bool) (bool, error)
	GetRegistrationInviteOnly() (bool, error)
	GetRegistrationEmailDomainPolicy() (enabled bool, allowedDomains []string, err error)
	GetByKey(key string) (interface{}, error)
	GetActiveHomeAnnouncement() (jsonmap.JSON, bool)
}

// PaymentChannels 公开支付渠道端口。
type PaymentChannels interface {
	GetOrderPaymentChannels() ([]map[string]interface{}, error)
}

// CaptchaPublic 公开验证码配置端口。
type CaptchaPublic interface {
	GetPublicSetting() (jsonmap.JSON, error)
}

// TelegramAuthPublic Telegram 登录公开配置端口。
type TelegramAuthPublic interface {
	PublicConfig() map[string]interface{}
}

// TelegramAuthFallback 无 TelegramAuthService 时的配置回退。
type TelegramAuthFallback struct {
	Enabled     bool
	BotUsername string
	MiniAppURL  string
}

// GoogleAuthPublic Google Identity Services 登录公开配置端口。
type GoogleAuthPublic interface {
	PublicConfig() map[string]interface{}
}

// GoogleAuthFallback 无 GoogleAuthService 时的配置回退。
type GoogleAuthFallback struct {
	Enabled  bool
	ClientID string
}

// ResellerOverlay 分销站配置叠加端口。
type ResellerOverlay interface {
	ApplyPublicConfigOverlay(ctx context.Context, tenant reseller.TenantContext, base map[string]interface{}) (map[string]interface{}, error)
}

// Handler 处理公开站点配置 HTTP 请求。
type Handler struct {
	cache          ConfigCache
	settings       Settings
	payments       PaymentChannels
	captcha        CaptchaPublic
	telegram       TelegramAuthPublic
	fallback       TelegramAuthFallback
	google         GoogleAuthPublic
	googleFallback GoogleAuthFallback
	overlay        ResellerOverlay
}

func NewHandler(
	cache ConfigCache,
	settings Settings,
	payments PaymentChannels,
	captcha CaptchaPublic,
	telegram TelegramAuthPublic,
	fallback TelegramAuthFallback,
	google GoogleAuthPublic,
	googleFallback GoogleAuthFallback,
	overlay ResellerOverlay,
) *Handler {
	if cache == nil {
		panic("public config handler: cache is nil")
	}
	if settings == nil {
		panic("public config handler: settings is nil")
	}
	if payments == nil {
		panic("public config handler: payments is nil")
	}
	return &Handler{
		cache:          cache,
		settings:       settings,
		payments:       payments,
		captcha:        captcha,
		telegram:       telegram,
		fallback:       fallback,
		google:         google,
		googleFallback: googleFallback,
		overlay:        overlay,
	}
}

// GetConfig 获取全局配置。
func (h *Handler) GetConfig(c *gin.Context) {
	defaults := map[string]interface{}{
		"languages":                              append([]string(nil), constants.SupportedLocales...),
		constants.SettingFieldSiteCurrency:       constants.SiteCurrencyDefault,
		constants.SettingFieldStorefrontTemplate: constants.StorefrontTemplateDefault,
		"brand": map[string]interface{}{
			"site_name": "ZX AI Hub",
			"site_url":  "https://shop.zxaihub.com",
			"site_icon": "/zx.svg",
			"site_description": map[string]interface{}{
				"zh-CN": "精选商品与在线服务平台，以清晰的流程提供便捷、可靠的购买体验。",
				"zh-TW": "精選商品與線上服務平台，以清晰的流程提供便捷、可靠的購買體驗。",
				"en-US": "A curated marketplace for products and online services, built around a clear and reliable purchase experience.",
			},
		},
		"contact": map[string]interface{}{
			"telegram": "https://t.me/ultra537",
		},
		"about": map[string]interface{}{
			"hero": map[string]interface{}{
				"title": map[string]interface{}{
					"zh-CN": "关于我们", "zh-TW": "關於我們", "en-US": "About Us",
				},
				"subtitle": map[string]interface{}{
					"zh-CN": "了解 ZX AI Hub", "zh-TW": "了解 ZX AI Hub", "en-US": "Meet ZX AI Hub",
				},
			},
			"introduction": map[string]interface{}{
				"zh-CN": "ZX AI Hub 是一个面向多种商品与在线服务的自助购买平台。我们重视清晰的信息、可靠的交易流程和可持续的服务体验。",
				"zh-TW": "ZX AI Hub 是一個面向多種商品與線上服務的自助購買平台。我們重視清晰的資訊、可靠的交易流程和可持續的服務體驗。",
				"en-US": "ZX AI Hub is a self-service marketplace for a wide range of products and online services. We value clear information, reliable transactions, and a service experience built to last.",
			},
			"services": map[string]interface{}{
				"title": map[string]interface{}{
					"zh-CN": "我们的服务", "zh-TW": "我們的服務", "en-US": "Our Services",
				},
				"items": []interface{}{
					map[string]interface{}{"zh-CN": "精选商品与服务", "zh-TW": "精選商品與服務", "en-US": "Curated products and services"},
					map[string]interface{}{"zh-CN": "清晰透明的购买流程", "zh-TW": "清晰透明的購買流程", "en-US": "A clear and transparent purchase flow"},
					map[string]interface{}{"zh-CN": "可靠的订单与交付支持", "zh-TW": "可靠的訂單與交付支援", "en-US": "Reliable order and delivery support"},
					map[string]interface{}{"zh-CN": "持续更新的产品信息", "zh-TW": "持續更新的產品資訊", "en-US": "Continuously updated product information"},
				},
			},
			"contact": map[string]interface{}{
				"title": map[string]interface{}{
					"zh-CN": "联系我们", "zh-TW": "聯絡我們", "en-US": "Contact Us",
				},
				"text": map[string]interface{}{
					"zh-CN": "如果您对商品或订单有任何疑问，欢迎通过站内提供的方式联系我们。",
					"zh-TW": "如果您對商品或訂單有任何疑問，歡迎透過站內提供的方式聯絡我們。",
					"en-US": "If you have questions about a product or order, contact us through the channels provided on this site.",
				},
			},
		},
		"seo": map[string]interface{}{
			"title": map[string]interface{}{
				"zh-CN": "ZX AI Hub", "zh-TW": "ZX AI Hub", "en-US": "ZX AI Hub",
			},
			"description": map[string]interface{}{
				"zh-CN": "精选商品与在线服务，购买流程清晰，订单状态可查询。",
				"zh-TW": "精選商品與線上服務，購買流程清晰，訂單狀態可查詢。",
				"en-US": "Curated products and online services with clear purchasing and trackable orders.",
			},
		},
		"scripts": make([]interface{}, 0),
	}

	tenant, _ := reseller.TenantFromContext(c.Request.Context())
	cacheKey := h.cache.CacheKey(tenant.ResellerID)

	var cached map[string]interface{}
	if hit, err := h.cache.GetJSON(c.Request.Context(), cacheKey, &cached); err == nil && hit {
		cached["server_time"] = time.Now().UnixMilli()
		cached["app_version"] = version.Version
		response.Success(c, cached)
		return
	}

	data, err := h.settings.GetConfig(defaults)
	if err != nil {
		ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", err)
		return
	}

	publicChannels, err := h.payments.GetOrderPaymentChannels()
	if err != nil {
		ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", err)
		return
	}
	data["payment_channels"] = publicChannels

	walletRechargeChannelIDs := h.settings.GetWalletRechargeChannelIDs()
	if len(walletRechargeChannelIDs) > 0 {
		data["wallet_recharge_channel_ids"] = walletRechargeChannelIDs
	}
	if h.settings.GetWalletOnlyPayment() {
		data["wallet_only_payment"] = true
	}

	if h.captcha != nil {
		publicCaptcha, captchaErr := h.captcha.GetPublicSetting()
		if captchaErr != nil {
			ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", captchaErr)
			return
		}
		data["captcha"] = publicCaptcha
	}
	telegramAuthConfig := map[string]interface{}{
		"enabled":      false,
		"bot_username": "",
		"mini_app_url": "",
	}
	if h.telegram != nil {
		telegramAuthConfig = h.telegram.PublicConfig()
	} else {
		telegramAuthConfig["enabled"] = h.fallback.Enabled
		telegramAuthConfig["bot_username"] = strings.TrimSpace(h.fallback.BotUsername)
		telegramAuthConfig["mini_app_url"] = strings.TrimSpace(h.fallback.MiniAppURL)
	}
	data["telegram_auth"] = telegramAuthConfig

	data["google_auth"] = resolveGoogleAuthPublicConfig(h.google, h.googleFallback)

	affiliateSetting, err := h.settings.GetAffiliateSettingMap()
	if err != nil {
		ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", err)
		return
	}
	data["affiliate"] = affiliateSetting

	data["smtp_enabled"] = h.settings.GetSMTPEnabled()
	registrationEnabled, _ := h.settings.GetRegistrationEnabled(true)
	emailVerificationEnabled, _ := h.settings.GetEmailVerificationEnabled(true)
	data["registration_enabled"] = registrationEnabled
	data["email_verification_enabled"] = emailVerificationEnabled
	inviteOnly, _ := h.settings.GetRegistrationInviteOnly()
	data["invite_only_enabled"] = registrationEnabled && inviteOnly
	enabled, allowedDomains, policyErr := h.settings.GetRegistrationEmailDomainPolicy()
	if policyErr != nil {
		ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", policyErr)
		return
	}
	data["email_domain_allowlist_enabled"] = enabled
	data["allowed_email_domains"] = allowedDomains

	navConfigVal, _ := h.settings.GetByKey(constants.SettingKeyNavConfig)
	if navConfigVal != nil {
		data["nav_config"] = navConfigVal
	} else {
		data["nav_config"] = map[string]interface{}{
			"builtin":      map[string]interface{}{"blog": true, "notice": true, "about": true},
			"custom_items": make([]interface{}, 0),
		}
	}

	if announcement, ok := h.settings.GetActiveHomeAnnouncement(); ok {
		data["announcement"] = announcement
	}

	if h.overlay != nil {
		overlaid, overlayErr := h.overlay.ApplyPublicConfigOverlay(c.Request.Context(), tenant, data)
		if overlayErr != nil {
			ginutil.RespondError(c, response.CodeInternal, "error.config_fetch_failed", overlayErr)
			return
		}
		data = overlaid
	} else if tenant.ResellerID == nil {
		data["tenant"] = map[string]interface{}{"mode": "main", "host": tenant.Host}
	}

	_ = h.cache.SetJSON(c.Request.Context(), cacheKey, data, publicConfigCacheTTL)
	data["server_time"] = time.Now().UnixMilli()
	data["app_version"] = version.Version
	response.Success(c, data)
}

func stringValue(value interface{}) string {
	text, _ := value.(string)
	return text
}

func resolveGoogleAuthPublicConfig(source GoogleAuthPublic, fallback GoogleAuthFallback) map[string]interface{} {
	enabled := fallback.Enabled
	clientID := strings.TrimSpace(fallback.ClientID)
	if source != nil {
		publicGoogle := source.PublicConfig()
		clientID = strings.TrimSpace(stringValue(publicGoogle["client_id"]))
		enabled, _ = publicGoogle["enabled"].(bool)
	}
	return map[string]interface{}{
		"enabled":   enabled && clientID != "",
		"client_id": clientID,
	}
}
