'use strict';

var webpack = require('webpack');
var fs = require('fs');
var Config = require('./hosting_config');
var banner = fs.readFileSync('./src/core/sqlite-cloud-licence.js', 'utf8');
banner = banner.replace('<VERSION>', Config.version);

module.exports = {
  entry: './src/core/index.js',
  plugins: [
    new webpack.BannerPlugin({ banner: banner, raw: true }),
    new webpack.DefinePlugin({
      VERSION: JSON.stringify(Config.version),
      // CDN_HTTP: JSON.stringify(Config.cdn_http),
      // CDN_HTTPS: JSON.stringify(Config.cdn_https),
    })
  ],
};