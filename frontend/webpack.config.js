'use strict';

const process = require('process');
const path = require('path');
const buildPath = path.resolve(__dirname, 'static');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const CopyWebpackPlugin = require('copy-webpack-plugin');
//const TransferWebpackPlugin = require('transfer-webpack-plugin'); // dev-server only
let devtool = 'source-map'; // Render source-map file for final build
const plugins = [
  new CopyWebpackPlugin({
    patterns: [
      { context: '.', from: 'css/**/*.css' },
      { context: '.', from: 'html/**/*.html' },
      { context: '.', from: 'image/**/*.png' },
      { context: '.', from: 'image/**/*.svg' },
      { context: nodeModulesPath, from: 'bootstrap/dist/**', to: 'lib' },
      { context: nodeModulesPath, from: 'bootstrap-icons/**', to: 'lib' },
      { context: nodeModulesPath, from: 'bootswatch/dist/yeti/**', to: 'lib' },
      { context: nodeModulesPath, from: 'jquery/dist/**', to: 'lib' },
      { context: nodeModulesPath, from: 'react/umd/**', to: 'lib' },
      { context: nodeModulesPath, from: 'react-dom/umd/**', to: 'lib' },
    ],
  }),
];

if (process.env.WEBPACK_DEV_SERVER === 'true') {
  console.log('WEBPACK_DEV_SERVER is true');
  devtool = 'eval';
}

const config = {
  mode: process.env.MINIFY === 'true' ? 'production' : 'development',
  entry: {
    main: './js/main.tsx',
    setting: './js/setting.tsx',
  },
  resolve: {
    //When require, do not have to add these extensions to file's name
    extensions: ['.js', '.jsx', '.json', '.css', '.ts', '.tsx'],
    modules: ['web_modules', 'node_modules'] // (Default Settings)
  },
  output: {
    path: path.join(buildPath, process.env.VERSION_HASH),
    publicPath: path.join('/static', process.env.VERSION_HASH),
    filename: 'js/[name].bundle.js',
  },
  externals: {
    jquery: 'jQuery',
    react: 'React',
    'react-dom': 'ReactDOM',
    bootstrap: 'bootstrap',
    bootswatch: 'bootswatch',
  },
  optimization: {
    runtimeChunk: {
      name: 'vendor',
    },
    splitChunks: {
      name: 'vendor',
      chunks: 'initial',
    },
  },
  plugins: plugins,
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        //include: paths.appSrc,
        loader: require.resolve('babel-loader'),
        exclude: /node_modules/,
        options: {
          // This is a feature of `babel-loader` for Webpack (not Babel itself).
          // It enables caching results in ./node_modules/.cache/babel-loader/
          // directory for faster rebuilds.
          cacheDirectory: true,
          //plugins: ['react-hot-loader/babel'],
          presets: [
            ['@babel/react'],
            [
              '@babel/env',
              {
                targets: {
                  browsers: ['last 2 versions', 'safari >= 7'],
                },
              },
            ],
          ],
        },
      },
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.png$/,
        use: 'url-loader?limit=100000',
      },
      {
        test: /\.jpg$/,
        use: 'file-loader',
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url?limit=10000&mimetype=application/font-woff',
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url?limit=10000&mimetype=application/octet-stream',
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        use: 'file',
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url?limit=10000&mimetype=image/svg+xml',
      },
    ],
  },
  // Dev server Configuration options
  devServer: {
    //    contentBase: 'frontend', // Relative directory for base of server
    hot: true, // Live-reload
    port: 4000,
    host: '0.0.0.0', // Change to '0.0.0.0' for external facing server
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
      },
    },
    static: {
      directory: path.resolve(__dirname, '.'),
    },
  },
};

module.exports = config;
