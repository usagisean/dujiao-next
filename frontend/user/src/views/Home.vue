<template>
  <div class="home-page min-h-screen bg-background text-foreground">

    <!-- ==================== LIST MODE ==================== -->
    <template v-if="templateMode === 'list'">
      <!-- Main: Left Categories + Right Product List -->
      <section class="relative z-10 pb-4 pt-20 md:pt-24">
        <div class="container mx-auto px-4">
          <h1 class="sr-only">{{ t('nav.products') }}</h1>
          <div class="flex flex-col lg:flex-row gap-6">

            <CategorySidebar
              :categories="listCategoryGroups"
              :selected-category="listSelectedCategory"
              :expanded-parent-ids="listExpandedParentIds"
              :show-drawer="listShowFilterDrawer"
              compact
              @select-category="listSelectCategory"
              @toggle-parent="listToggleParentCategory"
              @update:show-drawer="listShowFilterDrawer = $event"
            />

            <!-- Right: Product List -->
            <main class="flex-1 min-w-0">
              <!-- Search Bar -->
              <div class="relative mb-4">
                <div class="absolute inset-y-0 left-3.5 flex items-center pointer-events-none">
                  <Search class="w-4 h-4 text-muted-foreground" />
                </div>
                <input
                  v-model="listSearchQuery"
                  type="text"
                  class="w-full h-10 pl-10 pr-10 rounded-xl border bg-card text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 text-foreground placeholder:text-muted-foreground transition-shadow"
                  :placeholder="t('products.searchBoxPlaceholder')"
                  :aria-label="t('products.searchLabel')"
                  @keydown.enter="listOnSearch"
                />
                <button
                  v-if="listSearchQuery"
                  type="button"
                  class="absolute inset-y-0 right-3 flex items-center text-muted-foreground hover:text-foreground transition-colors"
                  @click="listClearSearch"
                >
                  <X class="w-4 h-4" />
                </button>
              </div>

              <!-- Loading Skeleton -->
              <div v-if="listLoading" class="space-y-6">
                <div v-for="i in 3" :key="i">
                  <div class="flex items-center gap-2 mb-3 px-0.5">
                    <div class="h-5 w-5 rounded theme-skeleton"></div>
                    <div class="h-4 w-28 rounded theme-skeleton"></div>
                  </div>
                  <div class="space-y-2">
                    <div v-for="j in 3" :key="j"
                      class="bg-card rounded-xl border flex items-center h-[72px]">
                      <div class="w-14 h-14 m-2 rounded-lg theme-skeleton flex-shrink-0"></div>
                      <div class="flex-1 px-3 py-2 space-y-2">
                        <div class="h-3.5 w-1/3 rounded theme-skeleton"></div>
                        <div class="h-3 w-1/4 rounded theme-skeleton"></div>
                      </div>
                      <div class="px-4 py-2">
                        <div class="h-4 w-14 rounded theme-skeleton"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Grouped Product List -->
              <div v-else-if="listProductGroups.length > 0" class="space-y-6">
                <div v-for="group in listProductGroups" :key="group.categoryId ?? 'uncategorized'">
                  <!-- Category Header -->
                  <div class="flex items-center gap-2 mb-3 px-0.5">
                    <span class="w-1 h-5 rounded-full bg-primary flex-shrink-0"></span>
                    <img v-if="group.categoryIcon" :src="getImageUrl(group.categoryIcon)"
                      :alt="group.categoryName" loading="lazy" class="h-5 w-5 rounded object-cover flex-shrink-0" />
                    <span class="text-sm font-semibold text-foreground truncate">{{ group.categoryName }}</span>
                    <span class="text-xs text-muted-foreground">({{ group.products.length }})</span>
                  </div>
                  <!-- Products in this category -->
                  <div class="space-y-2">
                    <ProductListItem
                      v-for="(product, idx) in group.products"
                      :key="product.id"
                      :product="product"
                      :index="idx"
                      :animation-step="20"
                      @click="goToProduct"
                      @quick-buy="openQuickBuy"
                    />
                  </div>
                </div>

                <PaginationNav
                  :current-page="listCurrentPage"
                  :total-pages="listTotalPages"
                  :loading="listLoading"
                  compact
                  @change-page="listChangePage"
                />
              </div>

              <!-- Empty State -->
              <EmptyState v-else variant="soft" icon="package" :title="t('products.empty')" />
            </main>
          </div>
        </div>
      </section>

      <!-- Brand banners follow the catalog -->
      <section v-if="showHeroSection" class="relative z-10 py-6 md:py-8">
        <div class="container mx-auto px-4">
          <div class="relative overflow-hidden rounded-2xl border bg-card"
            @touchstart="onBannerTouchStart"
            @touchend="onBannerTouchEnd">
            <Transition name="banner-fade" mode="out-in">
              <img v-if="!bannerLoading && heroImage" :src="heroImage" :key="heroImage" class="absolute inset-0 h-full w-full object-cover" alt="" />
            </Transition>
            <div class="zx-list-hero-overlay absolute inset-0"></div>
            <div v-if="bannerLoading" class="relative flex min-h-[175px] flex-col justify-between p-5 md:min-h-[215px] md:p-7">
              <div class="space-y-3">
                <div class="h-5 w-24 theme-skeleton rounded-full" style="background: rgba(255,255,255,0.35)"></div>
                <div class="h-8 max-w-3xl theme-skeleton rounded-xl md:h-10" style="background: rgba(255,255,255,0.35)"></div>
                <div class="h-4 max-w-2xl theme-skeleton rounded-lg" style="background: rgba(255,255,255,0.3)"></div>
              </div>
            </div>
            <div v-else class="relative flex min-h-[175px] flex-col justify-between p-5 md:min-h-[215px] md:p-7">
              <div v-if="bannerCount > 1" class="mb-3 flex items-center justify-end gap-2">
                <button type="button"
                  class="zx-hero-nav inline-flex h-8 w-8 items-center justify-center rounded-full border transition md:h-9 md:w-9"
                  @click="handlePrevHeroBanner" :aria-label="t('common.previousBanner')">
                  <ChevronLeft class="h-3.5 w-3.5" />
                </button>
                <button type="button"
                  class="zx-hero-nav inline-flex h-8 w-8 items-center justify-center rounded-full border transition md:h-9 md:w-9"
                  @click="handleNextHeroBanner" :aria-label="t('common.nextBanner')">
                  <ChevronRight class="h-3.5 w-3.5" />
                </button>
              </div>
              <div class="space-y-2 sm:space-y-3">
                <span :class="heroBadgeClass">
                  <span class="h-2 w-2 rounded-full bg-emerald-300"></span>
                  {{ heroBadge }}
                </span>
                <h2 class="zx-list-hero-title max-w-4xl text-xl font-semibold tracking-[-0.02em] sm:text-2xl md:text-3xl">
                  {{ heroTitle }}
                </h2>
                <p class="zx-list-hero-subtitle max-w-3xl text-xs leading-relaxed sm:text-sm">
                  {{ heroSubtitle }}
                </p>
              </div>
              <div v-if="bannerCount > 1" class="mt-4 flex items-center gap-2">
                <button v-for="(_, bIdx) in banners" :key="`list-dot-${bIdx}`" type="button"
                  class="h-2 rounded-full transition-all"
                  :class="bIdx === currentBannerIndex ? 'zx-hero-dot-active w-6' : 'zx-hero-dot w-2'"
                  :aria-label="t('common.switchBanner', { n: bIdx + 1 })"
                  @click="selectHeroBanner(bIdx)"></button>
              </div>
            </div>
          </div>
        </div>
      </section>

    </template>

    <!-- ==================== CARD MODE (default) ==================== -->
    <template v-else>
    <section v-if="productsLoading || products.length" id="featured" :aria-busy="productsLoading" class="relative z-10 pb-4 pt-20 md:pb-6 md:pt-28">
      <div class="container mx-auto px-4">
        <div class="mb-4 flex flex-wrap items-center justify-between gap-x-4 gap-y-2 md:mb-6">
          <h1 class="theme-section-heading text-xl sm:text-2xl md:text-3xl">{{ t('home.featured.title') }}</h1>
          <router-link
            to="/products"
            class="inline-flex min-h-11 items-center gap-1.5 text-xs font-semibold text-primary transition-colors hover:text-foreground sm:text-sm"
          >
            {{ t('home.featured.viewAll') }}
            <ArrowRight class="h-4 w-4" />
          </router-link>
        </div>

        <div v-if="productsLoading" aria-hidden="true" class="grid grid-cols-1 gap-3 min-[375px]:grid-cols-2 md:grid-cols-3 md:gap-4 lg:grid-cols-4 xl:grid-cols-5">
          <div v-for="index in 2" :key="index" class="overflow-hidden rounded-2xl border bg-card">
            <div class="theme-skeleton aspect-[4/3]"></div>
            <div class="space-y-2 p-3 md:p-4">
              <div class="theme-skeleton h-4 w-2/3 rounded"></div>
              <div class="theme-skeleton h-4 w-1/3 rounded"></div>
            </div>
          </div>
        </div>
        <div v-else-if="products.length > 0" class="grid grid-cols-1 gap-3 min-[375px]:grid-cols-2 md:grid-cols-3 md:gap-4 lg:grid-cols-4 xl:grid-cols-5">
          <ProductCard
            v-for="(product, idx) in products"
            :key="product.id"
            :product="product"
            :index="idx"
            :animation-step="60"
            @click="goToProduct"
            @quick-buy="openQuickBuy"
          />
        </div>
      </div>
    </section>

    <section v-if="!bannerLoading && showHeroSection" class="zx-home-hero relative z-10 py-5 md:py-8" :class="{ '!pt-20 md:!pt-28': !productsLoading && !products.length }">
      <div class="container mx-auto px-4">
        <div class="zx-hero-shell relative overflow-hidden rounded-2xl border"
          @touchstart="onBannerTouchStart"
          @touchend="onBannerTouchEnd">
          <Transition name="banner-fade" mode="out-in">
            <img v-if="heroImage" :src="heroImage" :key="heroImage" class="absolute inset-0 h-full w-full object-cover" alt="" />
          </Transition>
          <div v-if="heroImage" class="zx-hero-overlay absolute inset-0"></div>
          <div v-else class="zx-hero-mesh absolute inset-0" aria-hidden="true"></div>
          <div class="zx-hero-grid absolute inset-0" aria-hidden="true"></div>

          <div class="relative flex flex-col p-4 sm:p-6 md:p-8">
            <div class="flex min-h-9 items-center justify-between gap-2">
              <span :class="heroBadgeClass">
                <span class="h-1.5 w-1.5 shrink-0 rounded-full bg-orange-300"></span>
                {{ heroBadge }}
              </span>

              <div v-if="bannerCount > 1" class="z-20 flex shrink-0 items-center gap-2">
                <button
                  type="button"
                  class="zx-hero-nav inline-flex h-9 w-9 items-center justify-center rounded-full border backdrop-blur transition md:h-10 md:w-10"
                  @click="handlePrevHeroBanner"
                  :aria-label="t('common.previousBanner')"
                >
                  <ChevronLeft class="h-4 w-4" />
                </button>
                <button
                  type="button"
                  class="zx-hero-nav inline-flex h-9 w-9 items-center justify-center rounded-full border backdrop-blur transition md:h-10 md:w-10"
                  @click="handleNextHeroBanner"
                  :aria-label="t('common.nextBanner')"
                >
                  <ChevronRight class="h-4 w-4" />
                </button>
              </div>
            </div>

            <div class="relative z-10 max-w-3xl py-3 sm:py-4">
              <h2 class="zx-hero-title bg-clip-text text-xl font-medium leading-tight tracking-tight text-transparent sm:text-2xl md:text-3xl">
                {{ heroTitle }}
              </h2>
              <p class="zx-hero-subtitle mt-2 max-w-2xl text-sm leading-relaxed">
                {{ heroSubtitle }}
              </p>

              <div class="mt-4 flex flex-wrap items-center gap-2">
                <button
                  type="button"
                  @click="goToHeroLink"
                  class="zx-hero-primary inline-flex min-h-11 items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold transition"
                >
                  {{ heroPrimaryButtonText }}
                  <ArrowRight class="h-4 w-4" />
                </button>
                <a
                  v-if="telegram.href"
                  :href="telegram.href"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="zx-hero-secondary inline-flex min-h-11 items-center rounded-xl border px-4 py-2 text-sm font-semibold backdrop-blur transition"
                >
                  {{ telegram.label }}
                </a>
              </div>
            </div>

            <div>
              <div v-if="bannerCount > 1" class="flex items-center gap-2">
                <button
                  v-for="(_, index) in banners"
                  :key="`hero-dot-${index}`"
                  type="button"
                  class="h-1 rounded-full transition-all duration-300"
                  :class="index === currentBannerIndex ? 'zx-hero-dot-active w-10' : 'zx-hero-dot w-4'"
                  @click="selectHeroBanner(index)"
                  :aria-label="t('common.switchBanner', { n: index + 1 })"
                ></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <template v-if="latestContentVisible">
    <hr class="theme-section-divider mx-4 md:mx-auto md:max-w-6xl" />

    <section class="relative z-10 py-12">
      <div class="container mx-auto px-4">
        <div class="mb-6 flex items-end justify-between gap-4">
          <div>
            <h2 class="theme-section-heading text-[1.7rem]">{{ t('home.latest.title') }}</h2>
            <p class="mt-1 text-sm text-muted-foreground">{{ t('home.latest.description') }}</p>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <a v-if="blogEnabled" :href="EXTERNAL_BLOG_URL" target="_blank" rel="noopener noreferrer" class="text-muted-foreground transition-colors hover:text-foreground">{{ t('nav.blog') }}</a>
            <router-link v-if="noticeEnabled" to="/notice" class="text-muted-foreground transition-colors hover:text-foreground">{{ t('nav.notice') }}</router-link>
          </div>
        </div>

        <div v-if="posts.length > 0" class="grid grid-cols-1 gap-5 md:grid-cols-3">
          <article
            v-for="post in posts"
            :key="post.id"
            class="cursor-pointer rounded-xl border bg-card p-5 shadow-sm transition hover:shadow-md"
            @click="goToPost(post.slug)"
          >
            <div class="mb-2 text-xs text-muted-foreground">{{ formatDate(post.published_at) }}</div>
            <h3 class="line-clamp-2 text-base font-semibold">{{ getLocalizedText(post.title) }}</h3>
            <p class="mt-2 line-clamp-2 text-sm text-muted-foreground">{{ getLocalizedText(post.summary) }}</p>
            <div class="mt-4 text-sm font-medium text-primary">{{ t('blog.readMore') }}</div>
          </article>
        </div>
        <div v-else class="rounded-2xl border border-dashed py-12 text-center text-muted-foreground">
          {{ t('blog.empty') }}
        </div>
      </div>
    </section>
    </template>
    </template>

    <ProductQuickBuy
      v-if="quickBuyProduct"
      :product="quickBuyProduct"
      :visible="quickBuyVisible"
      @update:visible="quickBuyVisible = $event"
    />

    <AnnouncementModal
      v-if="activeAnnouncement"
      :announcement="activeAnnouncement"
      :visible="announcementVisible"
      @update:visible="announcementVisible = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  Search,
  X,
} from 'lucide-vue-next'
import { postAPI, productAPI } from '../api'
import { getImageUrl } from '../utils/image'
import { useLocalized } from '../composables/useProduct'
import { useBannerCarousel } from '../composables/useBannerCarousel'
import { useProductList } from '../composables/useProductList'
import { useProductListGroups } from '../composables/useProductListGroups'
import { usePageSeo } from '../composables/usePageSeo'
import { useAppStore } from '../stores/app'
import { resolveStoreContacts } from '../config/contact'
import { EXTERNAL_BLOG_URL } from '../config/links'
import ProductCard from '../components/ProductCard.vue'
import ProductListItem from '../components/ProductListItem.vue'
import ProductQuickBuy from '../components/ProductQuickBuy.vue'
import CategorySidebar from '../components/CategorySidebar.vue'
import PaginationNav from '../components/PaginationNav.vue'
import EmptyState from '../components/EmptyState.vue'
import AnnouncementModal from '../components/AnnouncementModal.vue'
import { useAnnouncement, type HomeAnnouncement } from '../composables/useAnnouncement'

