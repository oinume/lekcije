'use strict';

const process = require('process');
const webpack = require('webpack');
const path = require('path');
const buildPath = path.resolve(__dirname, 'static');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const TransferWebpackPlugin = require('transfer-webpack-plugin'); // dev-server only

var devtool = 'source-map'; // Render source-map file for final build
var plugins = [
  new webpack.NoErrorsPlugin(),
  new CopyWebpackPlugin([
    { context: 'src', from: '**/*.css' },
    { context: 'src', from: '**/*.html' },
    { context: 'src', from: '**/*.png' },
    { context: 'src', from: '**/*.jpg' },
    { context: 'src', from: '**/*.svg' },
    { context: nodeModulesPath, from: 'bootstrap/dist/**', to: 'lib' },
    { context: nodeModulesPath, from: 'bootswatch/**', to: 'lib' },
    { context: nodeModulesPath, from: 'jquery/dist/**', to: 'lib' },
  ])
];

if (process.env.MINIFY === 'true') {
  console.log('MINIFY = true');
  plugins.push(
    // Minify the bundle
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        //supresses warnings, usually from module minification
        warnings: false,
      }
    })
  );
}

if (process.env.WEBPACK_DEV_SERVER === 'true') {
  console.log('WEBPACK_DEV_SERVER is true');
  devtool = 'eval';
  plugins.push(
    new TransferWebpackPlugin([
      {from: 'css'},
      {from: 'html'},
      {from: 'image'},
      {from: nodeModulesPath + "/bootstrap", to: 'lib'},
      {from: nodeModulesPath + "/bootswatch", to: 'lib'},
      {from: nodeModulesPath + "/jquery", to: 'lib'},
    ], path.resolve(__dirname, "src"))
  );
}

const config = {
  entry: path.join(__dirname, '/src/js/main.js'),
  resolve: {
    //When require, do not have to add these extensions to file's name
    extensions: ["", ".js"],
    //node_modules: ["web_modules", "node_modules"]  (Default Settings)
  },
  devtool: devtool,
  //output config
  output: {
    path: path.join(buildPath, process.env.VERSION_HASH),
    publicPath: "/static/" + process.env.VERSION_HASH,
    filename: 'js/main.js',  // Name of output file
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
      // {
      //   test: /\//,
      //   loader: 'string-replace',
      //   query: {
      //     search: '$staticUrl$',
      //     replace: 'http://static.local.lekcije.com/static'
      //   }
      // },
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
        target: 'http://localhost:4001',
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
