// app/components/layouts/home-screens/sidebar/sidebar.js
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

import SidebarNav from './sidebar-nav';

<template>
  <!-- mobile -->
  <el-dialog>
    <dialog id="sidebar" class="backdrop:bg-transparent xl:hidden">
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
          <div class="relative flex grow flex-col gap-y-5 overflow-y-auto bg-gray-50 px-6">
            <div class="relative flex h-16 shrink-0 items-center">
              <img src="{{@logoLight}}" alt="OttoApp" class="h-8 w-auto dark:hidden" />
              <img src="{{@logoDark}}" alt="OttoApp" class="h-8 w-auto not-dark:hidden" />
            </div>
            <SidebarNav @isRelative=true />
          </div>
        </el-dialog-panel>
      </div>
    </dialog>
  </el-dialog>

  <!-- desktop -->
  <div class="hidden xl:fixed xl:inset-y-0 xl:z-50 xl:flex xl:w-72 xl:flex-col">
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-50 px-6 ring-1 ring-gray-200">
      <div class="flex h-16 shrink-0 items-center">
        <img src="{{@logoLight}}" alt="OttoApp" class="h-8 w-auto dark:hidden" />
        <img src="{{@logoDark}}" alt="OttoApp" class="h-8 w-auto not-dark:hidden" />
      </div>
      <SidebarNav @isRelative=false />
    </div>
  </div>
</template>
