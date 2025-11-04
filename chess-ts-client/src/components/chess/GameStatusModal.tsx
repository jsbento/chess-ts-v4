import React from 'react'

import Modal from '@components/common/modal/Modal'

import { useAppDispatch, useAppSelector } from '@hooks'
import { closeGameStatusModal } from '@reducers'

import { saveGame } from '@behavior'

interface GameStatusModalProps {
  resetBoard: () => void
}

const GameStatusModal: React.FC<GameStatusModalProps> = ({ resetBoard }) => {
  const dispatch = useAppDispatch()
  const { isOpen, message, user, movesList } = useAppSelector((state) => ({
    isOpen: state.gameStatusModal.isOpen,
    message: state.gameStatusModal.message,
    user: state.auth.user,
    movesList: state.chessMoves.moves,
  }))

  const close = () => {
    resetBoard()
    dispatch(closeGameStatusModal())
  }

  const onSaveGame = async () => {
    if (!user || !message) return

    const moves = movesList.join(',')
    const game = await saveGame({
      playerId: user.id,
      moves,
      result: message,
    })
    if (!game) return

    close()
  }

  return (
    <Modal isOpen={isOpen} showCloseButton={true} close={close}>
      <div className='text-center'>
        <h2 className='text-2xl font-bold'>Game Over!</h2>
        <p className='text-lg'>{message}</p>
        <button className='mt-10 bg-[#242424] text-white' onClick={close}>
          Reset
        </button>
        {user && (
          <button className='mt-10 bg-[#242424] text-white' onClick={onSaveGame}>
            Save Game
          </button>
        )}
      </div>
    </Modal>
  )
}

export default GameStatusModal
