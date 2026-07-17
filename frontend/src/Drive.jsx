import React, {useEffect, useState, useRef} from 'react'

export default function Drive({token}){
  const [files,setFiles] = useState([])
  const fileRef = useRef()

  const listFiles = async ()=>{
    const res = await fetch('/api/drive/list',{headers:{'Authorization':'Bearer '+token}})
    if(res.ok){ setFiles(await res.json()) }
  }
  useEffect(()=>{ listFiles() },[])

  const upload = async ()=>{
    const f = fileRef.current.files[0]
    if(!f) return
    const fd = new FormData()
    fd.append('file', f)
    const res = await fetch('/api/drive/upload',{method:'POST',body:fd, headers:{'Authorization':'Bearer '+token}})
    if(res.ok){ listFiles() }
  }

  return (
    <div>
      <section className="mb-6">
        <h3 className="text-lg font-medium">Drive</h3>
        <div className="mt-3">
          <input ref={fileRef} type="file" className="text-sm" />
          <button className="ml-2 px-3 py-1 border border-white" onClick={upload}>Upload</button>
        </div>
      </section>
      <section>
        <h4 className="text-md font-medium mb-2">Fichiers</h4>
        <ul className="list-none">
          {files.map(f=> (
            <li key={f.ID} className="py-2 border-b border-white/10 flex justify-between">
              <div>
                <div className="font-medium">{f.Filename}</div>
                <div className="text-xs opacity-70">{f.OwnerID} • {f.Size} bytes</div>
              </div>
              <div>
                <a className="px-2 py-1 border border-white" href={`/api/drive/download/${f.ID}`} >Download</a>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}

