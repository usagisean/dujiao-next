export const DEFAULT_BRAND_NAME = 'ZX AI Hub'
export const DEFAULT_BRAND_MARK = 'ZX'

/**
 * Keep the upstream default configurable while giving this fixed fork its own
 * identity before the site settings are populated.
 */
export function resolveBrandName(configuredName: unknown): string {
  const name = typeof configuredName === 'string' ? configuredName.trim() : ''
  if (!name || name === 'Dujiao-Next') return DEFAULT_BRAND_NAME
  return name
}

