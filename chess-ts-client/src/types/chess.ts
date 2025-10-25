export type PromotionPiece = 'q' | 'r' | 'b' | 'n'

export type Move = {
  from: string
  to: string
}

export type SelectedPiece = {
  id: number
  moves: number[]
}

export type EvaluationReq = {
  fen: string
}

export type EvaluationResp = {
  score: number
}

export type SearchPositionReq = {
  fen: string
  depth: number
  moveTime: number
}

export type SearchPositionResp = {
  move: string
}
