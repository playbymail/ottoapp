// app/components/site/simple-section/heading.gjs
// Requires a TailwindCSS Plus license.
// https://tailwindcss.com/plus/ui-blocks/marketing/sections/feature-sections#simple

<template>
  <div class="mx-auto max-w-2xl lg:mx-0">
    <h2 class="text-4xl font-semibold tracking-tight text-pretty text-gray-900 sm:text-5xl">{{@title}}</h2>
    <p class="mt-6 text-lg/8 text-gray-700">{{yield}}</p>
  </div>
</template>
