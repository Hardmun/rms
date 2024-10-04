package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"rms/logs"
	"rms/sqLite"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed rms@.service
var embedFiles embed.FS

type fieldMap map[string]string

type tableStruct struct {
	Collection fieldMap `json:"collection"`
	Images     fieldMap `json:"images"`
	Details    fieldMap `json:"details"`
}

type emailStruct struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	SMTPServer string `json:"SMTPServer"`
	SMTPPort   string `json:"SMTPPort"`
	EmailTo    string `json:"emailTo"`
}

type SettingsStruct struct {
	MU                 sync.RWMutex
	ServerPGSQL        string        `json:"serverPgsql"`
	PortPGSQL          string        `json:"portPgsql"`
	DatabasePGSQL      string        `json:"databasePgsql"`
	LoginPGSQL         string        `json:"loginPgsql"`
	PasswordPGSQL      string        `json:"passwordPgsql"`
	HTTPServer         string        `json:"httpServer"`
	HttpPort           string        `json:"httpPort"`
	UrlPort            string        `json:"urlPort"`
	HttpRedirectPrefix string        `json:"httpRedirectPrefix"`
	FileDirectory      string        `json:"fileDirectory"`
	RequestBlockMin    int64         `json:"requestBlockMin"`
	ConnLimit          int           `json:"connLimit"`
	UserTtlMin         time.Duration `json:"userTtlMin"`
	EmailTech          emailStruct   `json:"emailTech"`
	Tables             tableStruct   `json:"tables"`
}

type user struct {
	rateLimit chan struct{}
	timeLogin time.Time
	hashPass  string
}

type usersCache struct {
	users     sync.Map
	ttl       time.Duration
	connLimit int
}

func (u *usersCache) getUser(username string) (user, bool) {
	client, exists := u.users.Load(username)
	if exists {
		return client.(user), true
	}

	return user{}, false
}

func (u *usersCache) addUser(username, dbPassword string) user {
	usr := user{
		rateLimit: make(chan struct{}, u.connLimit),
		timeLogin: time.Now().UTC(),
		hashPass:  dbPassword,
	}
	u.users.Store(username, usr)

	return usr
}

func (u *usersCache) StartCleanupJob() {
	ticker := time.NewTicker(u.ttl * time.Minute)
	go func() {
		for range ticker.C {
			u.cleanupExpiredClients()
		}
	}()
}

func (u *usersCache) cleanupExpiredClients() {
	now := time.Now().UTC()
	u.users.Range(func(key, value interface{}) bool {
		usr := value.(user)
		if now.Sub(usr.timeLogin) > u.ttl {
			u.users.Delete(key)
		}
		return true
	})
}

func newClientCache() *usersCache {
	Params.MU.RLock()
	defer Params.MU.RUnlock()

	return &usersCache{
		ttl:       Params.UserTtlMin,
		connLimit: Params.ConnLimit,
	}
}

var (
	Params = GetSettings()
	Users  = newClientCache()
)

func GetSettings() *SettingsStruct {
	var (
		jsonFile *os.File
		err      error
		byteJSON []byte
		result   = SettingsStruct{}
		errLog   = logs.ErrLog
	)

	//DEFAULT
	//result.mu = sync.RWMutex{}
	result.ServerPGSQL = "localhost"
	result.PortPGSQL = "5432"
	result.HTTPServer = "http://localhost"
	result.HttpPort = "8080"
	result.ConnLimit = 100
	result.UserTtlMin = 15
	result.Tables.Collection = fieldMap{}
	result.Tables.Details = fieldMap{}

	jsonFilePath := filepath.Join(errLog.AbsPath, "settings.json")
	jsonFile, err = os.Open(jsonFilePath)
	if err != nil {
		if jsn, bErr := json.Marshal(&result); bErr != nil {
			errLog.Fatal(bErr, "\nGetSettings(errLog *logging.LogStruct)",
				"\nif b,bErr := json.Marshal(settingsStruct{}); bErr != nil")
		} else {
			var prettyJSON bytes.Buffer
			if bErr = json.Indent(&prettyJSON, jsn, "", "\t"); bErr != nil {
				errLog.Fatal(bErr, "\nGetSettings(errLog *logging.LogStruct)",
					"\nif bErr = json.Indent(&prettyJSON, jsn, \"\", \"\t\"); bErr != nil")
			}
			if bErr = os.WriteFile(jsonFilePath, prettyJSON.Bytes(), 0644); bErr != nil {
				errLog.Fatal(bErr, "\nGetSettings(errLog *logging.LogStruct)",
					"\nif bErr = os.WriteFile(jsonDir, prettyJSON.Bytes(), 0644); bErr != nil")
			}
		}
		return &result
	}

	byteJSON, err = io.ReadAll(jsonFile)
	if err != nil {
		errLog.Write(err, "\nGetSettings(errLog *logging.LogStruct)",
			"\nSON, err = io.ReadAll(jsonFile)")
		return &result
	}

	err = jsonFile.Close()
	if err != nil {
		errLog.Write(err)
	}

	if !json.Valid(byteJSON) {
		errLog.Write(err, "invalid JSON string: ", string(byteJSON))
		return &result
	}

	err = json.Unmarshal(byteJSON, &result)
	if err != nil {
		errLog.Write(err, "Unmarshal error: ", string(byteJSON))
		return &result
	}

	return &result
}

