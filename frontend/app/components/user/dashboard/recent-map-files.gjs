// app/components/user/dashboard/recent-map-files.gjs

import Component from "@glimmer/component";
import {LinkTo} from "@ember/routing";

export default class RecentMapFiles extends Component {
  <template>
    <section class="flex h-full flex-col rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
      <header class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-900">
          Recent Map Files
        </h2>

        {{!-- Optional overflow menu --}}
        <span class="inline-flex h-8 w-8 items-center justify-center rounded-full text-gray-400 hover:bg-gray-100">
          <span class="sr-only">More</span>
          •••
        </span>
      </header>

      {{#if @files.length}}
        <ul class="flex-1 space-y-3">
          {{#each @files as |file|}}
            <li
              class="flex items-center justify-between gap-4 rounded-md border border-gray-100 bg-gray-50 px-3 py-2 text-sm">
              <div class="min-w-0">
                <p class="truncate text-sm font-medium text-gray-900">
                  {{file.documentName}}
                </p>
                <p class="mt-0.5 text-xs text-gray-500">
                  {{file.createdAt}}
                </p>
              </div>

              <div class="shrink-0">
                <a
                  href={{file.links.contents.href}}
                  class="inline-flex items-center rounded-md border border-indigo-500 px-3 py-1.5 text-xs font-medium text-indigo-600 hover:bg-indigo-50"
                >
                  Download
                </a>
              </div>
            </li>
          {{/each}}
        </ul>
      {{else}}
        <p class="flex-1 text-sm text-gray-500">
          You don’t have any map files yet.
        </p>
      {{/if}}

      <footer class="mt-4 border-t border-gray-100 pt-4">
        <LinkTo
          @route="user.maps"
          class="text-sm font-medium text-indigo-600 hover:text-indigo-500"
        >
          View all map files &rarr;
        </LinkTo>
      </footer>
    </section>
  </template>
}
