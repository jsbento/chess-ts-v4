import { post, get } from '@utils'
import type { AppDispatch } from 'state/store'
import { login, logout } from '@reducers'
import type { User, SignInReq, SignUpReq, AuthResp } from '@types'

export const signOut = (dispatch: AppDispatch) => {
  window.localStorage.removeItem('token')
  dispatch(logout())
}

export const signIn = async (
  dispatch: AppDispatch,
  data: SignInReq,
): Promise<User | null> => {
  try {
    const resp = await post<SignInReq, AuthResp>('/users/signin', data)

    if (!resp) {
      return null
    }

    window.localStorage.setItem('token', resp.accessToken)

    dispatch(login(resp))
    return resp
  } catch (err) {
    console.log(err)
    return null
  }
}

export const signUp = async (
  dispatch: AppDispatch,
  data: SignUpReq,
): Promise<User | null> => {
  try {
    const resp = await post<SignUpReq, AuthResp>('/users/signup', data)

    if (!resp) {
      return null
    }

    window.localStorage.setItem('token', resp.accessToken)

    dispatch(login(resp))
    return resp
  } catch (err) {
    console.log(err)
    return null
  }
}

export const checkHealthz = async (): Promise<boolean> => {
  try {
    const resp = await get<null, { message: string }>('/ping', null)
    if (!resp) {
      return false
    }

    return resp.message === 'pong'
  } catch (err) {
    console.log(err)
    return false
  }
}
