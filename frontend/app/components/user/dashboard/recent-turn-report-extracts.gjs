// app/components/user/dashboard/recent-turn-report-extracts.gjs

import Component from "@glimmer/component";
import { LinkTo } from "@ember/routing";

import DashboardTurnReportStatus from 'frontend/components/user/dashboard/turn-report-status';

export default class RecentTurnReportExtracts extends Component {
  <template>
    <section class="flex h-full flex-col rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200 dark:bg-gray-800 dark:ring-gray-700">
      <header class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          Recent Turn Report Extracts
        </h2>

        {{!-- Optional overflow menu --}}
        <span class="inline-flex h-8 w-8 items-center justify-center rounded-full text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700">
          <span class="sr-only">More</span>
          •••
        </span>
      </header>

      {{#if @files.length}}
        <ul class="flex-1 space-y-3">
          {{#each @files as |file|}}
            <li class="rounded-md border border-gray-100 bg-gray-50 px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-900">
              <div class="flex items-center justify-between gap-4">
                <div class="min-w-0">
                  <p class="truncate text-sm font-medium text-gray-900 dark:text-white">
                    <LinkTo @route="user.extracts.show" @model={{file.id}}
                       class="text-sm font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300">
                      {{file.documentName}}
                    </LinkTo>
                  </p>
                  <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                    {{file.createdAt}}
                  </p>
                </div>

                <DashboardTurnReportStatus
                  @status={{file.processingStatus}}
                />
              </div>
            </li>
          {{/each}}
        </ul>
      {{else}}
        <p class="flex-1 text-sm text-gray-500 dark:text-gray-400">
          You don’t have any turn report extracts yet.
        </p>
      {{/if}}

      <footer class="mt-4 border-t border-gray-100 pt-4 dark:border-gray-700">
        <LinkTo
          @route="user.extracts"
          class="text-sm font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300"
        >
          View all turn report extracts &rarr;
        </LinkTo>
      </footer>
    </section>
  </template>
}
