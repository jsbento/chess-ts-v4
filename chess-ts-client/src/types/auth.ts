import type { User } from '@types'

export type AuthResp = User & {
  accessToken: string
  refreshToken: string
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
