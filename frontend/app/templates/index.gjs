import { pageTitle } from 'ember-page-title';

import SiteShell from 'frontend/components/site/shell';
import Header from 'frontend/components/site/header';
import Footer from 'frontend/components/site/footer';
import Image from 'frontend/components/site/image';

<template>
  {{pageTitle "Tailwind Home Screen"}}
  <div class="bg-white">
    <Header />
    <main class="isolate">
      <div class="mt-32 sm:mt-40 xl:mx-auto xl:max-w-7xl xl:px-8">
        <img src="/img/hero.jpg" />
      </div>
    </main>
    <Footer />
  </div>
</template>
