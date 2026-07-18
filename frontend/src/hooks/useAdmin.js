import { useState, useCallback } from 'react'
import { adminService } from '../services/adminService'

// Hook pour les opérations admin
export const useAdmin = () => {
  const [error, setError] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState(null)

  const validateUser = useCallback(async (userId) => {
    setIsLoading(true)
    setError(null)
    setMessage(null)

    try {
      await adminService.validateUser(userId)
      setMessage('User validated successfully')
      return true
    } catch (err) {
      setError('Failed to validate user: ' + (err.message || 'Unknown error'))
      console.error(err)
      return false
    } finally {
      setIsLoading(false)
    }
  }, [])

  const clearMessages = useCallback(() => {
    setError(null)
    setMessage(null)
  }, [])

  return {
    validateUser,
    error,
    message,
    isLoading,
    clearMessages,
  }
}

