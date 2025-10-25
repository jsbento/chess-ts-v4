import { Chess } from 'chess.js'
import { Move, PromotionPiece } from '@types'
export const getCharBoard = (chess: Chess): string[] => {
  const charBoard: string[] = []

  for (let rank = 0; rank < 8; rank++) {
    for (let file = 0; file < 8; file++) {
      const piece = chess?.board()[rank][file] || null
      if (piece) {
        charBoard.push(`${piece.color}${piece.type}`)
      } else {
        charBoard.push('empty')
      }
    }
  }

  return charBoard
}

export const getPossiblePromotions = (chess: Chess, from: string, to: string) =>
  chess
    .moves({ verbose: true })
    .filter(
      (move) => move.promotion && `${move.from}:${move.to}` === `${from}:${to}`,
    )

export const indexToRankFile = (index: number): [number, number] => {
  const rank = 8 - Math.floor(index / 8)
  const file = index % 8

  return [rank, file]
}

export const indexToSquare = (index: string | number): string => {
  const rank = 8 - Math.floor(Number(index) / 8)
  const file = 'abcdefgh'[Number(index) % 8]

  return `${file}${rank}`
}

export const buildMove = (
  from: string | number,
  to: string | number,
): { from: string; to: string } => {
  return {
    from: indexToSquare(from),
    to: indexToSquare(to),
  }
}

export const squareToIndex = (square: string): number => {
  const [file, rank] = square.split('')
  return (8 - Number(rank)) * 8 + 'abcdefgh'.indexOf(file)
}

export const getPieceImage = (piece: string, color: string): string => {
  return `./src/assets/chess-pieces//Chess_${piece}${color === 'black' ? 'd' : 'l'}t60.png`
}

export const getGameStatus = (chess: Chess): string => {
  if (chess.isCheckmate()) {
    return `Checkmate! ${chess.turn() === 'w' ? 'Black' : 'White'} wins!`
  } else if (chess.isDraw()) {
    return 'Draw!'
  } else if (chess.isInsufficientMaterial()) {
    return 'Draw by Insufficient Material!'
  } else if (chess.isStalemate()) {
    return 'Draw by Stalemate!'
  } else if (chess.isThreefoldRepetition()) {
    return 'Draw by Threefold Repetition!'
  } else {
    return 'Something strange happened!'
  }
}

export const parseEngineMove = (
  move: string,
): { move: Move; promotion?: PromotionPiece } => {
  const from = move.slice(0, 2)
  const to = move.slice(2, 4)
  const promotion = move.length > 4 ? move.slice(4) : undefined

  return {
    move: {
      from,
      to,
    },
    promotion: promotion as PromotionPiece,
  }
}
