<template>
  <section class="container mx-auto max-w-3xl px-4 pb-20 pt-28 text-foreground">
    <p class="mb-3 text-sm font-semibold text-primary">{{ brandName }}</p>
    <h1 class="mb-8 text-3xl font-bold tracking-tight sm:text-4xl">{{ t('afterSales.title') }}</h1>
    <div class="space-y-6 rounded-2xl border bg-card p-6 text-muted-foreground sm:p-8">
      <p>{{ t('afterSales.intro') }}</p>
      <section class="space-y-2">
        <h2 class="text-lg font-semibold text-foreground">{{ t('afterSales.checkoutTitle') }}</h2>
        <p>{{ t('afterSales.checkoutText') }}</p>
        <RouterLink to="/guest/orders" class="inline-block font-medium text-primary underline underline-offset-4">{{ t('navbar.guestOrders') }}</RouterLink>
      </section>
      <section class="space-y-2">
        <h2 class="text-lg font-semibold text-foreground">{{ t('afterSales.requestTitle') }}</h2>
        <p>{{ t('afterSales.requestText') }}</p>
        <p>{{ t('afterSales.privacyHint') }}</p>
        <div class="flex flex-wrap gap-4 pt-2">
          <a v-if="supportEmail" :href="`mailto:${supportEmail}`" class="font-medium text-primary underline underline-offset-4">{{ supportEmail }}</a>
          <a v-if="telegram.href" :href="telegram.href" target="_blank" rel="noopener noreferrer" class="font-medium text-primary underline underline-offset-4">{{ telegram.label }}</a>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/app'
import { usePageSeo } from '../composables/usePageSeo'
import { resolveBrandName } from '../config/brand'
import { resolveStoreContacts, supportEmail as mainSupportEmail } from '../config/contact'

const { t } = useI18n()
const app = useAppStore()
const brandName = computed(() => resolveBrandName(app.config?.brand?.site_name))
const supportEmail = computed(() => app.config && !app.isResellerTenant ? mainSupportEmail : '')
const telegram = computed(() => resolveStoreContacts(app.config?.contact).telegram)
usePageSeo({
  title: () => t('afterSales.title'),
  description: () => t('afterSales.intro'),
  canonicalPath: () => '/refund',
})
</script>
