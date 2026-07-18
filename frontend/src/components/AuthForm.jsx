import React, { useState } from 'react'
import { Button } from './Button'
import { TextField } from './TextField'
import { Alert } from './Alert'
import { Spinner } from './Spinner'

// Composant pour l'authentification (réutilisable)
export const AuthForm = ({
  mode = 'login',
  onSubmit,
  isLoading = false,
  error = null,
  onModeChange = null,
}) => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (email && password) {
      await onSubmit(email, password)
    }
  }

  const isLogin = mode === 'login'
  const title = isLogin ? 'Se connecter' : 'Create Account'
  const submitText = isLogin ? 'Se connecter' : "S'enregistrer"

  return (
    <form onSubmit={handleSubmit} className="max-w-sm">
      {error && <Alert type="error" message={error} />}

      <TextField
        label="Email"
        type="email"
        value={email}
        onChange={setEmail}
        placeholder="Enter your email"
        disabled={isLoading}
      />

      <TextField
        label="Password"
        type="password"
        value={password}
        onChange={setPassword}
        placeholder="Enter your password"
        disabled={isLoading}
      />

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading || !email || !password}>
          {isLoading ? 'Processing...' : submitText}
        </Button>

        {onModeChange && (
          <Button
            type="button"
            variant="secondary"
            disabled={isLoading}
            onClick={() => onModeChange(isLogin ? 'register' : 'login')}
          >
            {isLogin ? "S'enregistrer" : 'Se connecter'}
          </Button>
        )}
      </div>

      {isLoading && <Spinner text="Please wait..." />}
    </form>
  )
}

