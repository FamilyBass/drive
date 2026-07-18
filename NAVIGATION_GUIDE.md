# 🗺️ Guide de Navigation - Architecture La Familly Bass

## 📂 Structure du Backend Go

### Pour Comprendre le Flux de Connexion Admin

1. **User Clique "Se connecter"** (Frontend)
   ↓
2. **AuthForm.jsx** → `useAuth()` hook
   ↓
3. **useAuth.login()** → `authService.login()` 
   ↓
4. **authService.login()** → `apiClient.login()`
   ↓
5. **Requête POST /api/login**
   ↓
6. **Backend: main.go** → Route `/api/login`
   ↓
7. **authHandler.Login()** → Reçoit email/password
   ↓
8. **authService.Login()** → Logique métier
   - Appelle `userRepo.GetByEmail()` (interface → SQLite)
   - Valide le password avec bcrypt
   - Génère JWT token
   ↓
9. **SQLiteUserRepository.GetByEmail()** → Query DB
   ↓
10. **DB retourne User** → AuthService crée JWT
    ↓
11. **JWT envoyé au Frontend** ← Réponse JSON
    ↓
12. **Frontend décide si admin** (examine `admin` claim)

---

## 📊 Architecture Backend

```
┌─────────────────────────────────────────────────┐
│                    HTTP Layer                    │
│    main.go + chi.Router + Middlewares           │
└─────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────┐
│              Adapter Layer (HTTP)                │
│  auth_handler.go, drive_handler.go, middleware  │
│  ← Gère les requêtes HTTP                       │
│  → Appelle Domain Services                      │
└─────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────┐
│           Domain Layer (Logique Métier)         │
│  auth_service.go, drive_service.go              │
│  ← Reçoit des handlers                          │
│  ↔ Utilise repositories (interfaces)            │
│  → Retourne résultats/erreurs                   │
└─────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────┐
│         Adapter Layer (Persistance)              │
│  sqlite_user_repository.go, sqlite_file_...     │
│  ← Reçoit appels de services                    │
│  → Query/Persist dans SQLite                    │
└─────────────────────────────────────────────────┘
           ↓
        ┌──────────┐
        │ SQLite   │
        │ Database │
        └──────────┘
```

---

## 🎨 Architecture Frontend

```
┌─────────────────────────────────────────────────┐
│                   App.jsx                        │
│  ← Gère état (token, isAdmin)                   │
│  → Rend Login ou Drive                          │
└─────────────────────────────────────────────────┘
       ↓              ↓
    Login.jsx      Drive.jsx + AdminPanel.jsx
       ↓              ↓
┌─────────────────────────────────────────────────┐
│            Composants Réutilisables             │
│  Button, TextField, AuthForm, FileUpload...     │
│  ← Props génériques (label, value, onChange)   │
│  → Callback pour action                        │
└─────────────────────────────────────────────────┘
       ↓
┌─────────────────────────────────────────────────┐
│              Hooks (Logique Métier)             │
│  useAuth, useDrive, useAdmin                    │
│  ← Reçoivent pas de paramètres                 │
│  → Retournent {state, methods, error}           │
└─────────────────────────────────────────────────┘
       ↓
┌─────────────────────────────────────────────────┐
│            Services (API + Métier)              │
│  authService, driveService, adminService       │
│  ← Reçoivent de paramètres métier              │
│  → Retournent résultats ou erreurs             │
└─────────────────────────────────────────────────┘
       ↓
┌─────────────────────────────────────────────────┐
│          API Client (HTTP Layer)                │
│  apiClient (centralisé dans index.js)           │
│  ← Endpoint, method, body                       │
│  → Réponse JSON ou erreur                      │
└─────────────────────────────────────────────────┘
       ↓
    HTTP/localhost:8080/api/...
```

---

## 🔍 Trace un Bug - Exemples

### Exemple 1: "Login Admin ne fonctionne pas"

**Checklist de débogage:**

1. **Frontend Console** (DevTools)
   ```
   - Vérifier API call: `console.log` dans api.js
   - Vérifier response: Onglet Network
   - Vérifier hook: `console.log` dans useAuth.js
   ```

2. **Backend Logs** (Terminal)
   ```
   - Vérifier route: main.go `/api/login`
   - Vérifier handler: auth_handler.go Login()
   - Vérifier service: auth_service.go Login()
   - Vérifier DB: SQLiteUserRepository.GetByEmail()
   ```

3. **Database** (SQLite)
   ```
   SELECT * FROM users WHERE email='admin@example.com';
   - Vérifier user existe
   - Vérifier is_active=1, is_admin=1
   - Vérifier password hash
   ```

**Chemins à vérifier:**
- Frontend: `api.js` → `authService.js` → `useAuth.js` → `Login.jsx`
- Backend: `main.go` → `auth_handler.go` → `auth_service.go` → `sqlite_user_repository.go` → DB

---

### Exemple 2: "Upload fichier échoue"

**Checklist:**

1. **Frontend** (DevTools Network)
   - Status 413? → Fichier > 200MB
   - Status 401? → Token expiré
   - Status 400? → Champ `file` manquant

