import React from 'react'

import { getPieceImage } from '@utils'
import { PromotionPiece } from '@types'

const PromotionPieces: PromotionPiece[] = ['q', 'r', 'b', 'n']

interface PromotionSquareProps {
  turn: string
  onPromote: (promotionPiece: PromotionPiece) => void
}

const PromotionSquare: React.FC<PromotionSquareProps> = ({
  turn,
  onPromote,
}) => {
  return (
    <div className='w-full grid grid-cols-2 grid-rows-2'>
      {PromotionPieces.map((piece) => (
        <img
          key={piece}
          src={getPieceImage(piece, turn === 'b' ? 'black' : 'white')}
          alt={piece}
          className='cursor-pointer'
          onClick={() => onPromote(piece)}
        />
      ))}
    </div>
  )
}

export default PromotionSquare
