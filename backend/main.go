package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/golang-jwt/jwt/v5"
    _ "modernc.org/sqlite"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

var db *sql.DB
var jwtSecret string

type key string

const (
    ctxUserID key = "user_id"
    ctxIsAdmin key = "is_admin"
)

func main() {
    dataDir := os.Getenv("DATA_DIR")
    if dataDir == "" {
        dataDir = "/data"
    }
    mediaDir := filepath.Join(dataDir, "media")
    if err := os.MkdirAll(mediaDir, 0o755); err != nil {
        log.Fatalf("failed create media dir: %v", err)
    }

    dbPath := filepath.Join(dataDir, "db.sqlite")
    var err error
    db, err = sql.Open("sqlite", fmt.Sprintf("file:%s?_fk=1", dbPath))
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    if err := db.Ping(); err != nil {
        log.Fatalf("ping db: %v", err)
    }
    if err := migrate(); err != nil {
        log.Fatalf("migrate: %v", err)
    }

    jwtSecret = os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "dev_secret_change_me"
    }

    // ensure admin user
    adminEmail := os.Getenv("ADMIN_EMAIL")
    adminPass := os.Getenv("ADMIN_PASSWORD")
    if adminEmail != "" && adminPass != "" {
        ensureAdmin(adminEmail, adminPass)
    }

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Post("/api/register", handleRegister)
    r.Post("/api/login", handleLogin)

    r.Group(func(r chi.Router) {
        r.Use(authMiddleware)
        r.Post("/api/drive/upload", handleUpload)
        r.Get("/api/drive/list", handleList)
        r.Get("/api/drive/download/{id}", handleDownload)

        r.Group(func(r chi.Router) {
            r.Use(adminOnly)
            r.Post("/api/admin/validate", handleValidateUser)
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port
    log.Printf("server listening %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("listen: %v", err)
    }
}

func migrate() error {
    stm := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE,
        password TEXT,
        is_active INTEGER DEFAULT 0,
        is_admin INTEGER DEFAULT 0,
        created_at INTEGER
    );
    CREATE TABLE IF NOT EXISTS files (
        id TEXT PRIMARY KEY,
        owner_id TEXT,
        filename TEXT,
        size INTEGER,
        path TEXT,
        created_at INTEGER
    );`
    _, err := db.Exec(stm)
    return err
}

func ensureAdmin(email, pass string) {
    var id string
    err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)
    if err == sql.ErrNoRows {
        // create
        hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
        uid := uuid.New().String()
        _, err := db.Exec("INSERT INTO users(id,email,password,is_active,is_admin,created_at) VALUES(?,?,?,?,?,?)",
            uid, email, string(hashed), 1, 1, time.Now().Unix())
        if err != nil {
            log.Printf("create admin error: %v", err)
        } else {
            log.Printf("created admin %s", email)
        }
    } else if err == nil {
        // make sure active/admin
        _, _ = db.Exec("UPDATE users SET is_active=1, is_admin=1 WHERE email=?", email)
    } else {
        log.Printf("ensureAdmin query error: %v", err)
    }
}

// --- Handlers ---

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var req struct{
        Email string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    if req.Email == "" || req.Password == "" {
        http.Error(w, "email and password required", http.StatusBadRequest)
        return
    }
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }
    id := uuid.New().String()
    _, err = db.Exec("INSERT INTO users(id,email,password,is_active,is_admin,created_at) VALUES(?,?,?,?,?,?)",
        id, req.Email, string(hashed), 0, 0, time.Now().Unix())
    if err != nil {
        http.Error(w, "email exists or db error", 500)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    var req struct{ Email, Password string }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", 400)
        return
    }
    var id, hash string
    var isActive, isAdmin int
    err := db.QueryRow("SELECT id,password,is_active,is_admin FROM users WHERE email = ?", req.Email).Scan(&id, &hash, &isActive, &isAdmin)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    if isActive == 0 {
        http.Error(w, "account not validated", http.StatusForbidden)
        return
    }
    if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": id,
        "admin": isAdmin==1,
        "exp": time.Now().Add(24*time.Hour).Unix(),
    })
    tokStr, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"token": tokStr})
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        h := r.Header.Get("Authorization")
        if h == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }
        var tok string
        fmt.Sscanf(h, "Bearer %s", &tok)
        t, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !t.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        claims, ok := t.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "invalid token claims", http.StatusUnauthorized)
            return
        }
        sub, _ := claims["sub"].(string)
        admin, _ := claims["admin"].(bool)
        ctx := context.WithValue(r.Context(), ctxUserID, sub)
        ctx = context.WithValue(ctx, ctxIsAdmin, admin)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func adminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if v := r.Context().Value(ctxIsAdmin); v == nil || v == false {
            http.Error(w, "admin only", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func handleValidateUser(w http.ResponseWriter, r *http.Request) {
    var req struct{ UserID string `json:"user_id"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
        http.Error(w, "bad request", 400)
        return
    }
    _, err := db.Exec("UPDATE users SET is_active=1 WHERE id=?", req.UserID)
    if err != nil {
        http.Error(w, "db error", 500)
        return
    }
    w.WriteHeader(200)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    uid, _ := ctx.Value(ctxUserID).(string)
    // limit size to 200MB to protect RAM
    r.Body = http.MaxBytesReader(w, r.Body, 200<<20)
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        http.Error(w, "file too large", 400)
        return
    }
    file, fh, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "missing file", 400)
        return
    }
    defer file.Close()
    id := uuid.New().String()
    dataDir := os.Getenv("DATA_DIR")
    if dataDir == "" { dataDir = "/data" }
    target := filepath.Join(dataDir, "media", id+"_"+filepath.Base(fh.Filename))
    out, err := os.Create(target)
    if err != nil {
        http.Error(w, "cannot save", 500)
        return
    }
    defer out.Close()
    // stream copy
    n, err := io.Copy(out, file)
    if err != nil {
        http.Error(w, "write error", 500)
        return
    }
    _, err = db.Exec("INSERT INTO files(id,owner_id,filename,size,path,created_at) VALUES(?,?,?,?,?,?)",
        id, uid, fh.Filename, n, target, time.Now().Unix())
    if err != nil {
        http.Error(w, "db error", 500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"id": id, "path": target})
}

