// app/components/site/simple-section/list.gjs
// Requires a TailwindCSS Plus license.
// https://tailwindcss.com/plus/ui-blocks/marketing/sections/feature-sections#simple

<template>
  <dl
    class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 text-base/7 sm:grid-cols-2 lg:mx-0 lg:max-w-none lg:grid-cols-3">
    {{yield}}
  </dl>
</template>
