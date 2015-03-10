#global module:false

"use strict"

module.exports = (grunt) ->
  grunt.loadNpmTasks "grunt-contrib-copy"

  grunt.initConfig
    copy:
      fontawesome:
        files: [{
          expand: true
          cwd: "bower_components/fontawesome/css/"
          src: [ "font-awesome.min.css", "font-awesome.css.map" ]
          dest: "assets/vendor/css/"
        },
        {
          expand: true
          cwd: "bower_components/fontawesome/fonts/"
          src: "*"
          dest: "assets/vendor/fonts/"
        }]
      jquery:
        files: [{
          expand: true
          cwd: "bower_components/jquery/dist/"
          src: [ "jquery.min.js", "jquery.min.map" ]
          dest: "assets/vendor/js"
        }]
      bootstrap:
        files: [{
          expand: true
          cwd: "bower_components/bootstrap/dist/css/"
          src: [ "bootstrap.min.css", "bootstrap.css.map" ]
          dest: "assets/vendor/css"
        },
        {
          expand: true
          cwd: "bower_components/bootstrap/dist/js/"
          src: [ "bootstrap.min.js", "bootstrap.min.js" ]
          dest: "assets/vendor/js"
        }]

  grunt.registerTask "build", [
    "copy"
  ]

  grunt.registerTask "default", [
    "build"
  ]
