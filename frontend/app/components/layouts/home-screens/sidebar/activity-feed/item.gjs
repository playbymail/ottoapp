// app/components/layouts/home-screens/sidebar/activity-feed/item.gjs
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

<template>
  <li class="px-4 py-4 sm:px-6 lg:px-8">
    <div class="flex items-center gap-x-3">
      <img src="{{@image}}" alt="" class="size-6 flex-none rounded-full bg-gray-100 outline -outline-offset-1 outline-black/5" />
      <h3 class="flex-auto truncate text-sm/6 font-semibold text-gray-900">{{@actor}}</h3>
      <time datetime="2023-01-23T11:00" class="flex-none text-xs text-gray-500">{{@time}}</time>
    </div>
    {{yield}}
  </li>
</template>
