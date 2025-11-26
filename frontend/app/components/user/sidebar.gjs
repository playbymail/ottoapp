// app/components/user/sidebar.gjs

import Component from '@glimmer/component';

import SidebarFlex from './sidebar/flex';
import SidebarRelative from './sidebar/relative';

export default class Sidebar extends Component {
  <template>
    {{!-- dynamic sidebar for smaller displays --}}
    <SidebarRelative />
    {{!-- static sidebar for desktops --}}
    <SidebarFlex />
  </template>
}
