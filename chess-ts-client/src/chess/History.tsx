import React from 'react'

import { useAppSelector } from '@hooks'

const History: React.FC = () => {
  const { moves } = useAppSelector((state) => state.chessMoves)

  const renderRows = () => {
    const rows = []
    let movesNum = 1
    for (let i = 0; i < moves.length; i += 2) {
      rows.push(
        <tr key={i} className='border-y-2 border-[#333]'>
          <td className='py-2'>{movesNum}</td>
          <td className='py-2'>{moves[i]}</td>
          <td className='py-2'>{moves[i + 1] || ''}</td>
        </tr>,
      )
      movesNum++
    }

    return rows
  }

  return (
    <div className='p-2 rounded-lg border-2 border-[#333] w-1/2 h-[500px]'>
      <h2 className='font-semibold text-2xl text-center border-b-2 border-[#333] pb-4 h-[50px]'>
        Game History
      </h2>
      <div className='overflow-auto h-[430px]'>
        <table className='w-full text-center border-collapse'>
          <thead className='sticky top-0 bg-[#242424]'>
            <tr>
              <th className='w-1/3 py-2 z-10'>Move</th>
              <th className='w-1/3 py-2 z-10'>White</th>
              <th className='w-1/3 py-2 z-10'>Black</th>
            </tr>
          </thead>
          <tbody>{renderRows()}</tbody>
        </table>
      </div>
    </div>
  )
}

export default History
