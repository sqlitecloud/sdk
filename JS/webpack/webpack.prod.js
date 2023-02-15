var path = require('path');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const TerserPlugin = require("terser-webpack-plugin");

module.exports = merge(common, {
  mode: 'production',
  devtool: 'source-map',
  output: {
    path: path.resolve(__dirname, '../dist'),
    filename: '[name].bundle.js',
    clean: true,
    library: {
      name: 'sqliteCloudJs',
      type: 'umd',
    },
  },
  optimization: {
    minimize: true,
    minimizer: [
      new TerserPlugin()
    ]
  }
});