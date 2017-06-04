module.exports = function (grunt) {
    grunt.initConfig({
        ts: {
            dist: {
                src: ["./src/background.ts"],
                outDir: "./dist/ts",
                options: {
                    "rootDir": "./src",
                    "noImplicitAny": true,
                    "target": "es5",
                    "allowJs": true
                }
            },
        },
        browserify: {
            dist: {
                files: {
                    "./dist/background.js": ["./dist/ts/background.js"],
                },
                options: {
                    plugin: ['browserify-derequire']
                }
            },
        },
        watch: {
            scripts: {
                files: ["./src/background.ts", "../node_module/dist/module.js"],
                tasks: ["build"]
            },
        },
    });

    grunt.loadNpmTasks("grunt-browserify");
    grunt.loadNpmTasks("grunt-ts");
    grunt.loadNpmTasks("grunt-contrib-watch");

    grunt.registerTask("default", ["watch"]);
    grunt.registerTask("build", ["ts:dist", "browserify:dist"]);
};