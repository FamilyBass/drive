import React from 'react'
import { Button } from './Button'
import { Spinner } from './Spinner'

// Composant de liste de fichiers réutilisable
export const FileList = ({
  files,
  onDownload,
  isLoading = false,
  formatSize,
  formatDate,
}) => {
  if (isLoading) {
    return <Spinner text="Loading files..." />
  }

  if (!files || files.length === 0) {
    return <p className="text-sm opacity-70">No files yet.</p>
  }

  return (
    <ul className="list-none">
      {files.map((file) => (
        <li key={file.ID} className="py-3 border-b border-white/10 flex justify-between items-center">
          <div className="flex-1">
            <div className="font-medium">{file.Filename}</div>
            <div className="text-xs opacity-70">
              {file.OwnerID} • {formatSize(file.Size)} • {formatDate(file.CreatedAt)}
            </div>
          </div>
          <div>
            <Button
              variant="secondary"
              onClick={() => onDownload(file.ID, file.Filename)}
              className="text-sm"
            >
              Download
            </Button>
          </div>
        </li>
      ))}
    </ul>
  )
}

