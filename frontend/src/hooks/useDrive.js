import { useState, useCallback } from 'react'
import { driveService } from '../services/driveService'

// Hook pour la gestion des fichiers
export const useDrive = () => {
  const [files, setFiles] = useState([])
  const [error, setError] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)

  const loadFiles = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const result = await driveService.listFiles()
      setFiles(result || [])
    } catch (err) {
      setError('Failed to load files')
      console.error(err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  const uploadFile = useCallback(async (file) => {
    setIsLoading(true)
    setError(null)
    setUploadProgress(0)

    try {
      // Simuler une progression (la véritable progression nécessiterait XMLHttpRequest)
      setUploadProgress(50)
      const result = await driveService.uploadFile(file)
      setUploadProgress(100)

      // Recharger les fichiers
      await loadFiles()
      return result
    } catch (err) {
      setError('Upload failed: ' + (err.message || 'Unknown error'))
      console.error(err)
      throw err
    } finally {
      setIsLoading(false)
      setUploadProgress(0)
    }
  }, [loadFiles])

  const downloadFile = useCallback(async (fileId, filename) => {
    setError(null)
    try {
      await driveService.downloadAndSave(fileId, filename)
    } catch (err) {
      setError('Download failed')
      console.error(err)
    }
  }, [])

  return {
    files,
    loadFiles,
    uploadFile,
    downloadFile,
    error,
    isLoading,
    uploadProgress,
    formatSize: driveService.formatSize,
    formatDate: driveService.formatDate,
  }
}

