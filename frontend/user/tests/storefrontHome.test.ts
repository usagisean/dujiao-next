import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const readSource = (path: string) => readFileSync(new URL(path, import.meta.url), 'utf8')

const assertBefore = (source: string, first: string, second: string) => {
  const before = source.indexOf(first)
  const after = source.indexOf(second)
  assert.ok(before >= 0 && after > before, `Expected ${first} before ${second}`)
}

test('Classic home renders the catalog before banners in both modes', () => {
  const source = readSource('../src/views/Home.vue')
  assertBefore(source, '<ProductListItem', 'zx-list-hero-overlay')
  assertBefore(source, '<ProductCard', 'zx-home-hero')
  assertBefore(source, 'v-if="productsLoading"', 'v-else-if="products.length > 0"')
})

test('Vault home renders the catalog before categories and banners', () => {
  const source = readSource('../src/templates/vault/Home.vue')
  assertBefore(source, '<VaultProductCard', 'v-if="topCategories.length"')
  assertBefore(source, '<VaultProductListItem', '<VaultBannerHero />')
  assertBefore(source, '<VaultProductCard', '<VaultBannerHero />')
})
