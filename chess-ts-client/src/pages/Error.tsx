import React from 'react'
import { useRouteError } from 'react-router-dom'

const Error: React.FC = () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const error = useRouteError() as any

  return (
    <div className='flex flex-col items-center justify-center h-screen w-screen'>
      <h1>Something unexpected happened...</h1>
      <p className='mt-4'>
        <i className='text-red-500'>{error.statusText || error.message}</i>
      </p>
    </div>
  )
}

export default Error
