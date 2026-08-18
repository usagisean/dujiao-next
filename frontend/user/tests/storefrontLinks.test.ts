import assert from 'node:assert/strict'
import test from 'node:test'

import { EXTERNAL_BLOG_URL } from '../src/config/links.ts'

test('the storefront blog entry points to the independently hosted blog', () => {
  assert.equal(EXTERNAL_BLOG_URL, 'https://blog.zxaihub.com/')
})