func handleList(w http.ResponseWriter, r *http.Request) {
    uid, _ := r.Context().Value(ctxUserID).(string)
    isAdmin, _ := r.Context().Value(ctxIsAdmin).(bool)
    var rows *sql.Rows
    var err error
    if isAdmin {
        rows, err = db.Query("SELECT id,owner_id,filename,size,created_at FROM files ORDER BY created_at DESC LIMIT 100")
    } else {
        rows, err = db.Query("SELECT id,owner_id,filename,size,created_at FROM files WHERE owner_id=? ORDER BY created_at DESC LIMIT 100", uid)
    }
    if err != nil {
        http.Error(w, "db error", 500)
        return
    }
    defer rows.Close()
    type F struct{ ID, OwnerID, Filename string; Size int64; CreatedAt int64 }
    res := []F{}
    for rows.Next() {
        var f F
        rows.Scan(&f.ID, &f.OwnerID, &f.Filename, &f.Size, &f.CreatedAt)
        res = append(res, f)
    }
    json.NewEncoder(w).Encode(res)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    uid, _ := r.Context().Value(ctxUserID).(string)
    isAdmin, _ := r.Context().Value(ctxIsAdmin).(bool)
    var owner, path, filename string
    var size int64
    err := db.QueryRow("SELECT owner_id,path,filename,size FROM files WHERE id=?", id).Scan(&owner, &path, &filename, &size)
    if err != nil {
        http.Error(w, "not found", 404)
        return
    }
    if owner != uid && !isAdmin {
        http.Error(w, "forbidden", 403)
        return
    }
    f, err := os.Open(path)
    if err != nil {
        http.Error(w, "file error", 500)
        return
    }
    defer f.Close()
    w.Header().Set("Content-Length", strconv.FormatInt(size,10))
    w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
    http.ServeContent(w, r, filename, time.Now(), f)
}

