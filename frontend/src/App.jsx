import React, {useState} from 'react'
import Login from './Login'
import Drive from './Drive'
import AdminPanel from './AdminPanel'

export default function App(){
  const [token,setToken] = useState(null)
  const [isAdmin,setIsAdmin] = useState(false)

  if(!token) return <Login onLogin={(t,admin)=>{setToken(t); setIsAdmin(admin)}} />

  return (
    <div className="container">
      <header className="mb-8">
        <h1 className="text-4xl font-bold">La Familly Bass</h1>
        <p className="text-sm opacity-80">Minimal drive pour l'association — monochrome</p>
      </header>
      <main>
        <Drive token={token} />
        {isAdmin && <div className="mt-8"><AdminPanel token={token} /></div>}
      </main>
    </div>
  )
}

