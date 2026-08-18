import assert from 'node:assert/strict'
import test from 'node:test'

import { resolveStoreContacts } from '../src/config/contact.ts'

test('store contact labels never expose usernames and WhatsApp stays hidden', () => {
  const contacts = resolveStoreContacts({
    telegram: '@ultra537',
    whatsapp: '@ultrazx',
  })

  assert.deepEqual(contacts.telegram, {
    label: 'TG',
    href: 'https://t.me/ultra537',
  })
  assert.deepEqual(contacts.whatsapp, {
    label: '',
    href: '',
  })
})

test('Telegram URLs keep a generic public label', () => {
  const contacts = resolveStoreContacts({
    telegram: 'https://t.me/ultra537',
  })

  assert.equal(contacts.telegram.label, 'TG')
  assert.equal(contacts.telegram.href, 'https://t.me/ultra537')
})
