// app/templates/user/dashboard.gjs

import {LinkTo} from "@ember/routing";

import DashboardRecentMapFiles from 'frontend/components/user/dashboard/recent-map-files';
import DashboardRecentTurnReportExtracts from 'frontend/components/user/dashboard/recent-turn-report-extracts';
import DashboardRecentTurnReportFiles from 'frontend/components/user/dashboard/recent-turn-report-files';

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <div class="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <header class="mb-6 flex items-center justify-between">
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          Dashboard
        </h1>

        {{!-- Space for future buttons / filters --}}
        <div class="flex items-center gap-3"></div>
      </header>

      <div class="grid gap-6 lg:grid-cols-2">
        <DashboardRecentMapFiles
          @files={{@model.recentMapFiles}}
        />

        <DashboardRecentTurnReportFiles
          @files={{@model.recentTurnReportFiles}}
        />

        <DashboardRecentTurnReportExtracts
          @files={{@model.recentTurnReportExtracts}}
        />
      </div>
    </div>
  </div>
</template>

