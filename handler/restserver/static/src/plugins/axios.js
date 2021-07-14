import Axios from 'axios'

export default {
  install: (Vue) => {
    const axios = Axios.create()
    axios.interceptors.request.use(config => {
      if (!config.headers) {
        config.headers = {}
      }
      const token = sessionStorage.getItem("apiKey")
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })
    Vue.prototype.$axios = axios
  }
}
