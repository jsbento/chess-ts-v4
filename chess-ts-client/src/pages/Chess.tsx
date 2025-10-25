import React from 'react'

import { Board, History, EngineControl } from '@components'

const Chess: React.FC = () => {
  return (
    <div className='grid grid-cols-3'>
      <EngineControl />
      <div className='w-full'>
        <Board size={500} />
      </div>
      <div className='flex justify-center w-full'>
        <History />
      </div>
    </div>
  )
}

export default Chess
