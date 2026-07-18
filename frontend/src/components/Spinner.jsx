import React from 'react'

// Composant pour afficher un loader
export const Spinner = ({ text = 'Loading...' }) => {
  return (
    <div className="flex items-center gap-2 py-4">
      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
      <span className="text-sm">{text}</span>
    </div>
  )
}

