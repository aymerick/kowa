import Ember from 'ember';

// All that mess simply because of that issue:
//  https://github.com/emberjs/data/issues/1308
var ContentEditionController = Ember.Mixin.create({
  editionRelationships: Ember.A([ ]),
  editionDefaultTitle: null,
  editionSaveMsgOk: null,
  editionSaveMsgErr: null,

  setupEdition: function(fields) {
    if (!Ember.isNone(fields)) {
      this.set('editionRelationships', fields);
    }

    var model = this.get('model');
    var self = this;

    this.get('editionRelationships').forEach(function (field) {
      // eg: cover => previousCover
      self.set(self.previousFieldName(field), model.get(field));
    });

    this.defineIsDirtyProperty();
    this.defineNothingChangedProperty();
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
      return self.previousField(field).get('id') !== model.get(field).get('id');
    });
  },

  // define isDirty property
  defineIsDirtyProperty: function() {
    // check if property was already defined
    if (this.hasOwnProperty('isDirty')) {
      return;
    }

    var self = this;

    var watchProperties = Ember.A([ 'model.isDirty', 'model.isNew' ]);

    this.get('editionRelationships').forEach(function(field) {
      watchProperties.pushObject('model.' + field);
      watchProperties.pushObject(self.previousFieldName(field));
    });

    var prop = Ember.computed(function() {
      var model = self.get('model');

      return (model.get('isDirty') || model.get('isNew') || self.editionFieldChanged());
    });

    // javackish way of calling a method with a variable number of parameters. It would be easier that way:
    //    prop.property(watchProperties)
    prop = Ember.apply(prop, Ember.ComputedProperty.prototype.property, watchProperties);

    Ember.defineProperty(this, 'isDirty', prop);
  },

  // define nothingChanged property
  defineNothingChangedProperty: function() {
    // check if property was already defined
    if (this.hasOwnProperty('nothingChanged')) {
      return;
    }

    // nothingChanged property
    Ember.defineProperty(this, 'nothingChanged', Ember.computed.not('isDirty'));
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

export default ContentEditionController;
