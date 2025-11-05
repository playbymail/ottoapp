import EmberRouter from '@embroider/router';
import config from 'frontend/config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('about');
  this.route('calendar');
  this.route('dashboard');
  this.route('login');
  this.route('maps');
  this.route('my');
  this.route('profile');
  this.route('projects');
  this.route('reports');
  this.route('secure');
  this.route('settings');
  this.route('team');
  this.route('teams');
  this.route('teams/heroicons');
  this.route('teams/tailwindlabs');
  this.route('teams/workcation');
});
