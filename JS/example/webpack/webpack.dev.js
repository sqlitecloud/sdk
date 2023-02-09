const path = require('path');
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    allowedHosts: [
      'sqlitecloud'
    ],
    devMiddleware: {
      writeToDisk: true
    },
  },
  output: {
    filename: '[name].bundle.js',
    clean: true,
  }
});