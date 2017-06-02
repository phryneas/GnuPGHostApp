module.exports = function (grunt) {
    grunt.initConfig({
        browserify: {
            dist: {
                options: {
                    browserifyOptions: {
                        standalone: 'dist'
                    },
                },
                files: {
                    "./dist/module.js": ["./src/NativeOpenGpgMeClient.js"],
                }
            },
            specs: {
                files: {
                    "./dist/specs.js": ["specs/**/*.spec.js"],
                }
            },
            options: {
                configure: function (bundler) {
                    bundler.plugin(require('tsify'));
                    bundler.transform(require('babelify'), {
                        presets: ['es2015'],
                        extensions: ['.ts', '.js']
                    });
                }
            }
        },
        watch: {
            scripts: {
                files: ["./src/*.js"],
                tasks: ["browserify:dist", "browserify:specs"]
            },
            specs: {
                files: ["./specs/*.js"],
                tasks: ["browserify:specs"]
            },
        },
        jasmine: {
            pivotal: {
                src: 'src/**/*.js',
                options: {
                    specs: './dist/specs.js',
                }
            }
        }
    });

    grunt.loadNpmTasks("grunt-browserify");
    grunt.loadNpmTasks("grunt-contrib-watch");
    grunt.loadNpmTasks('grunt-contrib-jasmine');

    grunt.registerTask("default", ["watch"]);
    grunt.registerTask("build", ["browserify:dist"]);
    grunt.registerTask("test", ["browserify:specs", "jasmine"]);
};