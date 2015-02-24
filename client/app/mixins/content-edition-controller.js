import Ember from 'ember';

// All that mess simply because of that issue:
//  https://github.com/emberjs/data/issues/1308
var ContentEditionControllerMixin = Ember.Mixin.create({
  isDirty: false,
  nothingChanged: Ember.computed.not('isDirty'),

  watchProperties: Ember.A([ ]),

  // properties setup in controller
  editionRelationships: Ember.A([ ]),
  editionDefaultTitle: null,
  editionSaveMsgOk: null,
  editionSaveMsgErr: null,

  setupEdition: function(fields) {
    this.resetEdition();

    if (!Ember.isNone(fields)) {
      this.set('editionRelationships', fields);
    }

    var model = this.get('model');
    var self = this;

    this.get('editionRelationships').forEach(function (field) {
      // eg: cover => previousCover
      self.set(self.previousFieldName(field), model.get(field));
    });

    this.setWatchProperties();
    this.addDirtinessObserving();
    this.synchronizeIsDirtyProperty();
  },

  resetEdition: function() {
    this.removeDirtinessObserving();

    this.set('watchProperties', Ember.A([ ]));
    this.set('isDirty', false);
  },

  commitEdition: function(okMsg, errorMsg) {
    if (!this.get('isDirty')) {
      // This should never happen
      return;
    }

    // persist on server
    var model = this.get('model');
    var self = this;

    model.save().then(function (recordSaved) {
      self.get('flashes').success(okMsg);

      self.get('editionRelationships').forEach(function (field) {
        // eg: self.set('previousCover', model.get('cover'));
        self.set(self.previousFieldName(field), model.get(field));
      });

      return recordSaved;
    }).catch(function () {
      self.get('flashes').danger(errorMsg);
    });
  },

  rollbackEdition: function() {
    var model = this.get('model');

    if (model.get('isNew')) {
      model.deleteRecord();
    } else {
      var self = this;
      this.get('editionRelationships').forEach(function(field) {
        // eg: model.set('cover', self.get('previousCover'))
        model.set(field, self.previousField(field));
      });

      model.rollback();
    }
  },

  previousField: function(fieldName) {
    return this.get(this.previousFieldName(fieldName));
  },

  previousFieldName: function(fieldName) {
    return ('previous' + fieldName.classify());
  },

  editionFieldChanged: function() {
    var model = this.get('model');
    var self  = this;

    return this.get('editionRelationships').any(function(field) {
      var previousValue = self.previousField(field);
      var value = model.get(field);

      return ((Ember.isNone(previousValue) && !Ember.isNone(value)) ||
              (!Ember.isNone(previousValue) && Ember.isNone(value)) ||
              (!Ember.isNone(previousValue) && !Ember.isNone(value) && (previousValue.get('id') !== value.get('id'))));
    });
  },

  synchronizeIsDirtyProperty: function() {
    var model = this.get('model');
    var currentDirtinessValue = model.get('isDirty') || model.get('isNew') || this.editionFieldChanged();

    this.set('isDirty', currentDirtinessValue);
  },

  setWatchProperties: function() {
    var self = this;
    var watchProperties = Ember.A([ 'model.isDirty', 'model.isNew' ]);

    this.get('editionRelationships').forEach(function(field) {
      watchProperties.pushObject('model.' + field);
      watchProperties.pushObject(self.previousFieldName(field));
    });

    this.set('watchProperties', watchProperties);
  },

  addDirtinessObserving: function() {
    var self = this;

    this.get('watchProperties').forEach(function(watchProp) {
      self.addObserver(watchProp, self, 'synchronizeIsDirtyProperty');
    });
  },

  removeDirtinessObserving: function() {
    var self = this;

    this.get('watchProperties').forEach(function(watchProp) {
      self.removeObserver(watchProp, self, 'synchronizeIsDirtyProperty');
    });
  },

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    removeCover: function() {
      this.get('model').set('cover', null);
    },

    saveContent: function() {
      if (this.get('editionDefaultTitle') !== null) {
        // set a default title
        var model = this.get('model');
        if (!model.get('title')) {
            model.set('title', '(Untitled)');
        }
      }

      // save
      this.commitEdition(this.get('editionSaveMsgOk'), this.get('editionSaveMsgErr'));
    }
  }
});

export default ContentEditionControllerMixin;
