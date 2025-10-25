import { createSlice } from '@reduxjs/toolkit'

interface InitialChessMovesState {
  moves: string[]
}

const initialState: InitialChessMovesState = {
  moves: [],
}

export const chessMovesSlice = createSlice({
  name: 'chessMoves',
  initialState: initialState,
  reducers: {
    addMove: (state, action) => {
      state.moves = [...state.moves, action.payload]
    },
    clearMoves: (state) => {
      state.moves = []
    },
  },
})

export const { addMove, clearMoves } = chessMovesSlice.actions

export default chessMovesSlice.reducer
