// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from "@glimmer/component";
import {tracked} from "@glimmer/tracking";
import {action} from "@ember/object";
import {on} from "@ember/modifier";
import {fn} from "@ember/helper";

export default class TimezonePicker extends Component {
  @tracked showingAll = false;
  @tracked allTimezones = null;
  @tracked searchTerm = "";

  get filteredTimezones() {
    if (!this.allTimezones) return [];
    let q = this.searchTerm.trim().toLowerCase();
    if (!q) return this.allTimezones;
    return this.allTimezones.filter((tz) =>
      tz.label.toLowerCase().includes(q)
    );
  }

  isSelected = (label) => {
    return label === this.args.value;
  }

  @action async showAll() {
    this.showingAll = true;
    if (!this.allTimezones) {
      let res = await fetch("/api/timezones");
      this.allTimezones = await res.json();
    }
  }

  @action closeAll() {
    this.showingAll = false;
    this.searchTerm = "";
  }

  @action updateSearch(e) {
    this.searchTerm = e.target.value;
  }

  @action selectTimezone(label) {
    this.args.onChange?.(label);
    this.closeAll();
  }

  <template>
    <div class="space-y-3">
      <div class="block w-full rounded-md bg-gray-50 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 sm:max-w-xs sm:text-sm/6">
        {{@value}}
      </div>

      <button
        type="button"
        class="text-sm text-indigo-600 hover:text-indigo-500"
        {{on "click" this.showAll}}
      >
        Update timezone
      </button>
    </div>

    {{#if this.showingAll}}
      <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30">
        <div class="w-full max-w-lg rounded-lg bg-white p-4 shadow-lg">
          <div class="flex items-center justify-between gap-3 mb-3">
            <h2 class="text-base font-semibold">All timezones</h2>
            <button type="button" class="text-gray-400 hover:text-gray-500" {{on "click" this.closeAll}}>
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <input
            type="text"
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm mb-3 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
            placeholder="Search by name, region..."
            value={{this.searchTerm}}
            {{on "input" this.updateSearch}}
          />

          <div class="max-h-80 overflow-y-auto divide-y divide-gray-100 border border-gray-200 rounded-md">
            {{#each this.filteredTimezones as |tz|}}
              <button
                type="button"
                class="flex w-full items-center justify-between px-3 py-2 text-left hover:bg-gray-50"
                {{on "click" (fn this.selectTimezone tz.label)}}
              >
                <span class="text-sm">{{tz.label}}</span>
                {{#if (this.isSelected tz.label)}}
                  <span class="text-indigo-600 text-xs font-semibold">Selected</span>
                {{/if}}
              </button>
            {{/each}}
          </div>
        </div>
      </div>
    {{/if}}
  </template>
}
