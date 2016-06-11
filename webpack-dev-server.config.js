const webpack = require('webpack');
const path = require('path');
const buildPath = path.resolve(__dirname, 'build');
const nodeModulesPath = path.resolve(__dirname, 'node_modules');
const TransferWebpackPlugin = require('transfer-webpack-plugin');

const config = {
  //Entry points to the project
  entry: [
    'webpack/hot/dev-server',
    'webpack/hot/only-dev-server',
    path.join(__dirname, '/src/app/main.js'),
  ],
  //Config options on how to interpret requires imports
  resolve: {
    extensions: ["", ".js"],
    //node_modules: ["web_modules", "node_modules"]  (Default Settings)
  },
  //Server Configuration options
  devServer:{
    contentBase: 'src/www',  //Relative directory for base of server
    devtool: 'eval',
    hot: true,        //Live-reload
    inline: true,
    port: 4000,        //Port Number
    host: 'localhost',  //Change to '0.0.0.0' for external facing server
  },
  devtool: 'eval',
  output: {
    path: buildPath,    //Path of output file
    filename: 'main.js',
  },
  plugins: [
    //Enables Hot Modules Replacement
    new webpack.HotModuleReplacementPlugin(),
    //Allows error warnings but does not stop compiling. Will remove when eslint is added
    new webpack.NoErrorsPlugin(),
    //Moves files
    new TransferWebpackPlugin([
      {from: 'www'},
    ], path.resolve(__dirname, "src")),
  ],
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
  //eslint config options. Part of the eslint-loader package
  eslint: {
    configFile: '.eslintrc',
  },
};

module.exports = config;

