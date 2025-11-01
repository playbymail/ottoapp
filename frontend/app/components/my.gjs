// app/components/my.js

import Component from '@ember/component';
import { inject as service } from '@ember/service';

export default class MyComponent extends Component {
  @service session;
  @service currentUser;

  <template>
    <nav>
      {{#link-to 'index' classNames='navbar-brand'}}
        Home
      {{/link-to}}

      {{#if this.session.isAuthenticated}}
        <button {{on "click" this.logout}}>Logout</button>
        {{#if this.currentUser.user}}
          <p>Signed in as {{this.currentUser.user.name}}</p>
        {{/if}}
      {{else}}
        <button {{on "click" this.login}}>Login</button>
      {{/if}}
    </nav>
  </template>
}

