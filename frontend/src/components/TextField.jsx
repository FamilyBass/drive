import React from 'react'

// Composant réutilisable de champ de formulaire
export const TextField = ({
  label,
  value,
  onChange,
  type = 'text',
  placeholder = '',
  disabled = false,
  error = null,
}) => {
  return (
    <div className="mb-4">
      {label && <label className="block text-sm mb-1">{label}</label>}
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        disabled={disabled}
        className={`w-full p-2 bg-black border text-white ${
          error ? 'border-red-500' : 'border-white'
        }`}
      />
      {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
    </div>
  )
}

