// app/templates/about.gjs
import {pageTitle} from 'ember-page-title';

import Layout from 'frontend/components/site/layout';
import PageSections from 'frontend/components/site/page-sections';
import SimpleSectionHeading from 'frontend/components/site/simple-section/heading';
import SimpleSectionList from 'frontend/components/site/simple-section/list';
import SimpleSectionListItem from 'frontend/components/site/simple-section/list-item';

<template>
  {{pageTitle "About OttoMap"}}
  <Layout>
    <PageSections>
      <SimpleSectionHeading @title="About OttoApp">
        OttoApp is the name of the repository on Github.
        It is an open source hobby project is not intended for commercial use.
      </SimpleSectionHeading>
      <SimpleSectionList>
        <SimpleSectionListItem @title="Disclaimers">
          OttoApp is not associated with, affiliated with, or approved by TribeNet, Worldographer, Tailwind, or Caddy.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="OttoMap">
          OttoMap is written in Go and uses the Pigeon PEG parser generator to read the TribeNet report files.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="OttoWeb">
          The front end is written in EmberJS.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="API server">
          The API server is written in Go.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Worldographer">
          OttoMap renders to Worldographer because it is a great tool, has a free version, and the author has documented the data structures quite nicely.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Tailwind Plus UI Blocks">
          OttoApp uses Tailwind Plus UI Blocks for the parts of this site that look nice and work well on multiple browsers and devices.
          Tailwind Plus is not open source and you are not allowed to use the styles from this site in other projects unless you purchase a license from Tailwind.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Caddy">
          This site uses Caddy as the web server.
          It is nice and it works.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Powerful AI">
          Amp and ChatGPT were used to build, test, and document the project.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Database backups">
          Ha ha ha ha.
          Nope.
        </SimpleSectionListItem>
      </SimpleSectionList>
    </PageSections>
  </Layout>
</template>
