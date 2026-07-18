import { useState, useCallback } from 'react'
import { authService } from '../services/authService'

// Hook pour la gestion de l'authentification
export const useAuth = () => {
  const [error, setError] = useState(null)
  const [isLoading, setIsLoading] = useState(false)

  const register = useCallback(async (email, password) => {
    setIsLoading(true)
    setError(null)
    try {
      const result = await authService.register(email, password)
      return result
    } catch (err) {
      const message = err.message || 'Registration failed'
      setError(message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const login = useCallback(async (email, password) => {
    setIsLoading(true)
    setError(null)
    try {
      const result = await authService.login(email, password)
      return result
    } catch (err) {
      let message = 'Login failed'
      if (err.status === 403) {
        message = 'Account not validated by admin'
      } else if (err.status === 401) {
        message = 'Invalid credentials'
      } else {
        message = err.message || message
      }
      setError(message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const logout = useCallback(() => {
    authService.logout()
    setError(null)
  }, [])

  return {
    register,
    login,
    logout,
    error,
    isLoading,
  }
}

