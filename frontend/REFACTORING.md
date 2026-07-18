# Frontend Refactorisation

## Structure refactorisée

Le frontend a été refactorisé pour suivre une architecture modulaire avec séparation des préoccupations.

### Avant
- Logique mélangée avec le rendering dans Login.jsx, Drive.jsx, AdminPanel.jsx
- Appels API directement dans les composants
- Pas de réutilisabilité de code

### Après

#### 📁 Composants (`src/components/`)
- `Button.jsx` - Bouton réutilisable
- `TextField.jsx` - Champ de saisie réutilisable
- `Alert.jsx` - Affichage des alertes/messages
- `Spinner.jsx` - Indicateur de chargement
- `FileUpload.jsx` - Upload de fichiers réutilisable
- `FileList.jsx` - Liste de fichiers réutilisable
- `AuthForm.jsx` - Formulaire d'authentification
- `Login.jsx` (refactorisé) - Utilise AuthForm et les hooks
- `Drive.jsx` (refactorisé) - Utilise FileUpload, FileList et useDrive
- `AdminPanel.jsx` (refactorisé) - Utilise useAdmin et les composants

#### 🎣 Hooks (`src/hooks/`)
- `useAuth.js` - Gestion de l'authentification
- `useDrive.js` - Gestion des fichiers et uploads
- `useAdmin.js` - Gestion des opérations admin

#### 🔧 Services (`src/services/`)
- `api.js` - Client HTTP centralisé
- `authService.js` - Service d'authentification avec gestion des tokens
- `driveService.js` - Service de gestion des fichiers
- `adminService.js` - Service admin

### Avantages

✅ **Séparation des préoccupations** - Logique métier séparée des composants
✅ **Réutilisabilité** - Composants et hooks réutilisables
✅ **Testabilité** - Services facilement mockables
✅ **Maintenabilité** - Code plus lisible et facile à maintenir
✅ **Debugging** - Logique claire et traceable
✅ **Sécurité** - Tokens stockés en sessionStorage

### Migration

Les anciens fichiers (Login.jsx, Drive.jsx, AdminPanel.jsx) du dossier `src/` peuvent être supprimés.
Les nouveaux fichiers refactorisés se trouvent dans les sous-dossiers.

### Déploiement

```bash
npm run build
```

Le build produit les fichiers optimisés dans `dist/`.

