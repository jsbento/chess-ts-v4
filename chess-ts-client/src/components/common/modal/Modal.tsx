import React, { useEffect, useRef, KeyboardEvent } from 'react'
import './Modal.css'

interface ModalProps {
  isOpen: boolean
  showCloseButton?: boolean
  close: () => void
  children: React.ReactNode
}

const Modal: React.FC<ModalProps> = ({
  isOpen,
  showCloseButton,
  close,
  children,
}) => {
  const ref = useRef<HTMLDialogElement | null>(null)

  useEffect(() => {
    const modalElem = ref.current
    if (modalElem) {
      if (isOpen) {
        modalElem.showModal()
      } else {
        modalElem.close()
      }
    }
  }, [isOpen])

  const onKeyDown = (e: KeyboardEvent<HTMLDialogElement>) => {
    if (e.key === 'Escape') {
      close()
    }
  }

  return (
    <dialog ref={ref} className='modal' onKeyDown={onKeyDown}>
      <div className='modal-content'>
        {showCloseButton && (
          <button className='modal-close-btn' onClick={close}>
            &times;
          </button>
        )}
        {children}
      </div>
    </dialog>
  )
}

export default Modal
