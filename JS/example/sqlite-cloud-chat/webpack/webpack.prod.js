const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const TerserPlugin = require("terser-webpack-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

module.exports = merge(common, {
  mode: 'production',
  devtool: 'source-map',
  output: {
    filename: '[name].[contenthash].bundle.js',
    clean: true
  },
  optimization: {
    minimize: true,
    minimizer: [
      new CssMinimizerPlugin(),      
      new TerserPlugin()
    ]
  }
});