import store from '../state/store'

export type RootState = ReturnType<typeof store.getState>
