'use strict';

const process = require('process');
const webpack = require('webpack');
const path = require('path');
const buildPath = path.resolve(__dirname, 'static');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const TransferWebpackPlugin = require('transfer-webpack-plugin'); // dev-server only

var devtool = 'source-map'; // Render source-map file for final build
var entry = [];
var plugins = [
  new webpack.NoErrorsPlugin(),
  new CopyWebpackPlugin([
    { context: 'src', from: '**/*.html' },
    { context: 'src', from: '**/*.css' },
    { context: nodeModulesPath, from: 'bootstrap/dist/**', to: 'lib' },
    { context: nodeModulesPath, from: 'bootswatch/**', to: 'lib' },
    { context: nodeModulesPath, from: 'jquery/dist/**', to: 'lib' },
  ])
]
if (process.env.NODE_ENV === 'production') {
  console.log('NODE_ENV is production');
  plugins.push(
    // Minify the bundle
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        //supresses warnings, usually from module minification
        warnings: false,
      },
    })
  );
} else {
  console.log('NODE_ENV is ' + process.env.NODE_ENV);
  devtool = 'eval';
  entry.push(
    'webpack/hot/dev-server',
    'webpack/hot/only-dev-server'
  );

  plugins.push(
    // Enables Hot Modules Replacement
    new webpack.HotModuleReplacementPlugin(),
    new TransferWebpackPlugin([
      {from: 'html'},
      {from: 'css'},
      {from: nodeModulesPath + "/bootstrap", to: 'lib'},
      {from: nodeModulesPath + "/bootswatch", to: 'lib'},
      {from: nodeModulesPath + "/jquery", to: 'lib'},
    ], path.resolve(__dirname, "src"))
  );
}

entry.push(path.join(__dirname, '/src/app/main.js'));

const config = {
  entry: entry,
  resolve: {
    //When require, do not have to add these extensions to file's name
    extensions: ["", ".js"],
    //node_modules: ["web_modules", "node_modules"]  (Default Settings)
  },
  devtool: devtool,
  //output config
  output: {
    path: buildPath,
    publicPath: "/static/",
    filename: 'main.js',  // Name of output file
  },
  plugins: plugins,
  module: {
    loaders: [
      {
        //React-hot loader and
        test: /\.jsx?$/,  //All .js files
        loaders: ['react-hot', 'babel-loader'], //react-hot is like browser sync and babel loads jsx and es6-7
        // query: {
        //   presets: ['react', 'es2015']
        // },
        exclude: [nodeModulesPath]
      },
      {
        test: /\.css$/,
        loader: "style-loader!css-loader"
      },
      {
        test: /\.png$/,
        loader: "url-loader?limit=100000"
      },
      {
        test: /\.jpg$/,
        loader: "file-loader"
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=application/octet-stream'
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file'
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url?limit=10000&mimetype=image/svg+xml'
      }
    ],
  },
  // Eslint config
  eslint: {
    configFile: '.eslintrc', //Rules for eslint
  },
  // Dev server Configuration options
  devServer: {
    contentBase: 'src',  // Relative directory for base of server
    devtool: 'eval',
    hot: true,        // Live-reload
    inline: true,
    port: 4000,
    host: '0.0.0.0',  // Change to '0.0.0.0' for external facing server
    proxy: {
      '*': {
        target: 'http://localhost:5000',
        secure: false,
        // bypass: function (req, res, proxyOptions) {
        //     const accept = req.headers.accept
        //     console.log(accept);
        //     //if (accept.indexOf('html') !== -1 || accept.indexOf('js') !== -1 || accept.indexOf('css') !== -1) {
        //         console.log('Skipping proxy for browser request.');
        //         return false;
        //     //}
        // }
      }
    },
  },
};

module.exports = config;
