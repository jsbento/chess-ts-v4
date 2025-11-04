export type Game = {
  id: string
  playerId: string
  moves: string
  result: string
  createdAt: string
  updatedAt: string
}

export type SaveGameReq = {
  playerId: string
  moves: string
  result: string
}