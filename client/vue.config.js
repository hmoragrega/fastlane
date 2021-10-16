module.exports = {
    transpileDependencies: ["vuetify"],
    pluginOptions: {
        electronBuilder: {
            builderOptions: {
                directories: {
                    buildResources: "build",
                },
                appId: "com.hmoragrega.fastlane",
                mac: {
                    category: "public.app-category.developer-tools",
                    target: "dmg",
                    icon: "build/git.png"
                },
                extraResources: [
                    "src/assets/*"
                ],
                asar: true,
                asarUnpack: [
                    "node_modules/node-notifier/vendor/**"
                ],
            }
        }
    }
};
