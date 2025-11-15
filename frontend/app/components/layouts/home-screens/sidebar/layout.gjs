// app/components/layouts/home-screens/sidebar/layout.js
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar
//
// This example requires updating your template:
//   <html class="h-full bg-white">
//   <body class="h-full">

import Component from '@ember/component';
import {service} from '@ember/service';

import Sidebar from './sidebar';
import StickySearchHeader from './sticky-search-header';
import Deployments from './deployments';
import ActivityFeed from './activity-feed';

export default class Layout extends Component {
  <template>
    <Sidebar @logoLight="/img/logo-light.svg" @logoDark="/img/logo-dark.svg" />

    <div class="xl:pl-72">
      <StickySearchHeader />

      {{yield}}
    </div>
  </template>
}
