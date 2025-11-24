// app/templates/user/reports.gjs
import { LinkTo } from '@ember/routing';

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <div class="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <header class="mb-6 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <LinkTo @route="user.dashboard" class="text-gray-500 hover:text-gray-700 xl:hidden dark:text-gray-400 dark:hover:text-gray-300">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
            </svg>
          </LinkTo>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            Reports
          </h1>
        </div>
      </header>

      <section class="flex h-full flex-col rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200 dark:bg-gray-800 dark:ring-gray-700">
        {{#if @model.length}}
          <ul class="flex-1 space-y-3">
            {{#each @model as |file|}}
              <li class="flex items-center justify-between gap-4 rounded-md border border-gray-100 bg-gray-50 px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900">
                <div class="min-w-0">
                  <p class="truncate text-sm font-medium text-gray-900 dark:text-white">
                    {{file.documentName}}
                  </p>
                  <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                    {{file.updatedAt}}
                  </p>
                </div>

                <div class="shrink-0">
                  <a
                    href={{file.links.contents.href}}
                    class="inline-flex items-center rounded-md border border-indigo-500 px-3 py-1.5 text-xs font-medium text-indigo-600 hover:bg-indigo-50 dark:border-indigo-400 dark:text-indigo-400 dark:hover:bg-indigo-900/20"
                  >
                    Download
                  </a>
                </div>
              </li>
            {{/each}}
          </ul>
        {{else}}
          <p class="flex-1 text-sm text-gray-500 dark:text-gray-400">
            You don’t have any turn reports yet.
          </p>
        {{/if}}
      </section>
    </div>
  </div>
</template>
