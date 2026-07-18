import React, { useState } from 'react'
import { useAdmin } from '../hooks/useAdmin'
import { TextField } from './TextField'
import { Button } from './Button'
import { Alert } from './Alert'

export default function AdminPanel() {
  const { validateUser, error, message, isLoading, clearMessages } = useAdmin()
  const [userId, setUserId] = useState('')

  const handleValidate = async () => {
    if (userId.trim()) {
      const success = await validateUser(userId)
      if (success) {
        setUserId('')
      }
    }
  }

  return (
    <div className="p-4 border border-white">
      <h3 className="text-xl font-semibold mb-4">Admin — Validate Users</h3>

      {error && (
        <Alert type="error" message={error} onDismiss={clearMessages} />
      )}
      {message && (
        <Alert type="success" message={message} onDismiss={clearMessages} />
      )}

      <div className="max-w-sm">
        <TextField
          label="User ID"
          value={userId}
          onChange={setUserId}
          placeholder="Enter user ID to validate"
          disabled={isLoading}
        />

        <Button
          onClick={handleValidate}
          disabled={isLoading || !userId.trim()}
          variant="primary"
        >
          {isLoading ? 'Validating...' : 'Validate User'}
        </Button>
      </div>
    </div>
  )
}

