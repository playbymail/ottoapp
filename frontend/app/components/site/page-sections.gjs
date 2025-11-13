// app/components/site/page-sections.gjs
// Requires a TailwindCSS Plus license.
// https://tailwindcss.com/plus/ui-blocks/marketing/sections/feature-sections#simple

import Component from '@glimmer/component';
import SimpleSectionHeading from './simple-section/heading';
import SimpleSectionList from './simple-section/list';

export default class PageSections extends Component {
  <template>
    <div class="bg-white py-24 sm:py-32">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        {{yield}}
      </div>
    </div>
  </template>
}
