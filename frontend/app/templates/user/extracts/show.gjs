// app/templates/user/extracts/show.gjs

import { LinkTo } from '@ember/routing';

function formatLines(text) {
  if (!text) return [];
  return text.split(/\r?\n/).map((line, index) => ({
    number: index + 1,
    text: line
  }));
}

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <header class="mb-6">
        <div class="flex items-center gap-4">
          <LinkTo @route="user.extracts" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
            </svg>
          </LinkTo>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            {{@model.document.documentName}}
          </h1>
        </div>
      </header>

      <div class="rounded-lg bg-white shadow ring-1 ring-gray-200 dark:bg-gray-800 dark:ring-gray-700">
        <div class="overflow-x-auto p-4">
          <table class="min-w-full text-sm font-mono">
            <tbody>
              {{#each (formatLines @model.content) as |line|}}
                <tr>
                  <td class="w-12 select-none border-r border-gray-200 pr-4 text-right text-gray-400 py-0 align-baseline leading-tight dark:border-gray-700 dark:text-gray-500">
                    {{line.number}}
                  </td>
                  <td class="whitespace-pre pl-4 text-gray-900 py-0 align-baseline leading-tight dark:text-gray-100">
                    {{line.text}}
                  </td>
                </tr>
              {{/each}}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
