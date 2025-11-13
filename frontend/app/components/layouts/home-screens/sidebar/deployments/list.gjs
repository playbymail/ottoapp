// app/components/layouts/home-screens/sidebar/deployments/list.gjs
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

<template>
  <ul role="list" class="divide-y divide-gray-100">
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-gray-100 p-1 text-gray-400">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Planetaria</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">ios-app</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Initiated 1m 32s ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-green-500/10 p-1 text-green-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Planetaria</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">mobile-api</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 3m ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-indigo-50 px-2 py-1 text-xs font-medium text-indigo-500 ring-1 ring-indigo-200 ring-inset">Production</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-gray-100 p-1 text-gray-400">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Tailwind Labs</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">tailwindcss.com</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 3h ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-green-500/10 p-1 text-green-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Tailwind Labs</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">company-website</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 1d ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-green-500/10 p-1 text-green-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Protocol</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">relay-service</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 1d ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-indigo-50 px-2 py-1 text-xs font-medium text-indigo-500 ring-1 ring-indigo-200 ring-inset">Production</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-green-500/10 p-1 text-green-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Planetaria</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">android-app</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 5d ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-rose-500/10 p-1 text-rose-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Protocol</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">api.protocol.chat</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Failed to deploy 6d ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
    <li class="relative flex items-center space-x-4 px-4 py-4 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-auto">
        <div class="flex items-center gap-x-3">
          <div class="flex-none rounded-full bg-green-500/10 p-1 text-green-500">
            <div class="size-2 rounded-full bg-current"></div>
          </div>
          <h2 class="min-w-0 text-sm/6 font-semibold text-gray-900">
            <a href="#" class="flex gap-x-2">
              <span class="truncate">Planetaria</span>
              <span class="text-gray-400">/</span>
              <span class="whitespace-nowrap">planetaria.tech</span>
              <span class="absolute inset-0"></span>
            </a>
          </h2>
        </div>
        <div class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-gray-500">
          <p class="truncate">Deploys from GitHub</p>
          <svg viewBox="0 0 2 2" class="size-0.5 flex-none fill-gray-300">
            <circle r="1" cx="1" cy="1" />
          </svg>
          <p class="whitespace-nowrap">Deployed 6d ago</p>
        </div>
      </div>
      <div class="flex-none rounded-full bg-gray-50 px-2 py-1 text-xs font-medium text-gray-500 ring-1 ring-gray-200 ring-inset">Preview</div>
      <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 flex-none text-gray-400">
        <path d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
      </svg>
    </li>
  </ul>
</template>