func UpdateSettings(b io.ReadCloser) error {
	jsonBytes, err := io.ReadAll(b)
	if err != nil {
		return err
	}
	defer func() {
		if err = b.Close(); err != nil {
			logs.ErrLog.Write(err)
		}
	}()

	Params.MU.Lock()
	if err = json.Unmarshal(jsonBytes, &Params); err != nil {
		return err
	}
	Params.MU.Unlock()

	jsonFilePath := filepath.Join(logs.ErrLog.AbsPath, "settings.json")

	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, jsonBytes, "", "\t"); err != nil {
		logs.ErrLog.Write(err, "\nGetSettings(errLog *logging.LogStruct)",
			"\nif bErr = json.Indent(&prettyJSON, jsn, \"\", \"\t\"); bErr != nil")
		return err
	}
	if err = os.WriteFile(jsonFilePath, prettyJSON.Bytes(), 0644); err != nil {
		logs.ErrLog.Write(err, "\nGetSettings(errLog *logging.LogStruct)",
			"\nif bErr = os.WriteFile(jsonDir, prettyJSON.Bytes(), 0644); bErr != nil")
		return err
	}

	return nil
}

func IPFormRestriction(ip string) (int, error) {
	Params.MU.RLock()
	var requestBlockMin = Params.RequestBlockMin
	Params.MU.RUnlock()

	db := sqLite.GetDB()

	if requestBlockMin == 0 {
		return http.StatusOK, nil
	}

	query := "SELECT reqTime FROM reqFormIP WHERE ip = ? LIMIT 1"
	row := db.QueryRow(query, ip)

	var reqTime int64
	err := row.Scan(&reqTime)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logs.ErrLog.Write(err)
		return http.StatusInternalServerError, err
	}

	currentTime := time.Now().UTC().Unix()
	if tLeft := (currentTime - reqTime) / 60; tLeft < requestBlockMin {
		var duration string
		secLeft := (reqTime + requestBlockMin*60) - currentTime
		switch {
		case secLeft < 60:
			duration = fmt.Sprintf("%d seconds", secLeft)
		case secLeft < 3600:
			minutes := secLeft / 60
			duration = fmt.Sprintf("%d minutes", minutes)
		case secLeft < 86400:
			hours := secLeft / 3600
			duration = fmt.Sprintf("%d hours", hours)
		default:
			days := secLeft / 86400
			duration = fmt.Sprintf("%d days", days)
		}
		err = errors.New(fmt.Sprintf("Превышено максимальное количество запросов по IP: %s.\n"+
			"The maximum number of requests has been exceeded by IP: %s. "+
			"Please, repeat the request in %s", ip, ip, duration))

		logs.InfoLog.Write(err)
		return http.StatusTooManyRequests, err
	}

	query = "INSERT OR REPLACE INTO reqFormIP (ip, reqTime) VALUES(?,?)"
	if err = ExecuteQuery(query, ip, currentTime); err != nil {
		logs.ErrLog.Write(err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GenerateToken(length int) (string, error) {
	rndBytes := make([]byte, length)
	if _, err := rand.Read(rndBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(rndBytes), nil
}

func HashPassword(password string) (string, error) {
	hash := sha256.New()
	if _, err := io.WriteString(hash, password); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ExecuteQuery(query string, params ...any) error {
	db := sqLite.GetDB()
	stmtUsers, errUsers := db.Prepare(query)
	if errUsers != nil {
		return errUsers
	}
	defer func() {
		if err := stmtUsers.Close(); err != nil {
			logs.ErrLog.Write(err)
		}
	}()
	_, err := stmtUsers.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

func sendEMail(eMailTo, subject, message string) error {
	Params.MU.RLock()
	email := Params.EmailTech.Email
	pass := Params.EmailTech.Password
	smtpServ := Params.EmailTech.SMTPServer
	smtpPort := Params.EmailTech.SMTPPort
	Params.MU.RUnlock()

	auth := smtp.PlainAuth(
		"",
		email,    // Email
		pass,     // Password
		smtpServ, // SMTP server
	)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\nContent-Transfer-Encoding: 8bit\n\n"
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		mime+
		"\r\n"+
		"%s\r\n",
		eMailTo, subject, message))

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpServ, smtpPort), // SMTP server address with port
		auth,
		email, // Sender email
		strings.Split(eMailTo, ";"),
		msg,
	); err != nil {
		return err
	}

	return nil
}

func SendEmailRequest(company, email, applicant string) error {
	token, err := GenerateToken(16)
	if err != nil {
		return err
	}

	Params.MU.RLock()
	httpServer := Params.HTTPServer
	urlPort := Params.UrlPort
	emailTo := Params.EmailTech.EmailTo
	Params.MU.RUnlock()

	subject := fmt.Sprintf("API запрос от %s", company)
	msg := fmt.Sprintf("Запрос на доступ API от\n\n"+
		"Организация (Company): %s\n"+
		"E-Mail: %s\n"+
		"Заявитель (Applicant): %s\n\n"+
		"Для предоставления доступа перейдите по ссылке:\n\n"+
		fmt.Sprintf("%s%s%s/submit?id=%s", httpServer, urlPort, Params.HttpRedirectPrefix, token),
		company, email, applicant)

	if err = sendEMail(emailTo, subject, msg); err != nil {
		return err
	}

	query := "INSERT INTO requests (id, company, email,applicant,active,requestTime) VALUES(?,?,?,?,false,?)"
	if err = ExecuteQuery(query, token, company, email, applicant, time.Now().UTC().Unix()); err != nil {
		logs.ErrLog.Write(err)
		return err
	}

	return nil
}

func SubmitApplication(w http.ResponseWriter, id string) error {
	db := sqLite.GetDB()
	Params.MU.RLock()
	urlDir := fmt.Sprintf("%s%s%s/files", Params.HTTPServer, Params.UrlPort, Params.HttpRedirectPrefix)
	Params.MU.RUnlock()

	query := "SELECT email FROM requests WHERE id = ? and active = false LIMIT 1"
	row := db.QueryRow(query, id)

	var email string
	err := row.Scan(&email)
	if err != nil {
		return err
	}

	query = "SELECT COALESCE(max(id), 0) + 1 FROM users LIMIT 1"
	row = db.QueryRow(query, id)

	var userId int
	err = row.Scan(&userId)
	if err != nil {
		return err
	}

	username := fmt.Sprintf("user%s", strconv.Itoa(userId))
	password, err := GenerateToken(10)
	if err != nil {
		return err
	}

	hashPass, errHash := HashPassword(password)
	if errHash != nil {
		return err
	}

	query = "INSERT INTO users (user, password, active, requestID) VALUES (?,?,true,?)"
	if err = ExecuteQuery(query, username, hashPass, id); err != nil {
		return err
	}

	query = "UPDATE requests SET active = true, submitTime = ? WHERE id = ?"
	if err = ExecuteQuery(query, time.Now().UTC().Unix(), id); err != nil {
		return err
	}

	apiDesc := "https://github.com/Hardmun/rms/blob/main/install/api_ru.md"
	intro := "We are pleased to inform you that your request for access to our REST API " +
		"has been successfully granted."
	manuals := fmt.Sprintf("REST API documentation:\n\n%s\n\n%s\n%s", apiDesc,
		fmt.Sprintf("%s/manual_v1_en.docx", urlDir), fmt.Sprintf("%s/manual_v1_ru.docx", urlDir))
	loginPass := fmt.Sprintf("Username: %s\nPassword: %s", username, password)
	err = sendEMail(email, "The request has been approved.", fmt.Sprintf("%s\n\n%s\n\n\n%s",
		intro, loginPass, manuals))
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintf(w, fmt.Sprintf("The login credentials have been sent to %s\n\n%s",
		email, loginPass)); err != nil {
		return err
	}

	return nil
}

func dbPass(usr string) (pass string, err error) {
	db := sqLite.GetDB()

	query := "SELECT password FROM users WHERE user = ? and active = true LIMIT 1"
	row := db.QueryRow(query, usr)

	err = row.Scan(&pass)
	if err != nil {
		return
	}

	return
}

func AuthAndRateLimiter(next http.Handler, w http.ResponseWriter, r *http.Request) {
	var usr user
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "invalid authorization", http.StatusNonAuthoritativeInfo)
		return
	}
	if username == "" || password == "" {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	usr, ok = Users.getUser(username)
	if !ok {
		dbPassword, err := dbPass(username)
		if err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		usr = Users.addUser(username, dbPassword)
	}

	hashPass, err := HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNonAuthoritativeInfo)
		return
	}
	if !(hashPass == usr.hashPass) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	select {
	case usr.rateLimit <- struct{}{}:
		defer func() { <-usr.rateLimit }()
		next.ServeHTTP(w, r)
	default:
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
	}
}

