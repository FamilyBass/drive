import React, { useState } from 'react'
import { Button } from './Button'
import { TextField } from './TextField'
import { Spinner } from './Spinner'

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
  const title = isLogin ? 'Se connecter' : 'Créer un compte'
  const submitText = isLogin ? 'Connexion' : 'S\'enregistrer'

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <TextField
        label="Email"
        type="email"
        value={email}
        onChange={setEmail}
        placeholder="admin@example.com"
        disabled={isLoading}
      />

      <TextField
        label="Mot de passe"
        type="password"
        value={password}
        onChange={setPassword}
        placeholder="••••••••"
        disabled={isLoading}
      />

      <div className="flex gap-3 pt-4">
        <Button
          type="submit"
          disabled={isLoading || !email || !password}
          className="flex-1"
        >
          {isLoading ? 'Traitement...' : submitText}
        </Button>

        {onModeChange && (
          <Button
            type="button"
            variant="secondary"
            disabled={isLoading}
            onClick={() => onModeChange(isLogin ? 'register' : 'login')}
            className="flex-1"
          >
            {isLogin ? 'S\'enregistrer' : 'Se connecter'}
          </Button>
        )}
      </div>

      {isLoading && <Spinner text="Veuillez patienter..." />}
    </form>
  )
}
