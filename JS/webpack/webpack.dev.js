var path = require('path');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].bundle.js',
    clean: true,
    library: {
      name: 'sqliteCloudJS',
      type: 'umd',
    },
  },
  devServer: {
    static: './dist',
  },
});