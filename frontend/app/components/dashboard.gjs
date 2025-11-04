// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from "@glimmer/component";
import {service} from "@ember/service";
import {on} from '@ember/modifier';
import {action} from "@ember/object";

export default class Dashboard extends Component {
  @service session;

  <template>
    <h1>Dashboard</h1>
  </template>
}