2. **Backend Logs**
   - Vérifier route: `/api/drive/upload`
   - Vérifier middleware: Token valide?
   - Vérifier handler: `drive_handler.go Upload()`
   - Vérifier service: `drive_service.go UploadFile()`

3. **Filesystem**
   ```
   ls -la C:\dev\FamilyBass\drive\data\media\
   - Fichiers créés?
   - Permissions OK?
   ```

4. **Database**
   ```
   SELECT * FROM files WHERE owner_id='<user_id>';
   - Fichier enregistré?
   ```

---

## 🛠️ Ajouter une Nouvelle Fonctionnalité

### Cas: Ajouter "Supprimer un fichier"

**Backend (Bottom-Up)**

1. **Database** - Pas de changement (utiliser ID existant)

2. **Domain - Repository Interface** (`file_repository.go`)
   ```go
   type FileRepository interface {
       Delete(ctx context.Context, id string) error  // ← Ajouter
   }
   ```

3. **Domain - Service** (`drive_service.go`)
   ```go
   func (s *DriveService) DeleteFile(ctx context.Context, fileID, userID string, isAdmin bool) error {
       // Logique métier: vérifier propriétaire ou admin
       // Appeler fileRepo.Delete()
   }
   ```

4. **Adapter - Repository** (`sqlite_file_repository.go`)
   ```go
   func (r *SQLiteFileRepository) Delete(ctx context.Context, id string) error {
       // Query: DELETE FROM files WHERE id=?
   }
   ```

5. **Adapter - Handler** (`drive_handler.go`)
   ```go
   func (h *DriveHandler) Delete(w http.ResponseWriter, r *http.Request) {
       // Parser fileID du URL
       // Appeler driveService.DeleteFile()
   }
   ```

6. **Main** - Ajouter route
   ```go
   r.Delete("/api/drive/delete/{id}", driveHandler.Delete)
   ```

**Frontend (Top-Down)**

1. **Service** (`driveService.js`)
   ```js
   async deleteFile(fileId) {
       return apiClient.request(`/drive/delete/${fileId}`, { method: 'DELETE' })
   }
   ```

2. **Hook** (`useDrive.js`)
   ```js
   const deleteFile = useCallback(async (fileId) => {
       // Appeler driveService.deleteFile()
       // Mettre à jour files list
   }, [])
   ```

3. **Component** (`FileList.jsx`)
   ```jsx
   <button onClick={() => onDelete(file.ID)}>
       Delete
   </button>
   ```

4. **Integration** (`Drive.jsx`)
   ```jsx
   const { deleteFile } = useDrive()
   <FileList onDelete={deleteFile} />
   ```

---

## 📝 Conventions de Code

### Backend Go
```go
// Interfaces dans domain/repository
// Implémentations dans adapter/repository
// Logique dans domain/service
// HTTP dans adapter/http

// Erreurs descriptives
return fmt.Errorf("failed to get user: %w", err)

// Context pour timeouts
ctx := r.Context()  // Depuis request HTTP
```

### Frontend React
```jsx
// Composants: reçoivent props, return JSX
// Hooks: logique métier, return {state, methods}
// Services: appels API, return Promise

// Error handling dans hooks
try {
    await asyncOperation()
} catch (err) {
    setError(err.message)
}
```

---

## 🎓 Points Clés à Retenir

### Backend (Débogage)
1. ✅ Vérifier **type d'erreur** dans chaque couche
2. ✅ Utiliser **context** pour timeouts
3. ✅ **Logger** au bootstrap et erreurs critiques
4. ✅ **Interfaces** pour dépendances (testable)

### Frontend (Débogage)
1. ✅ **DevTools** Network tab pour requêtes
2. ✅ **Console** pour logs et erreurs
3. ✅ **Props drilling** vs Context (simple ici)
4. ✅ **Services** centralisés → faciles à tester

### Communication (Débogage)
1. ✅ Vérifier **status HTTP** (200, 401, 403, 500)
2. ✅ Vérifier **format JSON** request/response
3. ✅ Vérifier **Authorization header** (Bearer token)

---

## 🚀 Workflow Développement

```
1. Créer feature dans Backend (DB → Handler)
2. Tester avec curl/Postman
3. Créer feature dans Frontend (Service → Component)
4. Tester avec DevTools
5. Valider intégration complète
6. Committer + Documenter
```

---

## 📚 Fichiers Clés Par Use-Case

### "Je veux comprendre l'authentification"
→ `backend/internal/domain/service/auth_service.go`
→ `backend/internal/adapter/http/auth_handler.go`
→ `frontend/src/services/authService.js`

### "Je veux déboguer une erreur 401"
→ `backend/internal/adapter/http/middleware.go` (AuthMiddleware)
→ `frontend/src/services/authService.js` (restoreSession)
→ DevTools → Network tab

### "Je veux ajouter un champ utilisateur"
→ `backend/internal/domain/entity/user.go` (champ)
→ `backend/internal/adapter/repository/sqlite_user_repository.go` (SQL)
→ `frontend/src/components/TextField.jsx` (UI)

### "Je veux changer le design"
→ `frontend/src/index.css` (styles)
→ `frontend/src/components/*.jsx` (tailwind classes)

---

**Happy Debugging! 🐛**

