var path = require('path');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'sqlitecloud-sdk.js',
    clean: true,
    library: {
      name: 'SQLiteCloud',
      type: 'umd',
    },
  },
  devServer: {
    static: './dist',
  },
});