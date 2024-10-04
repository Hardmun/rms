package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"rms/logs"
	pgsql "rms/postgres"
	"rms/sqLite"
	"rms/utils"
	"runtime"
	"strings"
)

//go:embed html/auth.html
var embedFiles embed.FS

func applicationForm(w http.ResponseWriter) {
	contFile, errContFile := embedFiles.ReadFile("html/auth.html")
	if errContFile != nil {
		logs.ErrLog.ErrorHTTP(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderTmpl := strings.ReplaceAll(string(contFile), "%prefix", utils.Params.HttpRedirectPrefix)

	htmlTmpl, err := template.New("").Parse(renderTmpl)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = htmlTmpl.Execute(w, nil)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func requestApplicationHandler(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("source")
	if source != "button" || r.Method == http.MethodGet {
		applicationForm(w)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, "Invalid remote address", http.StatusInternalServerError)
		logs.ErrLog.Write(err)
		http.Error(w, "Invalid remote address", http.StatusInternalServerError)
		return
	}

	if httpStatus, errHTTP := utils.IPFormRestriction(ip); errHTTP != nil {
		http.Error(w, errHTTP.Error(), httpStatus)
		return
	}

	company := r.FormValue("company")
	email := r.FormValue("email")
	applicant := r.FormValue("applicant")

	err = utils.SendEmailRequest(company, email, applicant)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = fmt.Fprintf(w, "Ваши данные отправлены на рассмотрение."+
		"\nYour data has been submitted for review."+
		"\n Company: %v\n E-mail: %v\n Applicant: %v", company, email, applicant); err != nil {
		logs.ErrLog.Write(err)
		return
	}
}

func submitApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Не найден ID запроса.\n"+
			"The ID is required.", http.StatusBadRequest)
		return
	}

	if err := utils.SubmitApplication(w, id); err != nil {
		http.Error(w, fmt.Sprintf("Запрос по данной сессии рассмотрен либо истек.\n"+
			"The request for this session has either been processed or expired.\n\n"+
			"Error for support: %s", err.Error()),
			http.StatusBadRequest)
		if err.Error() != "sql: no rows in result set" {
			logs.ErrLog.Write(err)
		}
		return
	}
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	if err := utils.IsIpPrivate(r.RemoteAddr); err != nil {
		logs.InfoLog.Write(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(utils.Params); err != nil {
			logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		if err := utils.UpdateSettings(r.Body); err != nil {
			logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		}

		pgsql.UpdateSettings()

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func collectionBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := pgsql.GetCollectionRemaining(r)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func detailBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := pgsql.GetDetailRemaining(r)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func collectionImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := pgsql.GetCollectionImage(r)
	if err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		logs.ErrLog.ErrorHTTP(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func handleMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == wrapPattern("/request") || r.URL.Path == wrapPattern("/submit") ||
			r.URL.Path == "/settings" || r.URL.Path == wrapPattern("/files/manual_v1_en.docx") ||
			r.URL.Path == wrapPattern("/files/manual_v1_ru.docx") {
			next.ServeHTTP(w, r)
			return
		}

		utils.AuthAndRateLimiter(next, w, r)
	})
}

func wrapPattern(pattern string) string {
	return fmt.Sprintf("%s%s", utils.Params.HttpRedirectPrefix, pattern)
}

func main() {
	args := os.Args
	if len(args) == 2 && args[1] == "--install" {
		if runtime.GOOS != "linux" {
			fmt.Println("service can be installed only for Linux.")
			return
		}
		if err := utils.LinuxService(args[0]); err != nil {
			fmt.Println(err)
		}
		return
	}

	defer func() {
		err := sqLite.DB.Close()
		if err != nil {
			logs.ErrLog.Write(fmt.Errorf("failed to close SQLite DB: %w", err))
		}
		err = pgsql.DbPG.Close()
		if err != nil {
			logs.ErrLog.Write(fmt.Errorf("failed to close Postgres DB: %w", err))
		}
	}()

	params := utils.Params
	utils.Users.StartCleanupJob()

	mux := http.NewServeMux()

	mux.HandleFunc(wrapPattern("/request"), requestApplicationHandler)
	mux.HandleFunc(wrapPattern("/submit"), submitApplicationHandler)

	mux.HandleFunc("/settings", handleSettings)

	mux.HandleFunc(wrapPattern("/api/v1/collection"), collectionBalance)
	mux.HandleFunc(wrapPattern("/api/v1/details"), detailBalance)
	mux.HandleFunc(wrapPattern("/api/v1/images"), collectionImages)

	fs := http.FileServer(http.Dir(params.FileDirectory))
	mux.Handle(wrapPattern("/files/"), http.StripPrefix(wrapPattern("/files/"), fs))

	wrappedMux := handleMux(mux)
	exitOption := 0

	if crt, key, isSSL := utils.IsSSL(); isSSL {
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s", params.HttpPort), crt, key, wrappedMux); err != nil {
			logs.ErrLog.Write(err)
			exitOption = 1
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf("%s", params.HttpPort), wrappedMux); err != nil {
			logs.ErrLog.Write(err)
			exitOption = 1
		}
	}

	if err := logs.ErrLog.CloseLog(); err != nil {
		exitOption = 1
	}
	os.Exit(exitOption)
}
