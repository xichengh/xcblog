/**
 * whatwg-fetch implement standard Fetch specification, but to use in real project,
 * we need some addtional features
 */
import queryString from 'query-string'

async function request(req: any) {
  const loading: any = document.getElementById('root-loading-bar')
  if (loading) { loading.style.display = 'block'}
  const res = await req
  if (!res.ok) {
    return Promise.reject('服务端错误')
  }
  const resdata = await await res.json()
  if (resdata.msg === 'success') {
    if (loading) { loading.style.display = 'none' }
    return Promise.resolve(resdata.data)
  } else {
    if (loading) { loading.style.display = 'none' }
    return Promise.reject(resdata.msg)
  }
}
class Client {
  /**
   * send http GET request
   * @param  {String} url       absolute URL with another site`s host is performed in accordance with CORS
   * @param  {Object} [params]  url search
   * @param  {Boolean}[options] return full response data or parsed data
   * @param  {Object} [options] fetch options, see https://github.github.io/fetch/ for more details
   * @return {Promise}
   */
  get(url: any, params?: any, options?: any) {
    const formedUrl = params ? `${url}?${queryString.stringify(params)}` : url
    options = options || {}
    options.headers = options.headers || {}

    const req = fetch(formedUrl, {
      ...options,
      method: 'GET',
      credentials: 'same-origin',
    })
    return request(req)
  }
  /**
   * send http POST request
   * @param  {String} url       absolute URL with another site`s host is performed in accordance with CORS
   * @param  {Object} [data]    body to be sent
   * @param  {Boolean}[options] return full response data or parsed data
   * @param  {Object} [options] fetch options, see https://github.github.io/fetch/ for more details
   * @return {Promise}
   */
  post(url: any, data?: any, options?: any) {
    const isFormData = data instanceof FormData
    const headers = options ? options.headers : {}
    if (!isFormData) {
      headers['Content-Type'] = 'application/json'
    }
    options = options || {}
    const req = fetch(url, {
      ...options,
      method: 'POST',
      body: isFormData ? data : JSON.stringify(data),
      credentials: 'same-origin',
      headers: headers,
    })
    return request(req)
  }
  /**
   * send http PUT request
   * @param  {String} url       absolute URL with another site`s host is performed in accordance with CORS
   * @param  {Object} [data]    body to be sent
   * @param  {Boolean}[options] return full response data or parsed data
   * @param  {Object} [options] fetch options, see https://github.github.io/fetch/ for more details
   * @return {Promise}
   */
  put(url: any, data?: any, options?: any) {
    options = options || {}
    options.headers = options.headers || {}
    const req = fetch(url, {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data),
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    })
    return request(req)
  }
  /**
   * send http DELETE request
   * @param  {String} url       absolute URL with another site`s host is performed in accordance with CORS
   * @param  {Object} [data]    body to be sent
   * @param  {Boolean}[options] return full response data or parsed data
   * @param  {Object} [options] fetch options, see https://github.github.io/fetch/ for more details
   * @return {Promise}
   */
  del(url: any, data?: any, options?: any) {
    options = options || {}
    options.headers = options.headers || {}
    const req = fetch(url, {
      ...options,
      method: 'DELETE',
      body: JSON.stringify(data),
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    })
    return request(req)
  }
}
const client = new Client()

export default client
