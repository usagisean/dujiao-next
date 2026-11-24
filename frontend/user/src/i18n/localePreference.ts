export const DEFAULT_STOREFRONT_LOCALE = 'zh-CN'
export const storefrontLocales: readonly string[] = [DEFAULT_STOREFRONT_LOCALE, 'zh-TW', 'en-US']

export function resolveStorefrontLocale(saved: string | null): string {
    return saved && storefrontLocales.includes(saved) ? saved : DEFAULT_STOREFRONT_LOCALE
}
