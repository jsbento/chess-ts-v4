import axios from 'axios'

export const _axios = axios.create({
  baseURL: import.meta.env.VITE_API_HOST,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: false,
  timeout: 2500,
})

_axios.interceptors.request.use((config) => {
  const token = window.localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

_axios.interceptors.response.use((response) => {
  return response
})

export const post = async <Req, Res>(
  uri: string,
  data: Req,
): Promise<Res | null> => {
  try {
    const response = await _axios.post<Res>(uri, data)
    return response.data
  } catch (err) {
    console.log(err)
    return null
  }
}

export const get = async <Req, Res>(
  uri: string,
  params: Req,
): Promise<Res | null> => {
  try {
    const response = await _axios.get<Res>(uri, { params })
    return response.data
  } catch (err) {
    console.log(err)
    return null
  }
}

export const put = async <Req, Res>(
  uri: string,
  data: Req,
): Promise<Res | null> => {
  try {
    const response = await _axios.put<Res>(uri, data)
    return response.data
  } catch (err) {
    console.log(err)
    return null
  }
}

export const del = async <Res>(uri: string): Promise<Res | null> => {
  try {
    const response = await _axios.delete<Res>(uri)
    return response.data
  } catch (err) {
    console.log(err)
    return null
  }
}
