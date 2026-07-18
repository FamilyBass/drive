import React, { useRef } from 'react'
import { Button } from './Button'
import { Spinner } from './Spinner'

// Composant de téléchargement de fichier réutilisable
export const FileUpload = ({ onFileSelected, isLoading = false, uploadProgress = 0 }) => {
  const fileRef = useRef()

  const handleUpload = () => {
    const file = fileRef.current?.files?.[0]
    if (file) {
      onFileSelected(file)
    }
  }

  return (
    <div className="border border-white p-4 mb-6">
      <h3 className="text-lg font-medium mb-3">Upload a file</h3>

      <div className="flex gap-2">
        <input
          ref={fileRef}
          type="file"
          className="text-sm flex-1"
          disabled={isLoading}
        />
        <Button onClick={handleUpload} disabled={isLoading || !fileRef.current?.files?.length}>
          {isLoading ? 'Uploading...' : 'Upload'}
        </Button>
      </div>

      {uploadProgress > 0 && uploadProgress < 100 && (
        <div className="mt-3">
          <div className="w-full bg-black border border-white/20">
            <div
              className="bg-white h-1"
              style={{ width: `${uploadProgress}%`, transition: 'width 0.3s' }}
            />
          </div>
          <p className="text-xs mt-1 opacity-70">{uploadProgress}%</p>
        </div>
      )}

      {isLoading && uploadProgress === 0 && <Spinner text="Uploading..." />}
    </div>
  )
}

