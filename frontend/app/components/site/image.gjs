// Copyright (c) 2025 Michael D Henderson. All rights reserved.
import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import {LinkTo} from '@ember/routing';

export default class Image extends Component {
  <template>
    <div class="mt-32 sm:mt-40 xl:mx-auto xl:max-w-7xl xl:px-8">
      <img src="https://images.unsplash.com/photo-1529156069898-49953e39b3ac?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2832&q=80" alt="" class="aspect-5/2 w-full object-cover outline-1 -outline-offset-1 outline-black/5 xl:rounded-3xl" />
    </div>
  </template>
}
