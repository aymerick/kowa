#global module:false

"use strict"

module.exports = (grunt) ->
  grunt.loadNpmTasks "grunt-contrib-copy"
  grunt.loadNpmTasks "grunt-contrib-sass"
  grunt.loadNpmTasks "grunt-contrib-watch"
  grunt.loadNpmTasks "grunt-exec"

  grunt.initConfig
    sass:
      dist:
        options:
          bundleExec: true
        files: [{
          expand: true
          cwd: "sass/"
          src: ["*.scss"]
          dest: "static/css/"
          ext: ".css"
        }]

    copy:
      jquery:
        files: [{
          expand: true
          cwd: "bower_components/jquery/dist/"
          src: "jquery.min.js"
          dest: "static/js/"
        }]
      bootstrap:
        files: [{
          expand: true
          cwd: "bower_components/bootstrap/dist/css/"
          src: "bootstrap.min.css"
          dest: "static/css/"
        },
        {
          expand: true
          cwd: "bower_components/bootstrap/dist/js/"
          src: "bootstrap.min.js"
          dest: "static/js/"
        }]
      fontawesome:
        files: [{
          expand: true
          cwd: "bower_components/fontawesome/css/"
          src: "font-awesome.min.css"
          dest: "static/css/"
        },
        {
          expand: true
          cwd: "bower_components/fontawesome/fonts/"
          src: "*"
          dest: "static/fonts/"
        }]

    watch:
      source:
        files: [
          "sass/**/*"
        ]
        tasks: [
          "sass"
        ]

  grunt.registerTask "build", [
    "sass"
    "copy"
  ]

  grunt.registerTask "dev", [
    "build"
    "watch"
  ]

  grunt.registerTask "default", [
    "dev"
  ]
