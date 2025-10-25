import React from 'react'

import Modal from '@components/common/modal/Modal'

import { useAppDispatch, useAppSelector } from '@hooks'
import { closeGameStatusModal } from '@reducers'

interface GameStatusModalProps {
  resetBoard: () => void
}

const GameStatusModal: React.FC<GameStatusModalProps> = ({ resetBoard }) => {
  const dispatch = useAppDispatch()
  const { isOpen, message } = useAppSelector((state) => state.gameStatusModal)

  const close = () => {
    resetBoard()
    dispatch(closeGameStatusModal())
  }

  return (
    <Modal isOpen={isOpen} showCloseButton={true} close={close}>
      <div className='text-center'>
        <h2 className='text-2xl font-bold'>Game Over!</h2>
        <p className='text-lg'>{message}</p>
        <button className='mt-10 bg-gray-700' onClick={close}>
          Reset
        </button>
      </div>
    </Modal>
  )
}

export default GameStatusModal
