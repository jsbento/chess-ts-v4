import React, { useEffect } from 'react'
import { Outlet, Link } from 'react-router-dom'

import { signOut, checkHealthz } from '@behavior'
import { useAppDispatch, useAppSelector } from '@hooks'

import './Root.css'

const Root: React.FC = () => {
  const dispatch = useAppDispatch()

  const user = useAppSelector((state) => state.auth.user)

  useEffect(() => {
    const healthz = async () => {
      const isOk = await checkHealthz()
      if (!isOk) {
        console.log('Server is not running')
      } else {
        console.log('Server is running')
      }
    }
    healthz()
  }, [])

  const renderAuthorizedNav = () => {
    if (!user) {
      return (
        <li>
          <Link to={'/auth'}>Sign In</Link>
        </li>
      )
    }

    return (
      <>
        <li>
          <div className='flex flex-col'>
            <span className='font-semibold mb-3'>Hi {user.username}!</span>
            <button onClick={() => signOut(dispatch)}>Sign Out</button>
          </div>
        </li>
        <li>
          <Link to={'/profile'}>Profile</Link>
        </li>
      </>
    )
  }

  return (
    <>
      <div id='sidebar'>
        <h1>Chess TS</h1>
        <nav>
          <ul>
            {renderAuthorizedNav()}
            <li>
              <Link to={'/'}>Home</Link>
            </li>
            <li>
              <Link to={'/chess'}>Play</Link>
            </li>
          </ul>
        </nav>
      </div>
      <div id='detail'>
        <Outlet />
      </div>
    </>
  )
}

export default Root