// src -file, dst - directory
func copyFile(src, dst string) error {
	var (
		err     error
		srcFile *os.File
		dstFile *os.File
	)

	if _, err = os.Stat(dst); os.IsNotExist(err) {
		if err = os.MkdirAll(dst, os.ModePerm); err != nil {
			return err
		}
	}

	fileName := filepath.Base(src)
	srcFile, err = os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err = srcFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	dstFilePath := filepath.Join(dst, fileName)
	if _, err = os.Stat(dstFilePath); os.IsNotExist(err) {
		dstFile, err = os.Create(dstFilePath)
		if err != nil {
			return err
		}
	} else {
		dstFile, err = os.OpenFile(dstFilePath, os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	defer func() {
		if err = dstFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err = dstFile.Sync(); err != nil {
		return err
	}

	if err = exec.Command("chmod", "+x", dstFilePath).Run(); err != nil {
		return err
	}

	return nil
}

func LinuxService(src string) error {
	progDir := "/usr/local/rms"
	serviceFilePath := "/etc/systemd/system/rms.service"

	if err := copyFile(src, progDir); err != nil {
		return err
	}

	srvFile, errSrv := embedFiles.ReadFile("rms@.service")
	if errSrv != nil {
		return errSrv
	}

	if err := os.WriteFile(serviceFilePath, srvFile, 0644); err != nil {
		return err
	}

	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}

	if err := exec.Command("systemctl", "start", "rms.service").Run(); err != nil {
		return err
	}

	if err := exec.Command("systemctl", "enable", "rms.service").Run(); err != nil {
		return err
	}

	return nil
}

func IsIpPrivate(remoteIP string) error {
	ip, _, err := net.SplitHostPort(remoteIP)
	if err != nil {
		return err
	}

	if ip == "::1" {
		return nil
	}

	remote := net.ParseIP(ip)
	if !remote.IsPrivate() {
		return errors.New("remote address is not private")
	}

	return nil
}

func fileWithExt(ext string) (string, error) {
	sslPath := filepath.Join(logs.ErrLog.AbsPath, "ssl")
	var result string

	err := filepath.Walk(sslPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			result = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("file not found")
	}

	return result, nil
}

func IsSSL() (string, string, bool) {
	var (
		crtFile string
		keyFile string
		err     error
	)

	crtFile, err = fileWithExt(".crt")
	if err != nil {
		return "", "", false
	}
	keyFile, err = fileWithExt(".key")
	if err != nil {
		return "", "", false
	}

	return crtFile, keyFile, true
}
