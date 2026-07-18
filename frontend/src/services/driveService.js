import { apiClient } from './api'

// Service de gestion des fichiers
class DriveService {
  async listFiles() {
    return apiClient.listFiles()
  }

  async uploadFile(file) {
    return apiClient.uploadFile(file)
  }

  async downloadFile(fileId) {
    const response = await apiClient.downloadFile(fileId)
    return response
  }

  // Helper pour déclencher un téléchargement
  async downloadAndSave(fileId, filename) {
    try {
      const response = await this.downloadFile(fileId)
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename || 'file'
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (e) {
      console.error('Download failed:', e)
      throw e
    }
  }

  // Formater la taille en bytes lisible
  formatSize(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  // Formater la date
  formatDate(timestamp) {
    return new Date(timestamp * 1000).toLocaleDateString('fr-FR')
  }
}

export const driveService = new DriveService()

