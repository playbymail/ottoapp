// app/components/gm/sidebar/nav-list.gjs

import Component from '@glimmer/component';
import {service} from '@ember/service';
import not from 'frontend/helpers/not';

import {LinkTo} from '@ember/routing';

export default class NavList extends Component {
  @service router;
  @service session;

  activeLinkStyle = 'bg-gray-50 p-2 text-sm/6 font-semibold text-indigo-600 dark:bg-white/5 dark:text-white';
  inactiveLinkStyle = 'p-2 text-sm/6 font-semibold text-gray-700 hover:bg-gray-50 hover:text-indigo-600 dark:text-gray-400 dark:hover:bg-white/5 dark:hover:text-white';

  linkClass = (routeName, pfxLinkStyle) => {
    const current = this.router.currentRouteName;
    const isActive = current === routeName || current?.startsWith(`${routeName}.`);
    return `${pfxLinkStyle} ${isActive ? this.activeLinkStyle : this.inactiveLinkStyle}`;
  }

  activeIconStyle = 'text-indigo-600 dark:text-white';
  inactiveIconStyle = 'text-gray-400 group-hover:text-indigo-600 dark:group-hover:text-white';

  iconClass = (routeName, pfxLinkStyle) => {
    const current = this.router.currentRouteName;
    const isActive = current === routeName || current?.startsWith(`${routeName}.`);
    return `${pfxLinkStyle} ${isActive ? this.activeIconStyle : this.inactiveIconStyle}`;
  }

  <template>
    <ul role="list" class="flex flex-1 flex-col gap-y-7">
      <li>
        <ul role="list" class="-mx-2 space-y-1">
          <li>
            {{!-- Current: "bg-gray-50 dark:bg-white/5 text-indigo-600 dark:text-white", Default: "text-gray-700 dark:text-gray-400 hover:text-indigo-600 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-white/5" --}}
            <LinkTo @route="user.dashboard" class={{this.linkClass "user.dashboard" "group flex gap-x-3 rounded-md"}}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
                   aria-hidden="true" class={{this.iconClass "gm.dashboard" "size-6 shrink-0"}}>
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M8.25 21v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21m0 0h4.5V3.545M12.75 21h7.5V10.75M2.25 21h1.5m18 0h-18M2.25 9l4.5-1.636M18.75 3l-1.5.545m0 6.205 3 1m1.5.5-1.5-.5M6.75 7.364V3h-3v18m3-13.636 10.5-3.819" />
              </svg>
              Dashboard
            </LinkTo>
          </li>
        </ul>
      </li>

      <li class="mt-auto">
        {{#if this.session.canAccessAdminRoutes}}
          <LinkTo @route="admin.dashboard" class={{this.linkClass "user.settings"
                                                                  "group -mx-2 flex gap-x-3 rounded-md"}}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
                 aria-hidden="true" class={{this.iconClass "user.settings" "size-6 shrink-0"}}>
              <path stroke-linecap="round" stroke-linejoin="round"
                    d="M8.25 18.75a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 0 1-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 0 0-3.213-9.193 2.056 2.056 0 0 0-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 0 0-10.026 0 1.106 1.106 0 0 0-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12" />
            </svg>
            Admin Dashboard
          </LinkTo>
        {{/if}}
        {{#if this.session.canAccessUserRoutes}}
          <LinkTo @route="user.dashboard" class={{this.linkClass "user.dashboard"
                                                                 "group -mx-2 flex gap-x-3 rounded-md"}}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
                 aria-hidden="true" class={{this.iconClass "user.dashboard" "size-6 shrink-0"}}>
              <path
                d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
                stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            User Dashboard
          </LinkTo>
        {{/if}}
        <LinkTo @route="user.settings" class={{this.linkClass "user.settings" "group -mx-2 flex gap-x-3 rounded-md"}}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
               aria-hidden="true" class={{this.iconClass "user.settings" "size-6 shrink-0"}}>
            <path
              d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
              stroke-linecap="round" stroke-linejoin="round" />
            <path d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          Settings
        </LinkTo>
      </li>
    </ul>
  </template>
}
