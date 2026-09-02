import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { bannerAPI } from '../api'
import { getImageUrl } from '../utils/image'
import { useTheme } from '../utils/theme'
import { useLocalized } from './useProduct'

const localized = (simplified: string, traditional: string, english: string) => ({
  'zh-CN': simplified,
  'zh-TW': traditional,
  'en-US': english,
})

const DEFAULT_TECH_BANNERS = [
  {
    id: 'zx-nexus',
    image: '/banners/zx-nexus.svg',
    image_dark: '/banners/zx-nexus.svg',
    image_light: '/banners/zx-nexus-light.svg',
    badge: localized('新一代精选平台', '新一代精選平台', 'NEXT-GEN MARKETPLACE'),
    title: localized('ZX NEXUS', 'ZX NEXUS', 'ZX NEXUS'),
    subtitle: localized('探索 · 连接 · 体验', '探索 · 連結 · 體驗', 'DISCOVER · CONNECT · EXPERIENCE'),
    link_type: 'internal',
    link_value: '/products',
  },
  {
    id: 'zx-spectrum',
    image: '/banners/zx-spectrum.svg',
    image_dark: '/banners/zx-spectrum.svg',
    image_light: '/banners/zx-spectrum-light.svg',
    badge: localized('全球服务在线', '全球服務在線', 'SYSTEMS ONLINE'),
    title: localized('连接世界', '連結世界', 'GLOBAL ACCESS'),
    subtitle: localized('精选商品 · 在线服务', '精選商品 · 線上服務', 'CURATED PRODUCTS · DIGITAL SERVICES'),
    link_type: 'internal',
    link_value: '/products',
  },
  {
    id: 'zx-flow',
    image: '/banners/zx-flow.svg',
    image_dark: '/banners/zx-flow.svg',
    image_light: '/banners/zx-flow-light.svg',
    badge: localized('自动化交付', '自動化交付', 'AUTOMATED FULFILLMENT'),
    title: localized('即时交付', '即時交付', 'INSTANT DELIVERY'),
    subtitle: localized('快速 · 清晰 · 稳定', '快速 · 清晰 · 穩定', 'FAST · CLEAR · RELIABLE'),
    link_type: 'internal',
    link_value: '/products',
  },
]

