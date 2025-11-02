import { pageTitle } from 'ember-page-title';

<template>
  {{pageTitle "OttoMap Hero"}}
  <div class="bg-white">
    <div class="bg-white px-6 py-32 lg:px-8">
      <div class="mx-auto max-w-3xl text-base leading-7 text-gray-700">
        <h1 class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-6xl">
          Maps for your Tribenet clan
        </h1>

        <p class="mt-6 text-xl leading-8">
          OttoMap is a tool for TribeNet players to visualize their TribeNet maps.
          It reads your TribeNet turn reports and generates a Worldographer map that you can view on your
          computer.
        </p>

        <div class="mt-10 max-w-2xl">
          <figure class="mt-16">
            <img class="aspect-video rounded-xl bg-gray-50 object-cover"
                 src="/img/hero.jpg"
                 alt="">
            <figcaption class="mt-4 flex gap-x-2 text-sm leading-6 text-gray-500">
              <svg class="mt-0.5 h-5 w-5 flex-none text-gray-300" viewBox="0 0 20 20" fill="currentColor"
                   aria-hidden="true" data-slot="icon">
                <path fill-rule="evenodd"
                      d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Zm-7-4a1 1 0 1 1-2 0 1 1 0 0 1 2 0ZM9 9a.75.75 0 0 0 0 1.5h.253a.25.25 0 0 1 .244.304l-.459 2.066A1.75 1.75 0 0 0 10.747 15H11a.75.75 0 0 0 0-1.5h-.253a.25.25 0 0 1-.244-.304l.459-2.066A1.75 1.75 0 0 0 9.253 9H9Z"
                      clip-rule="evenodd"/>
              </svg>
              This Chief knows the way.
            </figcaption>
          </figure>

          <h2 class="mt-16 text-2xl font-bold tracking-tight text-gray-900">
            From tabula rasa to orbis terrarum in two shakes of a goat's tail
          </h2>
          <p class="mt-6">
            The OttoMap tool is designed to be as easy to use as possible, but it does require installing a Go
            compiler and working from the command line.
            This website was created to help you avoid all that hassle.
          </p>
          <figure class="mt-10 border-l border-indigo-600 pl-9">
            <blockquote class="font-semibold text-gray-900">
              <p>
                “OttoMap is a great tool for TribeNet players.
                It's easy to use and it generates a map that you can view on your computer.
                But this web site makes it so much nicer to use!”
              </p>
            </blockquote>
            <figcaption class="mt-6 flex gap-x-4">
              <img class="h-6 w-6 flex-none rounded-full bg-gray-50"
                   src="/img/hero-line.jpg"
                   alt="">
              <div class="text-sm leading-6">
                <strong class="font-semibold text-gray-900">Jelly Stormgoat</strong> – Chief, Clan 0987
              </div>
            </figcaption>
          </figure>
          <p class="mt-10">
            The OttoApp web server makes it easy to convert your reports to Worldographer maps.
            Just let the GM for TN3.1 know that you want to opt-in to uploading a copy of your turn report to the OttoApp server.
          </p>
        </div>

        <figure class="mt-16">
          <img class="aspect-video rounded-xl bg-gray-50 object-cover"
               src="/img/hexes-001.jpg"
               alt="">
          <figcaption class="mt-4 flex gap-x-2 text-sm leading-6 text-gray-500">
            <svg class="mt-0.5 h-5 w-5 flex-none text-gray-300" viewBox="0 0 20 20" fill="currentColor"
                 aria-hidden="true" data-slot="icon">
              <path fill-rule="evenodd"
                    d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Zm-7-4a1 1 0 1 1-2 0 1 1 0 0 1 2 0ZM9 9a.75.75 0 0 0 0 1.5h.253a.25.25 0 0 1 .244.304l-.459 2.066A1.75 1.75 0 0 0 10.747 15H11a.75.75 0 0 0 0-1.5h-.253a.25.25 0 0 1-.244-.304l.459-2.066A1.75 1.75 0 0 0 9.253 9H9Z"
                    clip-rule="evenodd"/>
            </svg>
            Explore strange new hexes.
          </figcaption>
        </figure>
      </div>
    </div>
  </div>
</template>
