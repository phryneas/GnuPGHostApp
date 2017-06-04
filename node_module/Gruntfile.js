module.exports = function (grunt) {
    grunt.initConfig({
        ts: {
            dist: {
                src: ["src/native-opengpgme-client.ts", "!node_modules/**"],
                outDir: "dist/ts",
                options: {
                    "rootDir": "./src",
                    "noImplicitAny": true,
                    "target": "es5",
                    declaration: true
                }
            }
        },
        dts_bundle: {
            dist: {
                options: {
                    name: 'native-opengpgme-client',
                    main: './dist/ts/native-opengpgme-client.d.ts',
                    out: '../module.d.ts'
                }
            }
        },
        browserify: {
            dist: {
                files: {
                    "./dist/module.js": ["dist/ts/*.js"],
                },
                options: {
                    browserifyOptions: {
                        standalone: 'native-opengpgme-client'
                    },
                    plugin: ['browserify-derequire']
                }
            },
            specs: {
                files: {
                    "./dist/specs.js": ["specs/**/*.spec.js"],
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
        },
        watch: {
            scripts: {
                files: ["./src/*.ts"],
                tasks: ["build", "browserify:specs"]
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
    grunt.loadNpmTasks("grunt-ts");
    grunt.loadNpmTasks("grunt-dts-bundle");
    grunt.loadNpmTasks("grunt-contrib-watch");
    grunt.loadNpmTasks('grunt-contrib-jasmine');

    grunt.registerTask("default", ["watch"]);
    grunt.registerTask("build", ["ts:dist", "dts_bundle:dist", "browserify:dist"]);
    grunt.registerTask("test", ["browserify:specs", "jasmine"]);
};