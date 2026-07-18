# 🚀 QUICK START - La Familly Bass Refactorisée

## 5 Minutes pour Démarrer

### Terminal 1️⃣ - Backend

```powershell
cd C:\dev\FamilyBass\drive
.\start-backend.ps1
```

**Output attendu:**
```
🚀 Starting backend on port 8080...
📧 Admin email: admin@example.com
🔐 Admin password: admin

server listening on :8080
```

### Terminal 2️⃣ - Frontend

```powershell
cd C:\dev\FamilyBass\drive
.\start-frontend.ps1
```

**Output attendu:**
```
🚀 Starting frontend development server...
📱 Open http://localhost:5173 in your browser

VITE v4.5.0 ready in XXX ms

➜ Local:   http://localhost:5173/
```

### 🌐 Browser

```
Ouvrir: http://localhost:5173
```

---

## 🔐 Credentials

```
Email: admin@example.com
Password: admin
```

---

## ✨ Testons!

1. **Login** → `admin@example.com` / `admin`
2. **Upload fichier** → Cliquer "Upload"
3. **Liste fichiers** → Voir fichier listés
4. **Download** → Cliquer "Download"
5. **Déconnexion** → Cliquer "Déconnexion"

---

## 🎨 Design Néon? Check!

- ✅ Titre avec glow vert
- ✅ Inputs cyan avec focus vert
- ✅ Boutons vert avec hover effect
- ✅ Transparence et blur effects
- ✅ Animations fluides

---

## 🐛 Déboguer? Voici comment:

### Backend Logs
- Terminal 1 affiche les logs
- Chercher erreurs avec "error" keyword

### Frontend Logs
- Ouvrir DevTools (F12)
- Onglet **Console** pour errors
- Onglet **Network** pour API calls

### Database
```powershell
# VSCode avec extension SQLite
# Ouvrir: C:\dev\FamilyBass\drive\data\db.sqlite
```

---

## 📁 Fichiers Importants

| Chemin | Raison |
|--------|--------|
| `backend/main.go` | Point d'entrée backend |
| `frontend/src/App.jsx` | Point d'entrée frontend |
| `frontend/src/index.css` | Styles néon |
| `frontend/src/components/` | Composants réutilisables |
| `backend/internal/domain/service/` | Logique métier |

---

## 🆘 Problèmes Courants

### "404 page not found"
- Backend lancé? Onglet Network dans DevTools
- Vérifier vite.config.js proxy

### "Login échoue"
- Admin@db existe? Vérifier Terminal 1 logs
- Token expiré? Rafraîchir page (Ctrl+F5)

### "Upload échoue"
- Fichier < 200MB?
- Token valide? Essayer re-login

### "Erreur de compilation"
```powershell
cd backend
go mod tidy
go build
```

---

## 📚 Documentation Complète

| Fichier | Contenu |
|---------|---------|
| `QUICK_START.md` | Ce fichier 👈 |
| `CHECKLIST.md` | Tous les fichiers créés ✅ |
| `NAVIGATION_GUIDE.md` | Comment déboguer 🔍 |
| `IMPLEMENTATION_SUMMARY.md` | Changements effectués 📝 |
| `README_IMPLEMENTATION.md` | Architecture complète 🏗️ |

---

## ⚡ Architecture en 2 Slides

### Backend Go
```
HTTP Request
    ↓
main.go (routes)
    ↓
authHandler/driveHandler (HTTP parsing)
    ↓
authService/driveService (business logic)
    ↓
sqliteUserRepository/sqliteFileRepository (DB)
    ↓
SQLite DB
```

### Frontend React
```
Component (JSX)
    ↓
Hook (useDrive, useAuth)
    ↓
Service (driveService, authService)
    ↓
API Client (apiClient)
    ↓
HTTP Request to Backend
```

---

## 🎯 Code Quality Checklist

- ✅ **Séparation** des responsabilités (hexagonal backend)
- ✅ **Réutilisabilité** des composants (modular frontend)
- ✅ **Testabilité** (interfaces Go, hooks React)
- ✅ **Maintenabilité** (code lisible et bien structuré)
- ✅ **Débogage** (logs clairs et structure claire)
- ✅ **Design** (néon moderne avec animations)

---

## 🚀 Prochaines Étapes

1. Ajouter tests unitaires
2. Déployer avec Docker
3. Optimiser pour Raspberry Pi
4. Ajouter refresh tokens
5. Implémenter caching

---

## 💡 Tips

- Utilisez VSCode avec Go extension
- Ouvrez DevTools toujours (F12)
- Consultez NAVIGATION_GUIDE.md pour déboguer
- Logs terminaux sont vos meilleurs amis

---

## 📞 Support

Besoin d'aide?
1. Vérifier NAVIGATION_GUIDE.md
2. Regarder Terminal logs (Backend)
3. Vérifier Console (Frontend DevTools)
4. Lire IMPLEMENTATION_SUMMARY.md

---

**Bon développement! 🎉**

*La Familly Bass - Minimal Drive*
*Version 1.0.0 - Fully Refactored*

