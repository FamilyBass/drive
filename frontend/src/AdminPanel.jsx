import React, {useState} from 'react'

export default function AdminPanel({token}){
  const [userId,setUserId] = useState('')
  const [msg,setMsg] = useState('')

  const validate = async ()=>{
    const res = await fetch('/api/admin/validate',{method:'POST',headers:{'Content-Type':'application/json','Authorization':'Bearer '+token},body:JSON.stringify({user_id:userId})})
    if(res.ok){ setMsg('Utilisateur validé') } else setMsg('Erreur')
  }

  return (
    <div className="p-4 border border-white">
      <h4 className="font-semibold">Admin — Valider utilisateur</h4>
      <div className="mt-2">
        <input className="p-2 bg-black border border-white" placeholder="user id" value={userId} onChange={e=>setUserId(e.target.value)} />
        <button className="ml-2 px-3 py-1 border border-white" onClick={validate}>Valider</button>
      </div>
      {msg && <div className="mt-2 text-sm">{msg}</div>}
    </div>
  )
}

