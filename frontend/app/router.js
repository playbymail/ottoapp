// app/router.js
import EmberRouter from '@embroider/router';
import config from 'frontend/config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  // public routes
  this.route('about');
  this.route('docs');
  this.route('login');
  this.route('privacy');

  // admin routes (authenticated, requires "admin" role)
  this.route('admin', function () {
    this.route('dashboard', { path: '/'});
    this.route('settings', function () {
      this.route('about');
      this.route('account', { path: '/'});
      this.route('notifications');
    });
    this.route('users', function () {
      this.route('index', { path: '/' });
      this.route('new');
      this.route('edit', { path: '/:user_id' });
    });
    this.route('park');
  });

  // gm routes (authenticated, requires "gm" role)
  this.route('gm', function () {
    this.route('dashboard', { path: '/' });
    this.route('turn-report-files', function () {
      this.route('upload');  // /gm/turn-report-files/upload
      // later: index, history, errors, etc.
    });

    // future GM-only routes:
    // this.route('maps');
    // this.route('ally-data');
  });

  // users routes (authenticated, requires "user" role)
  this.route('user', function () {
    this.route('dashboard', { path: '/' });
    this.route('documents', function () {
      this.route('show', { path: '/:document_id' });
    });
    this.route('extracts', function () {
      this.route('index', {path: '/'});
      this.route('show', { path: '/:document_id' });
    });
    this.route('maps');
    this.route('reports');
    this.route('settings', function () {
      this.route('about');
      this.route('account', { path: '/'});
      this.route('maps');
      this.route('notifications');
      this.route('teams');
    });
    this.route('team'); // obsolete route to be removed in a future sprint
    this.route('teams', function () {
      this.route('heroicons'); // obsolete route to be removed in a future sprint
      this.route('tailwindlabs'); // obsolete route to be removed in a future sprint
      this.route('workcation'); // obsolete route to be removed in a future sprint
    });
  });
});
