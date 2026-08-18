import { createI18n } from 'vue-i18n'
// English is the storefront default and is bundled with the first paint.
// Chinese locales stay available as on-demand language packs.
import enUS from './locales/en-US.json'

const supportedLocales = ['zh-CN', 'zh-TW', 'en-US']

// 非默认语言的懒加载器：切换语言时才下载对应语言包 chunk
const localeLoaders: Record<string, () => Promise<{ default: Record<string, unknown> }>> = {
    'zh-CN': () => import('./locales/zh-CN.json'),
    'zh-TW': () => import('./locales/zh-TW.json'),
}

const loadedLocales = new Set(['en-US'])

export function detectLocale(): string {
    const saved = localStorage.getItem('locale')
    if (saved && supportedLocales.includes(saved)) return saved

    return 'en-US'
}

const i18n = createI18n({
    legacy: false,
    locale: 'en-US',
    fallbackLocale: 'en-US',
    // Keep the locale key widened because Chinese packs are injected at runtime.
    messages: { 'en-US': enUS } as Record<string, typeof enUS>,
})

async function loadLocaleMessages(locale: string): Promise<void> {
    if (loadedLocales.has(locale)) return
    const loader = localeLoaders[locale]
    if (!loader) return
    const messages = await loader()
    i18n.global.setLocaleMessage(locale, (messages.default ?? messages) as any)
    loadedLocales.add(locale)
}

// 记录最新一次切换请求，防止快速连续切换时慢加载的旧请求覆盖新选择
let pendingLocale = ''

// 切换语言：先确保语言包加载完成再更新 locale，避免闪现缺失文案。
// 语言包加载失败时保持当前语言不变，避免整页文案退化为 key。
export async function setI18nLocale(locale: string): Promise<void> {
    if (!supportedLocales.includes(locale)) return
    pendingLocale = locale
    try {
        await loadLocaleMessages(locale)
    } catch (error) {
        console.error(`Failed to load locale messages: ${locale}`, error)
        return
    }
    if (pendingLocale !== locale) return
    i18n.global.locale.value = locale
}

// 空闲时预取其余语言包：让语言切换即时生效，也避免部署更新后旧页面
// 再切语言时请求已失效的旧 chunk（预取时已进入浏览器缓存）
export function warmupLocaleMessages(): void {
    if (typeof window === 'undefined') return

    const prefetch = () => {
        for (const locale of Object.keys(localeLoaders)) {
            void loadLocaleMessages(locale).catch(() => {
                // 预取失败静默忽略，实际切换时会重试并走 setI18nLocale 的错误处理
            })
        }
    }

    if ('requestIdleCallback' in window && typeof window.requestIdleCallback === 'function') {
        window.requestIdleCallback(prefetch, { timeout: 3000 })
        return
    }
    window.setTimeout(prefetch, 1000)
}

export default i18n
