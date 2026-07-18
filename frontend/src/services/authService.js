import { apiClient } from './api'

// Utilitaires pour JWT
export const decodeToken = (token) => {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) {
      throw new Error('Invalid token format')
    }
    const payload = JSON.parse(atob(parts[1]))
    return payload
  } catch (e) {
    console.error('Failed to decode token:', e)
    return null
  }
}

export const isTokenExpired = (token) => {
  const payload = decodeToken(token)
  if (!payload || !payload.exp) {
    return true
  }
  return Date.now() >= payload.exp * 1000
}

// Gestion du stockage des tokens (sessionStorage pour la sécurité)
const TOKEN_KEY = 'familybass_token'

export const saveToken = (token) => {
  if (typeof window !== 'undefined') {
    sessionStorage.setItem(TOKEN_KEY, token)
  }
}

export const loadToken = () => {
  if (typeof window !== 'undefined') {
    return sessionStorage.getItem(TOKEN_KEY)
  }
  return null
}

export const clearToken = () => {
  if (typeof window !== 'undefined') {
    sessionStorage.removeItem(TOKEN_KEY)
  }
}

// Service d'authentification
class AuthService {
  async register(email, password) {
    const result = await apiClient.register(email, password)
    return result
  }

  async login(email, password) {
    const result = await apiClient.login(email, password)
    const token = result.token

    if (token) {
      saveToken(token)
      apiClient.setToken(token)
    }

    return {
      token,
      isAdmin: this.isAdmin(token),
    }
  }

  logout() {
    clearToken()
    apiClient.setToken(null)
  }

  isAdmin(token) {
    const payload = decodeToken(token)
    return payload?.admin === true
  }

  restoreSession() {
    const token = loadToken()
    if (token && !isTokenExpired(token)) {
      apiClient.setToken(token)
      return {
        token,
        isAdmin: this.isAdmin(token),
      }
    }
    clearToken()
    return null
  }
}

export const authService = new AuthService()

