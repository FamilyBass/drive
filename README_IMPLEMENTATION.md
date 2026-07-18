# La Familly Bass — Minimal Drive 🎵

Un service de gestion de fichiers minimaliste optimisé pour une Raspberry Pi 3, avec un design néon moderne.

## 📋 Prérequis

- **Backend** : Go 1.21+
- **Frontend** : Node.js 16+, npm
- **Données** : Répertoire `/data` (créé automatiquement)

## 🚀 Démarrage Rapide

### Terminal 1 - Backend

```powershell
cd C:\dev\FamilyBass\drive
.\start-backend.ps1
```

Le backend démarre sur **http://localhost:8080** avec :
- 📧 Admin email : `admin@example.com`
- 🔐 Admin password : `admin`

### Terminal 2 - Frontend

```powershell
cd C:\dev\FamilyBass\drive
.\start-frontend.ps1
```

Le frontend démarre sur **http://localhost:5173**

## 🎨 Design

- **Néon vert/cyan** avec animations fluides
- **Transparence** pour un effet moderne
- **Monospace** (Courier New) pour l'atmosphère tech

## 📁 Architecture

### Backend (Go - Hexagonal)

```
internal/
├── domain/           # Entités métier
│   ├── entity/
│   ├── repository/   # Interfaces
│   └── service/      # Logique métier
├── adapter/          # Couche application
│   ├── http/         # Handlers & middlewares
│   └── repository/   # Implémentations SQLite
└── bootstrap/        # Injection de dépendances
```

### Frontend (React - Composants réutilisables)

```
src/
├── components/       # Composants réutilisables
│   ├── Button, TextField, Alert, Spinner
│   ├── FileUpload, FileList
│   ├── Login, Drive, AdminPanel
│   └── AuthForm
├── hooks/           # Custom hooks
│   ├── useAuth
│   ├── useDrive
│   └── useAdmin
├── services/        # Services métier
│   ├── api.js
│   ├── authService.js
│   ├── driveService.js
│   └── adminService.js
└── index.css        # Styles néon
```

## 🔑 Fonctionnalités

### Authentification
- Inscription/Login
- JWT tokens (24h d'expiration)
- Gestion des rôles (user/admin)
- Session persistent (sessionStorage)

### Drive
- Upload de fichiers (max 200MB)
- Listing des fichiers
- Téléchargement
- Affichage des métadonnées (taille, date)

### Admin
- Validation des nouveaux utilisateurs
- Vue d'ensemble des fichiers

## 🗄️ Base de Données

SQLite stockée dans `data/db.sqlite` avec tables :
- `users` : Authentification & rôles
- `files` : Métadonnées des fichiers

Fichiers physiques stockés dans `data/media/`

## 🌐 API Endpoints

### Auth
- `POST /api/register` - Créer un compte
- `POST /api/login` - Authentifier

### Drive
- `POST /api/drive/upload` - Télécharger un fichier
- `GET /api/drive/list` - Lister les fichiers
- `GET /api/drive/download/{id}` - Télécharger un fichier

### Admin
- `POST /api/admin/validate` - Valider un utilisateur

## 🔒 Sécurité

- Passwords hashés avec bcrypt
- JWT avec secret configurable
- Tokens en sessionStorage (pas localStorage)
- Middleware d'authentification sur routes protégées
- Limite de taille upload (200MB)

## 📝 Variables d'Environnement (Backend)

```
DATA_DIR = "/data"              # Répertoire des données
PORT = "8080"                   # Port du serveur
JWT_SECRET = "change_me"        # Secret JWT
ADMIN_EMAIL = "admin@..."       # Email admin initial
ADMIN_PASSWORD = "password"     # Password admin initial
```

## 🛠️ Développement

### Compiler le backend
```powershell
cd C:\dev\FamilyBass\drive\backend
go build -o backend.exe
```

### Build du frontend
```powershell
cd C:\dev\FamilyBass\drive\frontend
npm install
npm run build
```

Les fichiers produits se trouvent dans `frontend/dist/`

## 📦 Déploiement Docker

Le projet inclut des fichiers Docker pour déploiement sur Raspberry Pi :

```bash
docker-compose up --build
```

Voir `docker-compose.yml` pour les détails.

## 🎯 Notes d'Optimisation Raspberry Pi

- Backend en Go compilé statiquement (minimal footprint)
- Frontend servi par nginx en static
- SQLite pour la persistence légère
- Limitations :
  - Max 200MB par fichier upload
  - Max 100 fichiers listés
  - Tokens JWT 24h

## 📄 Licence

Privé - La Familly Bass

