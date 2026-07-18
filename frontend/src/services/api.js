// Service centralisé pour les appels API
const API_BASE = '/api'

class APIClient {
  constructor(token = null) {
    this.token = token
  }

  setToken(token) {
    this.token = token
  }

  async request(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      })

      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}`
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          // Si la réponse n'est pas JSON, utiliser le texte
          const text = await response.text()
          errorMessage = text || errorMessage
        }
        const error = new Error(errorMessage)
        error.status = response.status
        throw error
      }

      // Pour les téléchargements, retourner la réponse directement
      if (options.skipJSON) {
        return response
      }

      return response.json()
    } catch (err) {
      console.error(`API Error [${endpoint}]:`, err)
      throw err
    }
  }

  // Auth endpoints
  async register(email, password) {
    return this.request('/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  async login(email, password) {
    return this.request('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  // Drive endpoints
  async uploadFile(file) {
    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await fetch(`${API_BASE}/drive/upload`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${this.token}`,
        },
        body: formData,
      })

      if (!response.ok) {
        throw new Error('Upload failed')
      }

      return response.json()
    } catch (err) {
      console.error('Upload Error:', err)
      throw err
    }
  }

  async listFiles() {
    return this.request('/drive/list', { method: 'GET' })
  }

  async downloadFile(fileId) {
    return this.request(`/drive/download/${fileId}`, {
      skipJSON: true,
      method: 'GET',
    })
  }

  // Admin endpoints
  async validateUser(userId) {
    return this.request('/admin/validate', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  }
}

// Instance globale du client API
export const apiClient = new APIClient()
