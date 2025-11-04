import React, { useState, useEffect, useCallback } from 'react'

import { useAppSelector } from '@hooks'
import { getUserGames } from '@behavior'
import type { Game } from '@types'

const Profile: React.FC = () => {
  const user = useAppSelector((state) => state.auth.user)
  const [games, setGames] = useState<Game[]>([])
  if (!user) {
    return null
  }

  const fetchUserGames = useCallback(async () => {
    if (!user) return
    const games = await getUserGames()
    setGames(games)
  }, [user.id])

  useEffect(() => {
    fetchUserGames()
  }, [fetchUserGames])

  return (
    <div className='w-full'>
      <h1>Hi {user.username}!</h1>
      <div className='flex flex-col gap-4'>
        {games.map((game) => (
          <div key={game.id}>
            <h2>{game.moves}</h2>
            <p>{game.result}</p>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Profile
