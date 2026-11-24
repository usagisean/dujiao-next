import assert from 'node:assert/strict'
import test from 'node:test'
import { resolveStorefrontLocale } from '../src/i18n/localePreference.ts'

test('new storefront visitors default to Simplified Chinese', () => {
  assert.equal(resolveStorefrontLocale(null), 'zh-CN')
  assert.equal(resolveStorefrontLocale('invalid'), 'zh-CN')
})

test('a saved language preference still takes precedence', () => {
  for (const saved of ['zh-CN', 'zh-TW', 'en-US']) {
    assert.equal(resolveStorefrontLocale(saved), saved)
  }
})