const router = useRouter()
const { t } = useI18n()
const { getLocalizedText } = useLocalized()
const appStore = useAppStore()
const telegram = computed(() => resolveStoreContacts(appStore.config?.contact).telegram)

const heroBadgeClass = 'zx-hero-badge inline-flex items-center gap-2 rounded-md border px-2.5 py-0.5 text-xs font-semibold uppercase tracking-wider backdrop-blur-sm'

const templateMode = computed(() => appStore.config?.template_mode || 'card')
const navBuiltin = computed(() => (appStore.config?.nav_config as { builtin?: Record<string, boolean> } | undefined)?.builtin)
const blogEnabled = computed(() => navBuiltin.value?.blog !== false)
const noticeEnabled = computed(() => navBuiltin.value?.notice !== false)
const latestSectionVisible = computed(() => blogEnabled.value || noticeEnabled.value)

// ==================== Shared State ====================
const products = ref<any[]>([])
const productsLoading = ref(true)
const posts = ref<any[]>([])
const latestContentVisible = computed(() => latestSectionVisible.value && posts.value.length > 0)
const quickBuyProduct = ref<any>(null)
const quickBuyVisible = ref(false)

const { shouldShow } = useAnnouncement()
const activeAnnouncement = ref<HomeAnnouncement | null>(null)
const announcementVisible = ref(false)

