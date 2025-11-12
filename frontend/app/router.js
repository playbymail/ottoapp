import EmberRouter from '@embroider/router';
import config from 'frontend/config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  // public routes
  this.route('about');
  this.route('login');
  this.route('privacy');

  // admin routes (authenticated, requires "admin" role)
  this.route('admin', function () {
    this.route('dashboard');
    this.route('users', function () {
      this.route('index', { path: '/' });
      this.route('new');
      this.route('edit', { path: '/:user_id' });
    });
  });

  // users routes (authenticated, requires "user" role)
  this.route('user', function () {
    this.route('calendar'); // obsolete route to be removed in a future sprint
    this.route('dashboard');
    this.route('maps');
    this.route('my');
    this.route('profile');
    this.route('projects'); // obsolete route to be removed in a future sprint
    this.route('reports');
    this.route('secure'); // obsolete route to be removed in a future sprint
    this.route('settings');
    this.route('team'); // obsolete route to be removed in a future sprint
    this.route('teams', function () {
      this.route('heroicons'); // obsolete route to be removed in a future sprint
      this.route('tailwindlabs'); // obsolete route to be removed in a future sprint
      this.route('workcation'); // obsolete route to be removed in a future sprint
    });
  });
});
