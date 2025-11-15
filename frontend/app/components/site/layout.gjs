// app/components/site/layout.gjs
import Component from '@glimmer/component';

// Requires a TailwindCSS Plus license.
// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header

import {LinkTo} from '@ember/routing';
import Header from 'frontend/components/site/header';
import Hero from 'frontend/components/site/hero';
import Content from 'frontend/components/site/content';
import Image from 'frontend/components/site/image';
import Feature from 'frontend/components/site/feature';
import LogoCloud from 'frontend/components/site/logo-cloud';
import Team from 'frontend/components/site/team';
import Blog from 'frontend/components/site/blog';
import Footer from 'frontend/components/site/footer';

export default class SiteShell extends Component {
  <template>
    <div class="bg-white">
      <Header />
      <main class="isolate">
        {{yield}}
        {{!-- Hero / --}}
        {{!-- Content / --}}
        {{!-- Image / --}}
        {{!-- Feature / --}}
        {{!-- LogoCloud / --}}
        {{!-- Team / --}}
        {{!-- Blog / --}}
      </main>
      <Footer />
    </div>
  </template>
}