const openQuickBuy = (product: any) => {
  quickBuyProduct.value = product
  quickBuyVisible.value = true
}

// ==================== Banner Carousel ====================
const {
  banners,
  bannerLoading,
  currentBannerIndex,
  bannerCount,
  showHeroSection,
  heroImage,
  heroBadge,
  heroTitle,
  heroSubtitle,
  heroPrimaryButtonText,
  loadBanners,
  handleNextHeroBanner,
  handlePrevHeroBanner,
  selectHeroBanner,
  goToHeroLink,
  onBannerTouchStart,
  onBannerTouchEnd,
  stopHeroAutoPlay,
} = useBannerCarousel()

// ==================== List Mode ====================
const {
  loading: listLoading,
  products: listProducts,
  selectedCategory: listSelectedCategory,
  searchQuery: listSearchQuery,
  currentPage: listCurrentPage,
  totalPages: listTotalPages,
  showFilterDrawer: listShowFilterDrawer,
  expandedParentIds: listExpandedParentIds,
  categoryGroups: listCategoryGroups,
  categoryMap: listCategoryMap,
  selectCategory: listSelectCategory,
  toggleParentCategory: listToggleParentCategory,
  changePage: listChangePage,
  clearSearch: listClearSearch,
  onSearch: listOnSearch,
  initialize: listInitialize,
  cleanup: listCleanup,
} = useProductList({ pageSize: 20, homeRouteName: 'home' })

