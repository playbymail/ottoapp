// app/serializers/application.js
import JSONAPISerializer from '@ember-data/serializer/json-api';

export default class ApplicationSerializer extends JSONAPISerializer {
  // Default JSONAPISerializer should handle dirty tracking automatically
  // Keeping this file for future customization if needed
}
