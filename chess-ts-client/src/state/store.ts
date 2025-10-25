import { configureStore } from '@reduxjs/toolkit'
import storage from 'redux-persist/lib/storage'
import { persistReducer, persistStore } from 'redux-persist'
import {
  gameStatusModalReducer,
  chessMovesReducer,
  authReducer,
  gameSettingsReducer,
} from '@reducers'

const persistUserConfig = {
  key: 'auth',
  storage,
  whitelist: ['user'],
}
const persistAuthReducer = persistReducer(persistUserConfig, authReducer)

const store = configureStore({
  devTools: import.meta.env.MODE !== 'production',
  reducer: {
    gameStatusModal: gameStatusModalReducer,
    chessMoves: chessMovesReducer,
    auth: persistAuthReducer,
    gameSettings: gameSettingsReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
      },
    }),
})
export default store
export const persistor = persistStore(store)

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
