export interface StoreContact {
  label: string
  href: string
}

const DEFAULT_TELEGRAM: StoreContact = {
  label: 'TG',
  href: 'https://t.me/ultra537',
}

const DEFAULT_WHATSAPP: StoreContact = {
  label: '',
  href: '',
}

const isUpstreamPlaceholder = (value: string) => {
  const normalized = value.toLowerCase()
  return normalized.includes('dujiao')
    || normalized.includes('yourchannel')
    || normalized.includes('1234567890')
    || normalized.includes('example')
}

function resolveTelegram(value: unknown): StoreContact {
  const raw = String(value || '').trim()
  if (!raw || isUpstreamPlaceholder(raw)) return DEFAULT_TELEGRAM
  if (raw.startsWith('@')) {
    const username = raw.slice(1)
    return { label: 'TG', href: `https://t.me/${username}` }
  }
  if (/^https?:\/\//i.test(raw)) {
    return { label: 'TG', href: raw }
  }
  return { label: 'TG', href: '' }
}

export function resolveStoreContacts(contactConfig: unknown) {
  const config = contactConfig && typeof contactConfig === 'object'
    ? contactConfig as Record<string, unknown>
    : {}

  return {
    telegram: resolveTelegram(config.telegram),
    whatsapp: DEFAULT_WHATSAPP,
  }
}
