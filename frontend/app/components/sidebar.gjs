// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import SidebarFlex from 'frontend/components/sidebar/flex';
import SidebarRelative from 'frontend/components/sidebar/relative';

export default class Sidebar extends Component {
  <template>
    {{!-- dynamic sidebar for smaller displays --}}
    <SidebarRelative />
    {{!-- static sidebar for desktops --}}
    <SidebarFlex />
  </template>
}
