// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import NavList from 'frontend/components/sidebar/nav-list.gjs';

export default class SidebarRelative extends Component {
  <template>
    <el-dialog>
      <dialog id="sidebar" class="backdrop:bg-transparent lg:hidden">
        <el-dialog-backdrop class="fixed inset-0 bg-gray-900/80 transition-opacity duration-300 ease-linear data-closed:opacity-0"></el-dialog-backdrop>

        <div tabindex="0" class="fixed inset-0 flex focus:outline-none">
          <el-dialog-panel class="group/dialog-panel relative mr-16 flex w-full max-w-xs flex-1 transform transition duration-300 ease-in-out data-closed:-translate-x-full">
            <div class="absolute top-0 left-full flex w-16 justify-center pt-5 duration-300 ease-in-out group-data-closed/dialog-panel:opacity-0">
              <button type="button" command="close" commandfor="sidebar" class="-m-2.5 p-2.5">
                <span class="sr-only">Close sidebar</span>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon" aria-hidden="true" class="size-6 text-white">
                  <path d="M6 18 18 6M6 6l12 12" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>

            <!-- Sidebar component, swap this element with another sidebar if you like -->
            <div class="relative flex grow flex-col gap-y-5 overflow-y-auto bg-white px-6 pb-4 dark:bg-gray-900 dark:ring dark:ring-white/10 dark:before:pointer-events-none dark:before:absolute dark:before:inset-0 dark:before:bg-black/10">
              <div class="relative flex h-16 shrink-0 items-center">
                <img src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600" alt="Your Company" class="h-8 w-auto dark:hidden" />
                <img src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500" alt="Your Company" class="h-8 w-auto not-dark:hidden" />
              </div>
              <nav class="relative flex flex-1 flex-col">
                <NavList />
              </nav>
            </div>
          </el-dialog-panel>
        </div>
      </dialog>
    </el-dialog>
  </template>
}
