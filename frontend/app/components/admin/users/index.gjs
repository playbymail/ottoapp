// app/components/admin/users/index.gjs

import Component from "@glimmer/component";
import { LinkTo } from "@ember/routing";

export default class AdminUsersIndex extends Component {
  <template>
    <div>
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
          <div class="sm:flex-auto">
            <h1 class="text-base font-semibold text-gray-900">Users</h1>
            <p class="mt-2 text-sm text-gray-700">A list of all the users you have access to including their name, email, timezone and role.</p>
          </div>
          <div class="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
            <LinkTo @route="admin.users.new"
                    class="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
              Add user</LinkTo>
          </div>
        </div>
      </div>
      <div class="mt-8 flow-root overflow-hidden">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <table class="w-full text-left">
            <thead class="bg-white">
            <tr>
              <th scope="col" class="relative isolate py-3.5 pr-3 text-left text-sm font-semibold text-gray-900">
                Handle
                <div class="absolute inset-y-0 right-full -z-10 w-screen border-b border-b-gray-200"></div>
                <div class="absolute inset-y-0 left-0 -z-10 w-screen border-b border-b-gray-200"></div>
              </th>
              <th scope="col" class="hidden px-3 py-3.5 text-left text-sm font-semibold text-gray-900 lg:table-cell">
                Name</th>
              <th scope="col" class="hidden px-3 py-3.5 text-left text-sm font-semibold text-gray-900 md:table-cell">
                Email</th>
              <th scope="col" class="hidden px-3 py-3.5 text-left text-sm font-semibold text-gray-900 sm:table-cell">
                Timezone</th>
              <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Roles</th>
              <th scope="col" class="py-3.5 pl-3">
                <span class="sr-only">Edit</span>
              </th>
            </tr>
            </thead>
            <tbody>
            {{#each @model as |user|}}
              <tr>
                <td class="relative py-4 pr-3 text-sm font-medium text-gray-900">
                  {{user.handle}}
                  <div class="absolute right-full bottom-0 h-px w-screen bg-gray-100"></div>
                  <div class="absolute bottom-0 left-0 h-px w-screen bg-gray-100"></div>
                </td>
                <td class="hidden px-3 py-4 text-sm text-gray-500 lg:table-cell">{{user.username}}</td>
                <td class="hidden px-3 py-4 text-sm text-gray-500 md:table-cell">{{user.email}}</td>
                <td class="hidden px-3 py-4 text-sm text-gray-500 sm:table-cell">{{user.timezone}}</td>
                <td class="px-3 py-4 text-sm text-gray-500">
                  {{#if user.roles}}
                    {{#each user.roles as |role|}}
                      <span class="inline-flex items-center rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800 mr-1">
                        {{role}}
                      </span>
                    {{/each}}
                  {{else}}
                    <span class="text-sm text-gray-500">No roles assigned</span>
                  {{/if}}
                </td>
                <td class="py-4 pl-3 text-right text-sm font-medium">
                  <LinkTo @route="admin.users.edit" @model={{user}} class="text-indigo-600 hover:text-indigo-900">
                    Edit<span class="sr-only">, {{user.handle}}</span>
                  </LinkTo>
                </td>
              </tr>
            {{/each}}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </template>
}
