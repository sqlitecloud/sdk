var path = require('path');
var Config = require('./hosting_config');
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const TerserPlugin = require("terser-webpack-plugin");

var minimize = process.env.MINIMIZE === 'false' ? false : true;
var filename = minimize
  ? 'sqlitecloud-nodejs-sdk.min.js'
  : 'sqlitecloud-nodejs-sdk.js';


const distFolder = Config.version.slice(0,3);

let buildOption = {
  mode: 'production',
  output: {
    path: path.resolve(__dirname, '../dist/'+distFolder+'/'),
    filename: filename,
    clean: false,
    globalObject: 'this',
    library: {
      name: 'SQLiteCloud',
      type: 'umd',
    },
  },
  optimization: {
    minimize: minimize,
    minimizer: [
      new TerserPlugin()
    ]
  }
}

if (minimize) {
  buildOption['devtool'] = 'source-map'
}

module.exports = merge(common, buildOption);