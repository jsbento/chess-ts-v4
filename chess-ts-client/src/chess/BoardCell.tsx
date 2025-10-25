import React from 'react'
import { useDroppable } from '@dnd-kit/core'

import { indexToSquare } from '@utils'

interface BoardCellProps {
  id: string
  color: string
  size: number
  highlight?: boolean
  children?: React.ReactNode
}

const BoardCell: React.FC<BoardCellProps> = ({
  id,
  color,
  size,
  highlight,
  children,
}) => {
  const { isOver, setNodeRef } = useDroppable({
    id,
  })
  const style = {
    width: size,
    height: size,
    opacity: isOver ? 0.75 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`${highlight ? 'bg-teal-200 border-2 border-teal-400' : color} w-full flex justify-center items-center relative bg-blend-screen`}
    >
      <p
        style={{ position: 'absolute', top: 0, left: 1, fontSize: '0.5em' }}
        className='text-black'
      >
        {indexToSquare(parseInt(id))}
      </p>
      {children}
    </div>
  )
}

export default BoardCell
