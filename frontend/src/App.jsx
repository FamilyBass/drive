import React, { useState } from 'react'
import { Login, Drive, AdminPanel } from './components'
import { apiClient } from './services/api'

export default function App() {
  const [token, setToken] = useState(null)
  const [isAdmin, setIsAdmin] = useState(false)

  const handleLogin = (t, admin) => {
    // Définir le token dans le client API global
    apiClient.setToken(t)
    setToken(t)
    setIsAdmin(admin)
  }

  const handleLogout = () => {
    apiClient.setToken(null)
    setToken(null)
    setIsAdmin(false)
  }

  if (!token) {
    return <Login onLogin={handleLogin} />
  }

  return (
    <div className="container">
      <header className="mb-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-4xl font-bold">La Familly Bass</h1>
            <p className="text-sm opacity-80 text-cyan-300">Minimal drive pour l'association</p>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 border border-green-500 text-green-300 hover:bg-green-500/10 transition"
          >
            Déconnexion
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
