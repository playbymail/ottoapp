// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from "@glimmer/component";
import {service} from "@ember/service";

export default class Profile extends Component {
  @service session;

  get currentUser() {
    console.log('com/prof', this.session.data.authenticated.user)
    if (!this.session.isAuthenticated) {
      return {
        id: 0,
        username: 'Guest',
        roles: ['guest'],
      };
    }
    return this.session.data.authenticated.user;
  }

  <template>
    <h1>Profile</h1>
    {{#if this.session.isAuthenticated}}
      <p>
        Signed in as {{this.currentUser.username}}
      </p>
    {{else}}
      <p>
        Guest
      </p>
    {{/if}}
  </template>
}
