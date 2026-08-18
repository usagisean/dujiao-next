import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createHead } from '@unhead/vue/client'
import './style.css'
import App from './App.vue'
import router, { warmupCommonRoutes } from './router'
import i18n, { detectLocale, setI18nLocale, warmupLocaleMessages } from './i18n'
import { useTelegramMiniAppStore } from './stores/telegramMiniApp'
import { initTemplateOverride } from './templates/registry'

// 预览用：?template=vault 持久化激活模板（站长正式切换走站点配置）
initTemplateOverride()

const brandLog = (globalThis as any).console?.log?.bind(console)
brandLog?.(
  '%c ZX AI Hub %c Storefront ',
  'background:#5b5ce2;color:#fff;padding:4px 8px;border-radius:4px 0 0 4px;font-weight:bold;',
  'background:#101522;color:#f5f5f7;padding:4px 8px;border-radius:0 4px 4px 0;',
)

const app = createApp(App)
const head = createHead()
const pinia = createPinia()

app.use(pinia)
app.use(head)
app.use(router)
app.use(i18n)

// 非默认语言的语言包为懒加载 chunk，挂载前并行加载，避免首屏文案闪现兜底语言
Promise.all([
  useTelegramMiniAppStore(pinia).init(),
  setI18nLocale(detectLocale()),
]).then(() => {
  app.mount('#app')
})

void router.isReady().then(() => {
    warmupCommonRoutes()
    warmupLocaleMessages()
})
