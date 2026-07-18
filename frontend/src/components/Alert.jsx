import React from 'react'

// Composant pour afficher des messages d'erreur ou d'alerte
export const Alert = ({ type = 'error', message, onDismiss = null }) => {
  if (!message) return null

  const bgClass = type === 'error' ? 'bg-red-900/20 border-red-500' : 'bg-green-900/20 border-green-500'
  const textClass = type === 'error' ? 'text-red-400' : 'text-green-400'

  return (
    <div className={`border p-3 mb-4 ${bgClass} ${textClass}`}>
      <div className="flex justify-between items-center">
        <span>{message}</span>
        {onDismiss && (
          <button onClick={onDismiss} className="text-sm opacity-70 hover:opacity-100">
            ✕
          </button>
        )}
      </div>
    </div>
  )
}

