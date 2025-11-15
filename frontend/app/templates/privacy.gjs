// app/templates/privacy.gjs
import {pageTitle} from 'ember-page-title';

import Layout from 'frontend/components/site/layout';
import PageSections from 'frontend/components/site/page-sections';
import SimpleSectionHeading from 'frontend/components/site/simple-section/heading';
import SimpleSectionList from 'frontend/components/site/simple-section/list';
import SimpleSectionListItem from 'frontend/components/site/simple-section/list-item';

<template>
  {{pageTitle "Privacy"}}
  <Layout>
    <PageSections>
      <SimpleSectionHeading @title="Privacy">
        OttoApp is a non-commercial hobby supported by a part-time team of one person, but it tries to be nice with your data.
      </SimpleSectionHeading>
      <SimpleSectionList>
        <SimpleSectionListItem @title="Email">
          OttoApp stores your e-mail address.
          It's used to control access to the turn reports and map data that have been uploaded.
          If you don't consent to this, please post a request on the OttoMap Discord server and we'll delete your account and data.
        </SimpleSectionListItem>
        <SimpleSectionListItem @title="Logging">
          The application is hosted on a VPS that purges web traffic logs after thirty days.
        </SimpleSectionListItem>
      </SimpleSectionList>
    </PageSections>
  </Layout>
</template>
