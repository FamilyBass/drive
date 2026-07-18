import React, { useState, useEffect } from 'react'
import { useAuth } from '../hooks/useAuth'
import { authService } from '../services/authService'
import { AuthForm } from './AuthForm'
import { Alert } from './Alert'

export default function Login({ onLogin }) {
  const [mode, setMode] = useState('login')
  const { register, login, error, isLoading } = useAuth()
  const [localError, setLocalError] = useState(null)

  useEffect(() => {
    // Essayer de restaurer la session
    const session = authService.restoreSession()
    if (session) {
      onLogin(session.token, session.isAdmin)
    }
  }, [onLogin])

  const handleSubmit = async (email, password) => {
    setLocalError(null)
    try {
      if (mode === 'login') {
        const result = await login(email, password)
        onLogin(result.token, result.isAdmin)
      } else {
        await register(email, password)
        setLocalError(null)
        setMode('login')
      }
    } catch (err) {
      // L'erreur est géré par le hook
    }
  }

  const displayError = error || localError

  return (
    <div className="container">
      <h2 className="text-2xl font-semibold mb-6">La Familly Bass — {mode === 'login' ? 'Connexion' : 'Inscription'}</h2>

      {displayError && (
        <Alert type="error" message={displayError} onDismiss={() => setLocalError(null)} />
      )}

      <AuthForm
        mode={mode}
        onSubmit={handleSubmit}
        isLoading={isLoading}
        error={displayError}
        onModeChange={setMode}
      />

      <p className="mt-6 text-xs opacity-70">
        {mode === 'login'
          ? 'Inscription ouverte — un admin doit valider le compte.'
          : 'Un administrateur doit valider votre compte avant de vous connecter.'}
      </p>
    </div>
  )
}

