// app/templates/index.gjs
import {pageTitle} from 'ember-page-title';

import Layout from 'frontend/components/site/layout';

<template>
  {{pageTitle "OttoMap"}}
  <Layout>
    <div class="mt-32 sm:mt-40 xl:mx-auto xl:max-w-7xl xl:px-8">
      <img src="/img/hero.jpg" />
    </div>
  </Layout>
</template>
