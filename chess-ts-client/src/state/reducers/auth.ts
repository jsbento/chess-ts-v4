import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { User } from '@types'

interface InitialAuthState {
  user: User | null
}

const initialState: InitialAuthState = {
  user: null,
}

export const authSlice = createSlice({
  name: 'auth',
  initialState: initialState,
  reducers: {
    login: (state, action: PayloadAction<User | null>) => {
      state.user = action.payload
    },
    logout: (state) => {
      state.user = null
    },
  },
})

export const { login, logout } = authSlice.actions

export default authSlice.reducer
