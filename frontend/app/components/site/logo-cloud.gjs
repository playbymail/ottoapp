// Copyright (c) 2025 Michael D Henderson. All rights reserved.
import Component from '@glimmer/component';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import {LinkTo} from '@ember/routing';

export default class LogoCloud extends Component {
  <template>
    <div class="relative isolate -z-10 mt-32 sm:mt-48">
      <div class="absolute inset-x-0 top-1/2 -z-10 flex -translate-y-1/2 justify-center overflow-hidden mask-[radial-gradient(50%_45%_at_50%_55%,white,transparent)]">
        <svg aria-hidden="true" class="h-160 w-7xl flex-none stroke-gray-200">
          <defs>
            <pattern id="e9033f3e-f665-41a6-84ef-756f6778e6fe" width="200" height="200" x="50%" y="50%" patternUnits="userSpaceOnUse" patternTransform="translate(-100 0)">
              <path d="M.5 200V.5H200" fill="none" />
            </pattern>
          </defs>
          <svg x="50%" y="50%" class="overflow-visible fill-gray-50">
            <path d="M-300 0h201v201h-201Z M300 200h201v201h-201Z" stroke-width="0" />
          </svg>
          <rect width="100%" height="100%" fill="url(#e9033f3e-f665-41a6-84ef-756f6778e6fe)" stroke-width="0" />
        </svg>
      </div>
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <h2 class="text-center text-lg/8 font-semibold text-gray-900">Trusted by the world’s most innovative teams</h2>
        <div class="mx-auto mt-10 grid max-w-lg grid-cols-4 items-center gap-x-8 gap-y-10 sm:max-w-xl sm:grid-cols-6 sm:gap-x-10 lg:mx-0 lg:max-w-none lg:grid-cols-5">
          <img width="158" height="48" src="https://tailwindcss.com/plus-assets/img/logos/158x48/transistor-logo-gray-900.svg" alt="Transistor" class="col-span-2 max-h-12 w-full object-contain lg:col-span-1" />
          <img width="158" height="48" src="https://tailwindcss.com/plus-assets/img/logos/158x48/reform-logo-gray-900.svg" alt="Reform" class="col-span-2 max-h-12 w-full object-contain lg:col-span-1" />
          <img width="158" height="48" src="https://tailwindcss.com/plus-assets/img/logos/158x48/tuple-logo-gray-900.svg" alt="Tuple" class="col-span-2 max-h-12 w-full object-contain lg:col-span-1" />
          <img width="158" height="48" src="https://tailwindcss.com/plus-assets/img/logos/158x48/savvycal-logo-gray-900.svg" alt="SavvyCal" class="col-span-2 max-h-12 w-full object-contain sm:col-start-2 lg:col-span-1" />
          <img width="158" height="48" src="https://tailwindcss.com/plus-assets/img/logos/158x48/statamic-logo-gray-900.svg" alt="Statamic" class="col-span-2 col-start-2 max-h-12 w-full object-contain sm:col-start-auto lg:col-span-1" />
        </div>
      </div>
    </div>
  </template>
}
