import { createSlice } from '@reduxjs/toolkit'

interface GameSettingsState {
  engineActive: boolean
  depth: number
  moveTime: number
}

const initialState: GameSettingsState = {
  engineActive: false,
  depth: 3,
  moveTime: 1000,
}

const gameSettingsSlice = createSlice({
  name: 'gameSettings',
  initialState,
  reducers: {
    setEngineActive: (state, action) => {
      state.engineActive = action.payload
    },
    setDepth: (state, action) => {
      state.depth = action.payload
    },
    setMoveTime: (state, action) => {
      state.moveTime = action.payload
    },
  },
})

export const { setEngineActive, setDepth, setMoveTime } =
  gameSettingsSlice.actions

export default gameSettingsSlice.reducer
