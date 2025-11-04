// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from "@glimmer/component";
import {service} from "@ember/service";
import {on} from '@ember/modifier';
import {action} from "@ember/object";

export default class LogoutForm extends Component {
  @service session;
  @service router; // so we can redirect after logout

  @action async logout() {
    // in dev, we can hit timing issues with the response, so ignore errors and assume the library did its job.
    await this.session.invalidate().catch(() => {
    }); // ← this calls app/authenticators/server.js → invalidate()
    this.router.transitionTo("/");  // force the browser back to the login page
  }

  <template>
    <button type="button" {{on "click" this.logout}}>
      Logout
    </button>
  </template>
}

/*
          <button type="button" {{on "click" @onLogout}}>
            Logout
          </button>

 */