const listProductGroups = useProductListGroups(listProducts, listCategoryMap)

// ==================== SEO ====================
const route = useRoute()
const seoCategoryName = computed(() => {
  if (!listSelectedCategory.value) return ''
  const cat = listCategoryMap.value.get(listSelectedCategory.value)
  return cat ? getLocalizedText(cat.name) : ''
})
usePageSeo({
  canonicalPath: () => route.path,
  title: () => {
    if (route.name === 'category-products') {
      return seoCategoryName.value || t('nav.products')
    }
    if (route.name === 'products') return t('nav.products')
    return undefined
  },
})

// ==================== Card Mode ====================
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}

const goToProduct = (slug: string) => {
  router.push(`/products/${slug}`)
}

const goToPost = (slug: string) => {
  router.push(`/blog/${slug}`)
}

const loadFeaturedProducts = async () => {
  try {
    const response = await productAPI.list({ page: 1, page_size: 15 })
    products.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load products:', error)
  } finally {
    productsLoading.value = false
  }
}

const loadLatestPosts = async () => {
  if (!latestSectionVisible.value) return
  try {
    const params: Record<string, unknown> = { page: 1, page_size: 3 }
    if (blogEnabled.value && !noticeEnabled.value) params.type = 'blog'
    if (!blogEnabled.value && noticeEnabled.value) params.type = 'notice'
    const response = await postAPI.list(params)
    posts.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load posts:', error)
  }
}

// ==================== Lifecycle ====================
const showAnnouncementIfNeeded = () => {
  const announcement = appStore.config?.announcement as HomeAnnouncement | undefined
  if (announcement && shouldShow(announcement)) {
    activeAnnouncement.value = announcement
    announcementVisible.value = true
  }
}

onMounted(async () => {
  await appStore.loadConfig()
  if (templateMode.value === 'list') {
    await Promise.all([loadBanners(), listInitialize()])
  } else {
    await Promise.all([loadBanners(), loadFeaturedProducts(), loadLatestPosts()])
  }
  showAnnouncementIfNeeded()
})

onUnmounted(() => {
  stopHeroAutoPlay()
  listCleanup()
})
</script>

<style scoped>
.banner-fade-enter-active,
.banner-fade-leave-active {
  transition: opacity 300ms ease;
}
.banner-fade-enter-from,
.banner-fade-leave-to {
  opacity: 0;
}
</style>