export function useBannerCarousel() {
  const router = useRouter()
  const { t } = useI18n()
  const { getLocalizedText } = useLocalized()
  const { theme } = useTheme()

  const banners = ref<any[]>([])
  const bannerLoading = ref(true)
  const currentBannerIndex = ref(0)
  let heroAutoPlayTimer: ReturnType<typeof setInterval> | null = null

  // Touch swipe
  let touchStartX = 0
  const onBannerTouchStart = (e: TouchEvent) => {
    touchStartX = e.touches[0]?.clientX ?? 0
  }
  const onBannerTouchEnd = (e: TouchEvent) => {
    const diff = touchStartX - (e.changedTouches[0]?.clientX ?? 0)
    if (Math.abs(diff) > 50) {
      if (diff > 0) handleNextHeroBanner()
      else handlePrevHeroBanner()
    }
  }

  const bannerCount = computed(() => banners.value.length)
  const showHeroSection = computed(() => bannerLoading.value || bannerCount.value > 0)

  const heroBanner = computed(() => {
    if (banners.value.length === 0) return null
    if (currentBannerIndex.value >= banners.value.length) return banners.value[0]
    return banners.value[currentBannerIndex.value]
  })

  const goToBanner = (index: number) => {
    if (banners.value.length === 0) return
    const total = banners.value.length
    currentBannerIndex.value = ((index % total) + total) % total
  }

  const nextHeroBanner = () => goToBanner(currentBannerIndex.value + 1)
  const prevHeroBanner = () => goToBanner(currentBannerIndex.value - 1)

  const stopHeroAutoPlay = () => {
    if (heroAutoPlayTimer) {
      clearInterval(heroAutoPlayTimer)
      heroAutoPlayTimer = null
    }
  }

  const startHeroAutoPlay = () => {
    stopHeroAutoPlay()
    if (banners.value.length <= 1) return
    heroAutoPlayTimer = setInterval(() => {
      nextHeroBanner()
    }, 5000)
  }

  const handleNextHeroBanner = () => {
    nextHeroBanner()
    startHeroAutoPlay()
  }

  const handlePrevHeroBanner = () => {
    prevHeroBanner()
    startHeroAutoPlay()
  }

  const selectHeroBanner = (index: number) => {
    goToBanner(index)
    startHeroAutoPlay()
  }

  const heroImage = computed(() => {
    const banner = heroBanner.value
    if (!banner) return ''
    const isMobile = typeof window !== 'undefined' && window.innerWidth < 768
    const themedDesktopImage = theme.value === 'light'
      ? banner.image_light
      : banner.image_dark
    const themedMobileImage = theme.value === 'light'
      ? banner.mobile_image_light
      : banner.mobile_image_dark
    const selectedImage = isMobile
      ? (themedMobileImage || banner.mobile_image || themedDesktopImage || banner.image)
      : (themedDesktopImage || banner.image || themedMobileImage || banner.mobile_image)
    return getImageUrl(selectedImage || '')
  })

  const heroBadge = computed(() => {
    if (!heroBanner.value) return t('home.hero.badge')
    return getLocalizedText(heroBanner.value.badge) || t('home.hero.badge')
  })

  const heroTitle = computed(() => {
    if (!heroBanner.value) return t('home.hero.title')
    return getLocalizedText(heroBanner.value.title) || t('home.hero.title')
  })

  const heroSubtitle = computed(() => {
    if (!heroBanner.value) return t('home.hero.subtitle')
    return getLocalizedText(heroBanner.value.subtitle) || t('home.hero.subtitle')
  })

  const heroLink = computed(() => {
    const banner = heroBanner.value
    if (!banner || banner.link_type === 'none') return ''
    return banner.link_value || ''
  })

  const hasHeroLink = computed(() => heroLink.value.trim().length > 0)

  const heroPrimaryButtonText = computed(() => {
    if (!hasHeroLink.value) return t('home.hero.cta')
    const linkType = String(heroBanner.value?.link_type || '').toLowerCase()
    if (linkType === 'external') {
      return t('common.learnMore')
    }
    return t('common.viewDetails')
  })

  const heroOpenInNewTab = computed(() => Boolean(heroBanner.value?.open_in_new_tab))

  const goToHeroLink = () => {
    if (!heroLink.value) {
      router.push('/products')
      return
    }
    const isExternal = /^https?:\/\//i.test(heroLink.value)
    if (isExternal || heroOpenInNewTab.value) {
      window.open(heroLink.value, heroOpenInNewTab.value ? '_blank' : '_self')
      return
    }
    router.push(heroLink.value)
  }

  const loadBanners = async () => {
    bannerLoading.value = true
    try {
      const response = await bannerAPI.list({ position: 'home_hero', limit: 5 })
      const configuredBanners = response.data.data || []
      banners.value = configuredBanners.length > 0
        ? configuredBanners
        : DEFAULT_TECH_BANNERS.map((banner) => ({ ...banner }))
      currentBannerIndex.value = 0
      startHeroAutoPlay()
    } catch (error) {
      banners.value = DEFAULT_TECH_BANNERS.map((banner) => ({ ...banner }))
      currentBannerIndex.value = 0
      startHeroAutoPlay()
      console.error('Failed to load banners:', error)
    } finally {
      bannerLoading.value = false
    }
  }

  return {
    banners,
    bannerLoading,
    currentBannerIndex,
    bannerCount,
    showHeroSection,
    heroBanner,
    heroImage,
    heroBadge,
    heroTitle,
    heroSubtitle,
    heroLink,
    hasHeroLink,
    heroPrimaryButtonText,
    heroOpenInNewTab,
    loadBanners,
    handleNextHeroBanner,
    handlePrevHeroBanner,
    selectHeroBanner,
    goToHeroLink,
    onBannerTouchStart,
    onBannerTouchEnd,
    stopHeroAutoPlay,
  }
}
