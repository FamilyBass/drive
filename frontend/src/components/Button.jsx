import React from 'react'

// Composant réutilisable de bouton
export const Button = ({
  children,
  onClick,
  disabled = false,
  variant = 'primary',
  type = 'button',
  className = '',
}) => {
  const baseClass = 'px-4 py-2 border'

  let variantClass = 'border-white'
  if (variant === 'secondary') {
    variantClass = 'border-white/50'
  } else if (variant === 'danger') {
    variantClass = 'border-red-500 text-red-500'
  }

  const disabledClass = disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`${baseClass} ${variantClass} ${disabledClass} ${className}`}
    >
      {children}
    </button>
  )
}

