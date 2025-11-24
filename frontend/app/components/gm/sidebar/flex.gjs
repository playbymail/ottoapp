// app/components/gm/sidebar/flex.gjs

import Component from '@glimmer/component';

import NavList from './nav-list.gjs';

export default class SidebarFlex extends Component {
  <template>
    {{!-- Static sidebar for desktop --}}
    <div class="hidden bg-gray-900 lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
      {{!-- Sidebar component, swap this element with another sidebar if you like --}}
      <div class="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-6 pb-4 dark:border-white/10 dark:bg-black/10">
        <div class="flex h-16 shrink-0 items-center">
          <img src="/img/logo-light.svg" alt="OttoApp" class="h-8 w-auto dark:hidden" />
          <img src="/img/logo-dark.svg" alt="OttoApp" class="h-8 w-auto not-dark:hidden" />
        </div>
        <nav class="flex flex-1 flex-col">
          <NavList />
        </nav>
      </div>
    </div>
  </template>
}
