import React, { useState } from 'react'
import { Login, Drive, AdminPanel } from './components'

export default function App() {
  const [token, setToken] = useState(null)
  const [isAdmin, setIsAdmin] = useState(false)

  if (!token) {
    return <Login onLogin={(t, admin) => { setToken(t); setIsAdmin(admin) }} />
  }

  return (
    <div className="container">
      <header className="mb-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-4xl font-bold">La Familly Bass</h1>
            <p className="text-sm opacity-80">Minimal drive pour l'association — monochrome</p>
          </div>
          <button
            onClick={() => { setToken(null); setIsAdmin(false) }}
            className="px-4 py-2 border border-white text-sm"
          >
            Logout
          </button>
        </div>
      </header>

      <main>
        <Drive />
        {isAdmin && (
          <div className="mt-8">
            <AdminPanel />
          </div>
        )}
      </main>
    </div>
  )
}

