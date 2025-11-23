// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// app/controllers/user/dashboard.js

import Controller from '@ember/controller';
import {  service } from '@ember/service';

export default class UserDashboardController extends Controller {
  @service store;

  constructor(...args) {
    super(...args);
    //console.log('app/controllers/user/dashboard');
  }
}
