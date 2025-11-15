// app/components/layouts/home-screens/sidebar/nav-item.gjs
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

import Component from '@glimmer/component';
import {service} from '@ember/service';
import {LinkTo} from '@ember/routing';

export default class SidbarNavItem extends Component {
  @service router;

  /**
   * Fuzzy route match:
   * - exact match: "admin.users"
   * - or any child: "admin.users.*" (index, new, edit, etc)
   */
  isActiveRoute = (baseRoute) => {
    const current = this.router.currentRouteName;
    if (!current || !baseRoute) return false;
    return (
      current === baseRoute ||
      current.startsWith(`${baseRoute}.`)
    );
  }

  iconClassFor = (routeName) => {
    return this.isActiveRoute(routeName)
      ? 'size-6 shrink-0 text-indigo-600'
      : 'size-6 shrink-0 text-gray-400 group-hover:text-indigo-600';
  }

  linkClassFor = (routeName) => {
    return this.isActiveRoute(routeName)
      ? 'group flex gap-x-3 rounded-md bg-gray-100 p-2 text-sm/6 font-semibold text-indigo-600'
      : 'group flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold text-gray-700 hover:bg-gray-100 hover:text-indigo-600';
  }

// Current: "bg-gray-100 text-indigo-600"
// Default: "text-gray-700 hover:text-indigo-600 hover:bg-gray-100"

  <template>
    <li>
      <LinkTo @route="{{this.args.link}}"
              class={{this.linkClassFor this.args.link}}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" data-slot="icon"
             aria-hidden="true" class={{this.iconClassFor this.args.link}}>
          {{yield}}
        </svg>
        {{this.args.text}}
      </LinkTo>
    </li>
  </template>
}
