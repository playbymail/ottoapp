// app/serializers/document.js
import ApplicationSerializer from './application';

export default class DocumentSerializer extends ApplicationSerializer {
  normalize(modelClass, resourceHash) {
    if (resourceHash.links) {
      resourceHash.attributes = resourceHash.attributes || {};
      resourceHash.attributes.links = resourceHash.links;
    }
    return super.normalize(modelClass, resourceHash);
  }
}
