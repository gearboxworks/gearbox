const CopyPlugin = require('copy-webpack-plugin');
module.exports = {
  publicPath: '',
  outputDir: undefined,
  assetsDir: undefined,
  runtimeCompiler: undefined,
  productionSourceMap: undefined,
  parallel: undefined,
  css: undefined,
  lintOnSave: undefined,
  configureWebpack: {
    plugins: [
      new CopyPlugin([
        { from: '**/gears.json', to: 'gears.json' }
      ])
    ]
  }
}
