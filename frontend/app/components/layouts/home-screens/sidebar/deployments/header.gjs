// app/components/layouts/home-screens/sidebar/deployments/header.gjs
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

<template>
  <header class="flex items-center justify-between border-b border-gray-200 px-4 py-4 sm:px-6 sm:py-6 lg:px-8">
    <h1 class="text-base/7 font-semibold text-gray-900">Deployments</h1>

    <!-- Sort dropdown -->
    <el-dropdown class="relative">
      <button class="flex items-center gap-x-1 text-sm/6 font-medium text-gray-900">
        Sort by
        <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 text-gray-500">
          <path d="M10.53 3.47a.75.75 0 0 0-1.06 0L6.22 6.72a.75.75 0 0 0 1.06 1.06L10 5.06l2.72 2.72a.75.75 0 1 0 1.06-1.06l-3.25-3.25Zm-4.31 9.81 3.25 3.25a.75.75 0 0 0 1.06 0l3.25-3.25a.75.75 0 1 0-1.06-1.06L10 14.94l-2.72-2.72a.75.75 0 0 0-1.06 1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
        </svg>
      </button>
      <el-menu anchor="bottom end" popover class="w-40 origin-top-right rounded-md bg-white py-2 shadow-lg outline-1 outline-gray-900/5 transition transition-discrete [--anchor-gap:--spacing(2.5)] data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in">
        <a href="#" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden">Name</a>
        <a href="#" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden">Date updated</a>
        <a href="#" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden">Environment</a>
      </el-menu>
    </el-dropdown>
  </header>
</template>
