import { User } from '@types'

export type AuthResp = {
  user?: User
  token?: string
}

export type SignInReq = {
  identifier: string
  password: string
}

export type SignUpReq = {
  username: string
  email: string
  password: string
}
