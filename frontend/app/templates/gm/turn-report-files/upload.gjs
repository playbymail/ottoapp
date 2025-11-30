// app/templates/gm/turn-report-files/upload.gjs

import TurnReportUpload from 'frontend/components/gm/turn-report-files/upload';

<template>
  <section class="max-w-xl mx-auto py-8">
    <h1 class="text-2xl font-semibold mb-4">
      Upload Turn Report
    </h1>

    {{#if this.errorMessage}}
      <div class="mb-4 rounded-md border border-red-300 bg-red-50 px-4 py-2 text-sm text-red-800">
        {{this.errorMessage}}
      </div>
    {{/if}}

    {{#if this.successMessage}}
      <div class="mb-4 rounded-md border border-green-300 bg-green-50 px-4 py-2 text-sm text-green-800">
        {{this.successMessage}}
      </div>
    {{/if}}

    <TurnReportUpload
      @isUploading={{this.isUploading}}
      @onFileSelected={{@controller.handleFileSelected}}
    />
  </section>
</template>
