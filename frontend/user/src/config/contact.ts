export interface StoreContact {
  label: string
  href: string
}

const DEFAULT_TELEGRAM: StoreContact = {
  label: '@ultra537',
  href: 'https://t.me/ultra537',
}

const DEFAULT_WHATSAPP: StoreContact = {
  label: '@ultrazx',
  // WhatsApp usernames are still rolling out and do not yet have a stable
  // public profile URL format. Show the exact handle without inventing a URL.
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
    return { label: raw, href: `https://t.me/${username}` }
  }
  if (/^https?:\/\//i.test(raw)) {
    const username = raw.match(/(?:t\.me|telegram\.me)\/([^/?#]+)/i)?.[1]
    return { label: username ? `@${username}` : 'Telegram', href: raw }
  }
  return { label: raw, href: '' }
}

function resolveWhatsApp(value: unknown): StoreContact {
  const raw = String(value || '').trim()
  if (!raw || isUpstreamPlaceholder(raw)) return DEFAULT_WHATSAPP
  if (raw.startsWith('@')) return { label: raw, href: '' }
  if (/^https?:\/\//i.test(raw)) return { label: 'WhatsApp', href: raw }
  return { label: raw, href: '' }
}

export function resolveStoreContacts(contactConfig: unknown) {
  const config = contactConfig && typeof contactConfig === 'object'
    ? contactConfig as Record<string, unknown>
    : {}

  return {
    telegram: resolveTelegram(config.telegram),
    whatsapp: resolveWhatsApp(config.whatsapp),
  }
}

