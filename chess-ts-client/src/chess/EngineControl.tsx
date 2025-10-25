import React, { useCallback } from 'react'

import { useAppSelector, useAppDispatch } from '@hooks'
import { NumberInput, Checkbox } from '@components/common'
import {
  setEngineActive as _setEngineActive,
  setDepth as _setDepth,
  setMoveTime as _setMoveTime,
} from '@reducers'

const EngineControl: React.FC = () => {
  const dispatch = useAppDispatch()
  const { engineActive, depth, moveTime } = useAppSelector(
    (state) => state.gameSettings,
  )

  const setEngineActive = useCallback(
    (value: boolean) => {
      dispatch(_setEngineActive(value))
    },
    [dispatch],
  )

  const setDepth = useCallback(
    (value: number) => {
      dispatch(_setDepth(value))
    },
    [dispatch],
  )

  const setMoveTime = useCallback(
    (value: number) => {
      dispatch(_setMoveTime(value))
    },
    [dispatch],
  )

  return (
    <div className='flex flex-col p-2 rounded-lg border-2 border-[#333] w-1/2'>
      <h2 className='font-semibold text-2xl text-center border-b-2 border-[#333] pb-4'>
        Engine Controls
      </h2>
      <Checkbox
        label='Engine Active'
        labelClassName='mt-2'
        checked={engineActive}
        onChange={() => setEngineActive(!engineActive)}
      />
      <NumberInput
        label='Depth'
        labelClassName='mt-2'
        value={depth}
        onChange={(value) => setDepth(value)}
      />
      <NumberInput
        label='Move Time (ms)'
        labelClassName='mt-2'
        value={moveTime}
        onChange={(value) => setMoveTime(value)}
      />
    </div>
  )
}

export default EngineControl
