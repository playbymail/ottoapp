// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Admin shell - different from user shell to distinguish admin interface
// No sidebar for now, just a simple header

<template>
  <div class="min-h-screen bg-gray-100">
    <header class="bg-white shadow">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between">
          <h1 class="text-3xl font-bold tracking-tight text-gray-900">Admin Dashboard</h1>
          <span class="inline-flex items-center rounded-md bg-red-50 px-2 py-1 text-xs font-medium text-red-700 ring-1 ring-inset ring-red-600/10">
            Admin
          </span>
        </div>
      </div>
    </header>

    <main class="py-10">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {{outlet}}
      </div>
    </main>
  </div>
</template>
