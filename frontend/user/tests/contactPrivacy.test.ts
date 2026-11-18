import assert from 'node:assert/strict'
import test from 'node:test'

import { resolveStoreContacts } from '../src/config/contact.ts'

test('store contact labels never expose usernames and WhatsApp stays hidden', () => {
  const contacts = resolveStoreContacts({
    telegram: '@fixture_support',
    whatsapp: '@example_whatsapp',
  })

  assert.deepEqual(contacts.telegram, {
    label: 'TG',
    href: 'https://t.me/fixture_support',
  })
  assert.deepEqual(contacts.whatsapp, {
    label: '',
    href: '',
  })
})

test('Telegram URLs keep a generic public label', () => {
  const contacts = resolveStoreContacts({
    telegram: 'https://t.me/fixture_support',
  })

  assert.equal(contacts.telegram.label, 'TG')
  assert.equal(contacts.telegram.href, 'https://t.me/fixture_support')
})

test('Telegram fallback is the same support link as the production configuration', () => {
  assert.deepEqual(resolveStoreContacts(undefined).telegram, {
    label: 'TG',
    href: 'https://t.me/ultrazzx',
  })
})
