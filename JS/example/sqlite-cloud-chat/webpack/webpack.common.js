const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const devMode = process.env.NODE_ENV !== "production";
const ASSET_PATH = process.env.ASSET_PATH || '';

module.exports = {
  context: path.resolve(__dirname, '../app'),
  entry: {
    app: './src/entry/app.js',
  },
  output: {
    publicPath: ASSET_PATH,
    path: path.resolve(__dirname, '../public'),
  },
  plugins: [
    new Dotenv(),
    new HtmlWebpackPlugin({
      template: './template/index.html',
      title: 'SQLite Cloud Chats',
      filename: "index.html"
    })
  ].concat(
    devMode ?
      []
      :
      [
        new MiniCssExtractPlugin({
          filename: devMode ? "[name].css" : "[name].[contenthash].css",
          chunkFilename: devMode ? "[id].css" : "[id].[contenthash].css",
          ignoreOrder: false, // Enable to remove warnings about conflicting order
        })
      ]
  ),
  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env']
          }
        }
      },
      {
        test: /\.css$/i,
        use: [
          devMode ? "style-loader" : MiniCssExtractPlugin.loader,
          'css-loader',
          'postcss-loader'
        ],
      },
      {
        test: /\.(png|svg|jpg|jpeg|gif)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'static/image/[hash][ext][query]'
        }
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'static/fonts/[hash][ext][query]'
        }
      }
    ],
  },
  optimization: {
    runtimeChunk: 'single',
  },
  stats: {
    children: false,
  }
};

