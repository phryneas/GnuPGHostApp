module.exports = function (grunt) {
    grunt.initConfig({
        browserify: {
            dist: {
                files: {
                    "./src/background-compiled.js": ["./src/background.js"],
                }
            },
        },
        watch: {
            scripts: {
                files: [".src/background.js"],
                tasks: ["browserify:dist"]
            },
        },
    });

    grunt.loadNpmTasks("grunt-browserify");
    grunt.loadNpmTasks("grunt-contrib-watch");

    grunt.registerTask("default", ["watch"]);
    grunt.registerTask("build", ["browserify:dist"]);
};