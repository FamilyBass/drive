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
      // L'erreur est gérée par le hook
    }
  }

  const displayError = error || localError

  return (
    <div className="container min-h-screen flex items-center justify-center">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <h1 className="text-4xl font-bold mb-2">La Familly Bass</h1>
          <p className="text-sm text-cyan-300 opacity-80">
            {mode === 'login' ? 'Connexion' : 'Création de Compte'}
          </p>
        </div>

        <div className="p-6 rounded-lg border-2 border-cyan-500"
             style={{
               background: 'rgba(0, 255, 255, 0.05)',
               boxShadow: '0 0 20px rgba(0, 255, 255, 0.2), inset 0 0 10px rgba(0, 255, 255, 0.1)',
               backdropFilter: 'blur(10px)'
             }}>

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

          <div className="mt-6 text-center border-t border-green-500/30 pt-4">
            <p className="text-xs opacity-70 text-green-300">
              {mode === 'login'
                ? 'Inscription ouverte — Un admin doit valider le compte.'
                : 'Un administrateur doit valider votre compte avant de vous connecter.'}
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
