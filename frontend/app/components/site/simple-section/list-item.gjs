// app/components/site/simple-section/item.gjs
// Requires a TailwindCSS Plus license.
// https://tailwindcss.com/plus/ui-blocks/marketing/sections/feature-sections#simple

<template>
  <div>
    <dt class="font-semibold text-gray-900">{{@title}}</dt>
    <dd class="mt-1 text-gray-600">{{yield}}</dd>
  </div>
</template>
