const path = require('path')
module.exports = {
  pages: {
    index: {
      entry: 'src/main.js',
      title: "Role Based N Security"
    }
  },
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api/v1': {
        target: 'http://api-rbac-dev:9443',
        secure: false
      }
    },
    disableHostCheck: true
  },
  configureWebpack: {
    resolve: {
      extensions: ['.webpack.js', '.web.js', '.js', '.vue'],
      alias: {
        '@': path.resolve(__dirname, 'src'),
        '@assets': path.resolve(__dirname, 'src', "assets"),
        '@plugins': path.resolve(__dirname, 'src', 'plugins'),
        '@page': path.resolve(__dirname, 'src', 'components', 'pages'),
        '@tpl': path.resolve(__dirname, 'src', 'components', 'templates'),
        '@org': path.resolve(__dirname, 'src', 'components', 'organisms'),
        '@mixin': path.resolve(__dirname, 'src', 'mixins'),
      }
    }
  }
}
