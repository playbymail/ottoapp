// Copyright (c) 2025 Michael D Henderson. All rights reserved.
import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import {LinkTo} from '@ember/routing';

export default class Header extends Component {
  <template>
    <header class="absolute inset-x-0 top-0 z-50">
      <nav aria-label="Global" class="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8">
        <div class="flex lg:flex-1">
          <a href="/" class="-m-1.5 p-1.5">
            <span class="sr-only">OttoMap</span>
            <img src="/img/logo-light.svg" alt="OttoApp" class="h-8 w-auto dark:hidden" />
            <img src="/img/logo-dark.svg" alt="OttoApp" class="h-8 w-auto not-dark:hidden" />
          </a>
        </div>
        <div class="flex lg:hidden">
          <button type="button" command="show-modal" commandfor="mobile-menu"
                  class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700">
            <span class="sr-only">Open main menu</span>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
                 aria-hidden="true" class="size-6">
              <path d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
        </div>
        <div class="hidden lg:flex lg:gap-x-12">
          <a href="/docs" class="text-sm/6 font-semibold text-gray-900">OttoMap</a>
          <a href="https://tribenet.com.au/" class="text-sm/6 font-semibold text-gray-900">TribeNet</a>
          <a href="https://worldographer.com/" class="text-sm/6 font-semibold text-gray-900">Worldographer</a>
        </div>
        <div class="hidden lg:flex lg:flex-1 lg:justify-end">
          <a href="/login" class="text-sm/6 font-semibold text-gray-900">Log in <span aria-hidden="true">&rarr;</span></a>
        </div>
      </nav>
      <el-dialog>
        <dialog id="mobile-menu" class="backdrop:bg-transparent lg:hidden">
          <div tabindex="0" class="fixed inset-0 focus:outline-none">
            <el-dialog-panel
              class="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-white p-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10">
              <div class="flex items-center justify-between">
                <a href="/" class="-m-1.5 p-1.5">
                  <span class="sr-only">OttoMap</span>
                  <img src="/img/logo-light.svg" alt="OttoApp" class="h-8 w-auto dark:hidden" />
                  <img src="/img/logo-dark.svg" alt="OttoApp" class="h-8 w-auto not-dark:hidden" />
                </a>
                <button type="button" command="close" commandfor="mobile-menu"
                        class="-m-2.5 rounded-md p-2.5 text-gray-700">
                  <span class="sr-only">Close menu</span>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
                       aria-hidden="true" class="size-6">
                    <path d="M6 18 18 6M6 6l12 12" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
              </div>
              <div class="mt-6 flow-root">
                <div class="-my-6 divide-y divide-gray-500/10">
                  <div class="space-y-2 py-6">
                    <a href="/docs"
                       class="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-gray-900 hover:bg-gray-50">
                      OttoMap</a>
                    <a href="https://tribenet.com.au/"
                       class="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-gray-900 hover:bg-gray-50">
                      TribeNet</a>
                    <a href="https://worldographer.com/"
                       class="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-gray-900 hover:bg-gray-50">
                      Worldographer</a>
                  </div>
                  <div class="py-6">
                    <a href="/login"
                       class="-mx-3 block rounded-lg px-3 py-2.5 text-base/7 font-semibold text-gray-900 hover:bg-gray-50">
                      Log in</a>
                  </div>
                </div>
              </div>
            </el-dialog-panel>
          </div>
        </dialog>
      </el-dialog>
    </header>
  </template>
}
