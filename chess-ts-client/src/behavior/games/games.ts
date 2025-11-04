import { post, get } from '@utils'

import type {
  SaveGameReq,
  Game,
} from '@types'

export const saveGame = async (req: SaveGameReq): Promise<Game | null> => {
  try {
    const resp = await post<SaveGameReq, Game>('/games', req)
    return resp ?? null
  } catch (err) {
    console.log(err)
    return null
  }
}

export const getUserGames = async (): Promise<Game[]> => {
  try {
    const resp = await get<null, Game[]>('/games/user')
    return resp ?? []
  } catch (err) {
    console.log(err)
    return []
  }
}