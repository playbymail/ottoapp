// app/controllers/user/extracts.js

import Controller from '@ember/controller';
import {  service } from '@ember/service';

export default class UserExtractsController extends Controller {
  @service store;

  constructor(...args) {
    super(...args);
    // console.log('app/controllers/user/extracts');
  }
}
