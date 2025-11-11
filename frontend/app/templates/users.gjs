import RouteTemplate from 'ember-route-template';

export default RouteTemplate(
  <template>
    <div class="min-h-screen bg-gray-100">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {{outlet}}
      </div>
    </div>
  </template>
);
