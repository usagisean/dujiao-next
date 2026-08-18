<template>
  <div class="mx-auto w-full max-w-[840px] px-6 pb-10 pt-2">
    <header class="my-8 text-center">
      <h1 class="mb-2.5 text-4xl font-extrabold">{{ heroTitle }}</h1>
      <p class="text-muted-foreground">{{ heroSubtitle }}</p>
    </header>

    <Card class="p-7 sm:p-[34px]">
      <div v-if="hasIntroduction" class="mb-9 last:mb-0">
        <p class="whitespace-pre-line text-base leading-[1.8] text-muted-foreground">{{ introductionText }}</p>
      </div>

      <div v-if="hasServices" class="mb-9 last:mb-0">
        <h2 v-if="servicesTitle" class="mb-[18px] flex items-center gap-3 text-[22px] font-bold">
          <span class="h-[26px] w-[5px] flex-none rounded-full bg-primary"></span>{{ servicesTitle }}
        </h2>
        <div class="grid gap-3.5 sm:grid-cols-2">
          <div v-for="(service, index) in serviceItems" :key="`about-service-${index}`" class="flex items-start gap-3 rounded-md border bg-secondary px-4 py-3.5">
            <Check class="mt-0.5 h-5 w-5 flex-none text-[color:var(--teal-strong)]" />
            <span class="text-[14.5px] text-muted-foreground">{{ service }}</span>
          </div>
        </div>
      </div>

      <div v-if="hasContact" class="mb-9 last:mb-0">
        <h2 v-if="contactTitle" class="mb-[18px] flex items-center gap-3 text-[22px] font-bold">
          <span class="h-[26px] w-[5px] flex-none rounded-full bg-primary"></span>{{ contactTitle }}
        </h2>
        <p v-if="contactText" class="mb-[18px] whitespace-pre-line leading-relaxed text-muted-foreground">{{ contactText }}</p>
        <div v-if="hasContactLinks" class="max-w-sm">
          <component :is="storeContacts.telegram.href ? 'a' : 'div'" :href="storeContacts.telegram.href || undefined" :target="storeContacts.telegram.href ? '_blank' : undefined" rel="noopener noreferrer" class="flex items-center justify-center gap-3 rounded-md border bg-secondary p-4 font-bold text-muted-foreground transition hover:border-hairline-strong hover:text-foreground">
            <svg viewBox="0 0 24 24" fill="currentColor" class="h-6 w-6 flex-none text-[#0ea5e9]"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.295-.6.295-.002 0-.003 0-.005 0l.213-3.054 5.56-5.022c.24-.213-.054-.334-.373-.121l-6.869 4.326-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.538-.196 1.006.128.832.941z" /></svg>
            <span>{{ storeContacts.telegram.label }}</span>
          </component>
        </div>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { Check } from 'lucide-vue-next'
import { Card } from '@/components/ui/card'
import { useAbout } from '../../composables/useAbout'

const {
  storeContacts, heroTitle, heroSubtitle, introductionText, servicesTitle, contactTitle, contactText,
  serviceItems, hasIntroduction, hasServices, hasContactLinks, hasContact,
} = useAbout()
</script>
