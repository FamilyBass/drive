import React, {useState} from 'react'

export default function Login({onLogin}){
  const [email,setEmail] = useState('')
  const [password,setPassword] = useState('')
  const [msg,setMsg] = useState('')

  const doLogin = async ()=>{
    try{
      const res = await fetch('/api/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({email,password})})
      if(!res.ok){ setMsg('Erreur connexion'); return }
      const j = await res.json()
      // best-effort decode token to read admin claim
      const parts = j.token.split('.')
      let isAdmin=false
      if(parts.length>1){
        try{ const payload = JSON.parse(atob(parts[1])); isAdmin = payload.admin }catch(e){}
      }
      onLogin(j.token, isAdmin)
    }catch(e){ setMsg('Erreur réseau') }
  }

  return (
    <div className="container">
      <h2 className="text-2xl font-semibold mb-4">La Familly Bass — Connexion</h2>
      <div className="max-w-sm">
        <label className="block text-sm mb-1">Email</label>
        <input className="w-full p-2 bg-black border border-white text-white" value={email} onChange={e=>setEmail(e.target.value)} />
        <label className="block text-sm mb-1 mt-3">Mot de passe</label>
        <input type="password" className="w-full p-2 bg-black border border-white text-white" value={password} onChange={e=>setPassword(e.target.value)} />
        <div className="mt-4">
          <button className="px-4 py-2 border border-white" onClick={doLogin}>Se connecter</button>
        </div>
        {msg && <p className="mt-3 text-red-400">{msg}</p>}
      </div>
      <p className="mt-6 text-xs opacity-70">Inscription ouverte — un admin doit valider le compte.</p>
    </div>
  )
}

