import { createSlice } from '@reduxjs/toolkit'

export const gameStatusModalSlice = createSlice({
  name: 'gameStatusModal',
  initialState: {
    isOpen: false,
    message: undefined,
  },
  reducers: {
    openGameStatusModal: (state, action) => {
      state.isOpen = true
      state.message = action.payload
    },
    closeGameStatusModal: (state) => {
      state.isOpen = false
      state.message = undefined
    },
  },
})

export const { openGameStatusModal, closeGameStatusModal } =
  gameStatusModalSlice.actions

export default gameStatusModalSlice.reducer
