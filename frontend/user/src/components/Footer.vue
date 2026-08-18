<template>
  <footer class="zx-footer relative text-muted-foreground border-t overflow-hidden">
    <div class="container mx-auto px-4 py-12 md:py-16 relative">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8 md:gap-12 mb-12 md:mb-16">
        <!-- Brand -->
        <div class="col-span-2 space-y-6">
          <div class="flex items-center space-x-3">
            <img
              v-if="brandLogo"
              :src="brandLogo"
              :alt="brandSiteName"
              class="h-10 w-10 rounded-full object-cover ring-1 ring-border"
            />
            <div v-else class="zx-brand-mark" aria-hidden="true">
              <span>ZX</span>
            </div>
            <h3 class="text-foreground text-xl font-bold tracking-tight">{{ brandSiteName }}</h3>
          </div>
          <p class="text-sm leading-relaxed max-w-sm text-muted-foreground">
            {{ brandDescription || t('footer.description') }}
          </p>
          <div class="flex space-x-4">
            <!-- Social Icons (Placeholder) -->
            <!--
            <a href="#" class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"/></svg>
            </a>
            <a href="#" class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
            </a>
            -->
          </div>
        </div>

        <!-- Links -->
        <div>
          <h4 class="text-foreground font-bold mb-6 tracking-wide">{{ t('footer.quickLinks') }}</h4>
          <ul class="space-y-3 text-sm">
            <li v-for="item in quickLinks" :key="item.path">
              <component
                :is="item.external ? 'a' : 'router-link'"
                :to="item.external ? undefined : item.path"
                :href="item.external ? item.path : undefined"
                :target="item.external ? '_blank' : undefined"
                :rel="item.external ? 'noopener noreferrer' : undefined"
                class="text-muted-foreground hover:text-foreground transition-colors flex items-center gap-2 group"
              >
                <component :is="item.icon" class="w-4 h-4 shrink-0 opacity-50 group-hover:opacity-80 transition-opacity" />
                {{ t(item.label) }}
              </component>
            </li>
          </ul>
        </div>

        <!-- Contact -->
        <div>
          <h4 class="text-foreground font-bold mb-6 tracking-wide">{{ t('footer.contact') }}</h4>
          <div class="space-y-4">
            <component :is="contacts.telegram.href ? 'a' : 'div'" :href="contacts.telegram.href || undefined"
              :target="contacts.telegram.href ? '_blank' : undefined" rel="noopener noreferrer"
              class="flex items-center space-x-3 text-sm hover:text-gray-900 dark:hover:text-white transition-colors p-3 rounded-lg bg-gray-50 dark:bg-white/5 hover:bg-gray-100 dark:hover:bg-white/10 border border-gray-100 dark:border-white/5 hover:border-gray-200 dark:hover:border-white/10">
              <svg class="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.295-.6.295-.002 0-.003 0-.005 0l.213-3.054 5.56-5.022c.24-.213-.054-.334-.373-.121l-6.869 4.326-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.538-.196 1.006.128.832.941z" />
              </svg>
              <span>{{ contacts.telegram.label }}</span>
            </component>
          </div>
        </div>
      </div>

      <!-- Copyright -->
      <div
        class="border-t pt-8 flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
        <div class="space-y-1 text-center md:text-left">
          <p>&copy; {{ currentYear }} {{ brandSiteName }}. {{ t('footer.rights') }}</p>
        </div>
        <div class="flex flex-col items-center gap-2 md:items-end">
          <div class="flex flex-wrap items-center gap-x-4 gap-y-1 justify-center md:justify-end">
            <router-link to="/privacy" class="hover:text-gray-900 dark:hover:text-gray-400">{{ t('footer.privacy') || 'Privacy Policy' }}</router-link>
            <router-link to="/terms" class="hover:text-gray-900 dark:hover:text-gray-400">{{ t('footer.terms') || 'Terms of Service' }}</router-link>
          </div>
          <div v-if="footerLinks.length" class="flex flex-wrap items-center gap-x-4 gap-y-1 justify-center md:justify-end">
            <a
              v-for="link in footerLinks"
              :key="link.name"
              :href="link.url || 'javascript:void(0)'"
              :target="link.url ? '_blank' : undefined"
              rel="noopener noreferrer"
              class="hover:text-gray-900 dark:hover:text-gray-400"
            >{{ link.name }}</a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Home, LayoutGrid, Newspaper, Info } from 'lucide-vue-next'
import { useAppStore } from '../stores/app'
import { getImageUrl } from '../utils/image'
import { getLocalizedText } from '../utils/resellerSiteConfig'
import { resolveBrandName } from '../config/brand'
import { resolveStoreContacts } from '../config/contact'
import { EXTERNAL_BLOG_URL } from '../config/links'

const { t } = useI18n()
const appStore = useAppStore()

const config = computed(() => appStore.config)
const contacts = computed(() => resolveStoreContacts(config.value?.contact))

const brandSiteName = computed(() => {
  return resolveBrandName(config.value?.brand?.site_name)
})

const brandDescription = computed(() => {
  const desc = config.value?.brand?.site_description
  if (desc && typeof desc === 'object') {
    const val = desc[appStore.locale] || desc['zh-CN'] || ''
    return typeof val === 'string' ? val.trim() : ''
  }
  return ''
})

const brandLogo = computed(() => {
  const raw = String(config.value?.brand?.site_logo || '').trim()
  return raw && raw !== '/dj.svg' ? getImageUrl(raw) : '/zx-logo.png'
})

const isListMode = computed(() => config.value?.template_mode === 'list')
const navConfig = computed(() => config.value?.nav_config as { builtin?: Record<string, boolean> } | undefined)

interface QuickLink {
  path: string
  label: string
  icon: typeof Home
  external?: boolean
}

const quickLinks = computed<QuickLink[]>(() => {
  const items: QuickLink[] = [
    { path: '/', label: 'nav.home', icon: Home },
  ]
  if (!isListMode.value) {
    items.push({ path: '/products', label: 'nav.products', icon: LayoutGrid })
  }
  const builtin = navConfig.value?.builtin
  if (!builtin || builtin.blog !== false) {
    items.push({ path: EXTERNAL_BLOG_URL, label: 'nav.blog', icon: Newspaper, external: true })
  }
  if (!builtin || builtin.about !== false) {
    items.push({ path: '/about', label: 'nav.about', icon: Info })
  }
  return items
})

const footerLinks = computed(() => {
  const links = config.value?.footer_links
  if (!Array.isArray(links)) return []
  return links
    .map((item: any) => ({
      name: typeof item?.name === 'string' ? item.name.trim() : getLocalizedText(item?.name, appStore.locale),
      url: String(item?.url || '').trim(),
    }))
    .filter((item) => item.name)
})

const currentYear = new Date().getFullYear()
</script>
