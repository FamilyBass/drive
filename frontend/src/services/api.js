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

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const text = await response.text()
      const error = new Error(text || `HTTP ${response.status}`)
      error.status = response.status
      throw error
    }

    // Pour les téléchargements, retourner la réponse directement
    if (options.skipJSON) {
      return response
    }

    return response.json()
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

