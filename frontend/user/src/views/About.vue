<template>
  <div
    class="about-page min-h-screen bg-background text-foreground pt-20 pb-16">
    <div class="container mx-auto px-4 max-w-4xl">
      <div class="text-center mb-16 mt-12">
        <h1 class="text-4xl md:text-6xl font-black mb-6 tracking-tight text-foreground">{{ heroTitle }}</h1>
        <p
          class="text-muted-foreground max-w-2xl mx-auto text-lg leading-relaxed border-b pb-8">
          {{ heroSubtitle }}
        </p>
      </div>

      <Card
        class="backdrop-blur-xl rounded-3xl p-8 md:p-12 shadow-2xl relative overflow-hidden">

        <div v-if="hasIntroduction" class="mb-12">
          <p class="text-muted-foreground text-lg leading-relaxed whitespace-pre-line">
            {{ introductionText }}
          </p>
        </div>

        <div v-if="hasServices" class="mb-12">
          <h2 v-if="servicesTitle" class="text-2xl font-bold text-foreground mb-6 flex items-center gap-3">
            <span class="w-1.5 h-8 bg-primary rounded-full"></span>
            {{ servicesTitle }}
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="(service, index) in serviceItems"
              :key="`about-service-${index}`"
              class="flex items-start gap-3 p-4 bg-secondary rounded-xl border transition-colors">
              <Check class="w-6 h-6 text-primary flex-shrink-0 mt-0.5" />
              <span class="text-muted-foreground">{{ service }}</span>
            </div>
          </div>
        </div>

        <div v-if="hasContact">
          <h2 v-if="contactTitle" class="text-2xl font-bold text-foreground mb-6 flex items-center gap-3">
            <span class="w-1.5 h-8 bg-primary rounded-full"></span>
            {{ contactTitle }}
          </h2>
          <p v-if="contactText" class="text-muted-foreground mb-8 whitespace-pre-line">
            {{ contactText }}
          </p>
          <div v-if="hasContactLinks" class="max-w-sm">
            <component :is="storeContacts.telegram.href ? 'a' : 'div'" :href="storeContacts.telegram.href || undefined"
              :target="storeContacts.telegram.href ? '_blank' : undefined" rel="noopener noreferrer"
              class="group flex items-center justify-center space-x-3 bg-secondary border text-muted-foreground px-6 py-4 rounded-xl transition-colors hover:text-foreground">
              <svg class="w-6 h-6 text-sky-500" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.295-.6.295-.002 0-.003 0-.005 0l.213-3.054 5.56-5.022c.24-.213-.054-.334-.373-.121l-6.869 4.326-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.538-.196 1.006.128.832.941z" />
              </svg>
              <span class="font-bold">{{ storeContacts.telegram.label }}</span>
            </component>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Check } from 'lucide-vue-next'
import { Card } from '@/components/ui/card'
import { useAbout } from '../composables/useAbout'

const {
  storeContacts, heroTitle, heroSubtitle, introductionText, servicesTitle, contactTitle, contactText,
  serviceItems, hasIntroduction, hasServices, hasContactLinks, hasContact,
} = useAbout()
</script>
