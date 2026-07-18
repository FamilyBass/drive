import React, { useEffect } from 'react'
import { useDrive } from '../hooks/useDrive'
import { FileUpload } from './FileUpload'
import { FileList } from './FileList'
import { Alert } from './Alert'

export default function Drive() {
  const {
    files,
    loadFiles,
    uploadFile,
    downloadFile,
    error,
    isLoading,
    uploadProgress,
    formatSize,
    formatDate,
  } = useDrive()

  useEffect(() => {
    loadFiles()
  }, [loadFiles])

  const handleFileSelected = async (file) => {
    try {
      await uploadFile(file)
    } catch (err) {
      // L'erreur est gérée par le hook
    }
  }

  return (
    <div>
      <section className="mb-8">
        <h2 className="text-2xl font-semibold mb-4">Drive</h2>

        {error && (
          <Alert type="error" message={error} />
        )}

        <FileUpload
          onFileSelected={handleFileSelected}
          isLoading={isLoading}
          uploadProgress={uploadProgress}
        />
      </section>

      <section>
        <h3 className="text-xl font-medium mb-4">Fichiers</h3>
        <FileList
          files={files}
          onDownload={downloadFile}
          isLoading={isLoading}
          formatSize={formatSize}
          formatDate={formatDate}
        />
      </section>
    </div>
  )
}

