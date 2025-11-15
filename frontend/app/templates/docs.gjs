// app/templates/docs.gjs
import {pageTitle} from 'ember-page-title';

import Layout from 'frontend/components/site/layout';
import PageSections from 'frontend/components/site/page-sections';
import SimpleSectionHeading from 'frontend/components/site/simple-section/heading';
import SimpleSectionList from 'frontend/components/site/simple-section/list';
import SimpleSectionListItem from 'frontend/components/site/simple-section/list-item';

<template>
  {{pageTitle "Documentation"}}
  <Layout>
    <PageSections>
      <SimpleSectionHeading @title="Documentation">
        OttoApp is a non-commercial hobby supported by a part-time team of one person, and documentation is on the to-do
        list.
      </SimpleSectionHeading>
      <SimpleSectionList>
        <SimpleSectionListItem @title="Parsing">
          OttoMap parses turn reports.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Rendering">
          OttoMap renders map data to Worldographer files.
        </SimpleSectionListItem>
      </SimpleSectionList>
    </PageSections>
  </Layout>
</template>
