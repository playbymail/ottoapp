// Copyright (c) 2025 Michael D Henderson. All rights reserved.
import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

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

// <!-- Include this script tag or install `@tailwindplus/elements` via npm: -->
// <!-- <script src="https://cdn.jsdelivr.net/npm/@tailwindplus/elements@1" type="module"></script> -->

export default class SiteShell extends Component {
  <template>
    <div class="bg-white">
      <Header />
      <main class="isolate">
        {{outlet}}
        <!-- Hero / -->
        <!-- Content / -->
        <!-- Image / -->
        <!-- Feature / -->
        <!-- LogoCloud / -->
        <!-- Team / -->
        <!-- Blog / -->
      </main>
      <Footer />
    </div>
  </template>
}
