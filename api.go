package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	account "rttys/account"
	"rttys/cache"
	"rttys/config"
	tenant "rttys/tenant"
	"rttys/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    int8   `json:"admin"`
	Tenant   string `json:"tenant"`
}

type Telant struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

type SBatchTenants struct {
	Tenants []string `json:"tenants"`
}

var httpSessions *cache.Cache

//go:embed ui/dist
var staticFs embed.FS

func allowOrigin(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
}

func httpLogin(cfg *config.Config, creds *credentials) bool {
	if creds.Username == "" || creds.Password == "" {
		return false
	}

	db, err := instanceDB(cfg.DB)
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}
	defer db.Close()

	cnt := 0

	db.QueryRow("SELECT COUNT(*) FROM account WHERE username = ? AND password = ?", creds.Username, creds.Password).Scan(&cnt)

	return cnt != 0
}

func getUserTelnet(cfg *config.Config, creds *credentials) string {
	if creds.Username == "" || creds.Password == "" {
		return ""
	}

	db, err := instanceDB(cfg.DB)
	if err != nil {
		log.Error().Msg(err.Error())
		// return false
	}
	defer db.Close()

	cnt := ""

	db.QueryRow("SELECT tenant FROM account WHERE username = ? AND password = ?", creds.Username, creds.Password).Scan(&cnt)

	return cnt

}

func authorizedDev(devid string, cfg *config.Config) bool {
	if cfg.WhiteList == nil {
		return true
	}

	_, ok := cfg.WhiteList[devid]
	return ok
}

func isLocalRequest(c *gin.Context) bool {
	addr, _ := net.ResolveTCPAddr("tcp", c.Request.RemoteAddr)
	return addr.IP.IsLoopback()
}

func httpAuth(cfg *config.Config, c *gin.Context) bool {
	if !cfg.LocalAuth && isLocalRequest(c) {
		return true
	}

	cookie, err := c.Cookie("sid")
	if err != nil || !httpSessions.Have(cookie) {
		return false
	}

	httpSessions.Active(cookie, 0)

	return true
}

func isAdminUsername(cfg *config.Config, username string) bool {
	if username == "" {
		return false
	}

	db, err := instanceDB(cfg.DB)
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}
	defer db.Close()

	isAdmin := false

	if db.QueryRow("SELECT admin FROM account WHERE username = ?", username).Scan(&isAdmin) == sql.ErrNoRows {
		return false
	}

	return isAdmin
}

func getLoginUsername(c *gin.Context) string {
	cookie, err := c.Cookie("sid")
	if err != nil {
		return ""
	}

	username, ok := httpSessions.Get(cookie)
	if ok {
		return username.(string)
	}

	return ""
}

func apiStart(br *broker) {
	cfg := br.cfg

	httpSessions = cache.New(30*time.Minute, 5*time.Second)

	gin.SetMode(gin.ReleaseMode)
	gormDB, err := initGORM(cfg.DB)
	if err != nil {
		log.Err(err)
	}
	tenantRepo := tenant.NewTenantRepository(gormDB)
	accountRepo := account.NewAccountRepository(gormDB)
	r := gin.New()
	r.Use(gin.Recovery())
	// authz.NewAuthorizer(e)
	// r.Use(authz.NewAuthorizer(e))

	authorized := r.Group("/", func(c *gin.Context) {
		devid := c.Param("devid")
		if devid != "" && authorizedDev(devid, cfg) {
			return
		}

		if !httpAuth(cfg, c) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	// authorized.GET("/roles", func(c *gin.Context) {
	// 	roles := e.GetAllRoles()
	// 	c.JSON(http.StatusOK, gin.H{"roles": roles})
	// })

	authorized.GET("/tenants", func(c *gin.Context) {
		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()
		telants := tenantRepo.Query(tenant.SPTelant{})
		c.JSON(http.StatusOK, gin.H{"telant": telants})
	})

	authorized.POST("/tenants", func(c *gin.Context) {

		tens := tenant.STenant{}
		if err := c.BindJSON(&tens); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		id, err := tenantRepo.Create(tens)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	authorized.PUT("/tenants", func(c *gin.Context) {

		tens := tenant.STenant{}
		if err := c.BindJSON(&tens); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		where := tenant.SPTelant{}
		where.Name = tens.Name
		id, err := tenantRepo.Update(tens, where)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	authorized.DELETE("/tenants/:name", func(c *gin.Context) {
		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()
		name := c.Param("name")

		id, err := tenantRepo.Delete(tenant.STenant{
			Name: name,
		})
		if err != nil {
			log.Err(err)
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	authorized.GET("/tenants/users", func(c *gin.Context) {
		tenant := c.GetHeader("tenant")
		account := accountRepo.Query(account.SPAccount{
			SAccount: account.SAccount{
				Tenant: tenant,
			},
		})

		c.JSON(http.StatusOK, gin.H{"users": account})
	})

	authorized.POST("/tenants/users", func(c *gin.Context) {
		loginUser := account.SAccount{}
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		loginUser.Token = utils.GenUniqueID(loginUser.Password)
		id, err := accountRepo.Create(loginUser)
		if err != nil {
			log.Err(err)
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	authorized.PUT("/tenants/users", func(c *gin.Context) {
		loginUser := account.SAccount{}
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		spa := account.SPAccount{}
		spa.Username = loginUser.Username
		spa.Tenant = loginUser.Tenant
		id, err := accountRepo.Update(loginUser, spa)
		if err != nil {
			log.Err(err)
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	authorized.DELETE("/tenants/users/:username", func(c *gin.Context) {
		// tenant := c.GetHeader("tenant")
		username := c.Param("username")

		id, err := accountRepo.Delete(account.SAccount{
			Username: username,
		})
		if err != nil {
			// log.Err(err)
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		// id, err := result.LastInsertId()
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	authorized.GET("/fontsize", func(c *gin.Context) {
		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		value := "16"

		db.QueryRow("SELECT value FROM config WHERE name = 'FontSize'").Scan(&value)

		FontSize, _ := strconv.Atoi(value)

		c.JSON(http.StatusOK, gin.H{"size": FontSize})
	})

	authorized.POST("/fontsize", func(c *gin.Context) {
		data := make(map[string]int)

		err := c.BindJSON(&data)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		size, ok := data["size"]
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		if size < 12 {
			size = 12
		}

		_, err = db.Exec("DELETE FROM config WHERE name = 'FontSize'")
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO config values('FontSize',?)", fmt.Sprintf("%d", size))
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	authorized.GET("/connect/:devid", func(c *gin.Context) {
		if c.GetHeader("Upgrade") != "websocket" {
			c.Redirect(http.StatusFound, "/rtty/"+c.Param("devid"))
			return
		}
		serveUser(br, c)
	})

	authorized.GET("/devs", func(c *gin.Context) {
		type DeviceInfo struct {
			ID          string `json:"id"`
			Connected   uint32 `json:"connected"`
			Uptime      uint32 `json:"uptime"`
			Description string `json:"description"`
			Bound       bool   `json:"bound"`
			Online      bool   `json:"online"`
			Proto       uint8  `json:"proto"`
		}
		tenant := c.GetHeader("Tenant")
		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		sql := "SELECT id, description, username,tenant FROM device where 1=1"

		if cfg.LocalAuth || !isLocalRequest(c) {
			username := getLoginUsername(c)
			if username == "" {
				c.Status(http.StatusUnauthorized)
				return
			}

			if !isAdminUsername(cfg, username) {
				sql += fmt.Sprintf(" and username = '%s'", username)
			} else {

			}
		}

		sql += fmt.Sprintf(" and tenant = '%s'", tenant)
		devs := make([]DeviceInfo, 0)

		rows, err := db.Query(sql)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			id := ""
			desc := ""
			username := ""
			tenant := ""

			err := rows.Scan(&id, &desc, &username, &tenant)
			if err != nil {
				log.Error().Msg(err.Error())
				break
			}

			di := DeviceInfo{
				ID:          id,
				Description: desc,
				Bound:       username != "",
			}

			if dev, ok := br.devices[id]; ok {
				dev := dev.(*device)
				di.Connected = uint32(time.Now().Unix() - dev.timestamp)
				di.Uptime = dev.uptime
				di.Online = true
				di.Proto = dev.proto
			}

			devs = append(devs, di)
		}

		allowOrigin(c.Writer)

		c.JSON(http.StatusOK, devs)
	})

	authorized.POST("/cmd/:devid", func(c *gin.Context) {
		allowOrigin(c.Writer)

		handleCmdReq(br, c)
	})

	r.Any("/web/:devid/:proto/:addr/*path", func(c *gin.Context) {
		httpProxyRedirect(br, c)
	})

	r.GET("/authorized/:devid", func(c *gin.Context) {
		authorized := authorizedDev(c.Param("devid"), cfg) || httpAuth(cfg, c)
		c.JSON(http.StatusOK, gin.H{
			"authorized": authorized,
		})
	})

	r.POST("/signin", func(c *gin.Context) {
		var creds credentials
		err := c.BindJSON(&creds)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if httpLogin(cfg, &creds) {
			sid := utils.GenUniqueID("http")
			httpSessions.Set(sid, creds.Username, 0)
			tenant := getUserTelnet(cfg, &creds)
			c.SetCookie("sid", sid, 0, "", "", false, true)

			c.JSON(http.StatusOK, gin.H{
				"sid":      sid,
				"username": creds.Username,
				"tenant":   tenant,
			})
			return
		}

		c.Status(http.StatusForbidden)
	})

	r.GET("/alive", func(c *gin.Context) {
		if !httpAuth(cfg, c) {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Status(http.StatusOK)
		}
	})

	r.GET("/signout", func(c *gin.Context) {
		cookie, err := c.Cookie("sid")
		if err != nil || !httpSessions.Have(cookie) {
			return
		}

		httpSessions.Del(cookie)

		c.Status(http.StatusOK)
	})

	r.POST("/signup", func(c *gin.Context) {
		var creds credentials

		err := c.BindJSON(&creds)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		isAdmin := 0

		cnt := 0
		db.QueryRow("SELECT COUNT(*) FROM account").Scan(&cnt)
		if cnt == 0 {
			isAdmin = 1
		}

		db.QueryRow("SELECT COUNT(*) FROM account WHERE username = ?", creds.Username).Scan(&cnt)
		if cnt > 0 {
			c.Status(http.StatusForbidden)
			return
		}

		_, err = db.Exec("INSERT INTO account values(?,?,?)", creds.Username, creds.Password, isAdmin)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	r.GET("/isadmin", func(c *gin.Context) {
		isAdmin := true

		if cfg.LocalAuth || !isLocalRequest(c) {
			isAdmin = isAdminUsername(cfg, getLoginUsername(c))
		}

		c.JSON(http.StatusOK, gin.H{"admin": isAdmin})
	})

	r.GET("/users", func(c *gin.Context) {
		loginUsername := getLoginUsername(c)
		isAdmin := isAdminUsername(cfg, loginUsername)
		telant := c.GetHeader("Telant")
		username := c.Query("username")
		var sqlBuilder strings.Builder
		if cfg.LocalAuth || !isLocalRequest(c) {
			if !isAdmin {
				c.Status(http.StatusUnauthorized)
				return
			}
		}

		users := []string{}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		defer db.Close()
		// String
		sqlBuilder.WriteString("SELECT username FROM account where 1=1 ")
		if telant != "" {
			sqlBuilder.WriteString(fmt.Sprintf("and telant=%s ", telant))
		}
		if len(username) > 0 {
			sqlBuilder.WriteString(fmt.Sprintf("and username='%s' ", username))
		}
		rows, err := db.Query(sqlBuilder.String())
		if err != nil {
			log.Error().Msg(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			username := ""
			err := rows.Scan(&username)
			if err != nil {
				log.Error().Msg(err.Error())
				break
			}

			if isAdmin && username == loginUsername {
				continue
			}

			users = append(users, username)
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/bind", func(c *gin.Context) {
		if cfg.LocalAuth || !isLocalRequest(c) {
			username := getLoginUsername(c)
			if !isAdminUsername(cfg, username) {
				c.Status(http.StatusUnauthorized)
				return
			}
		}

		type binddata struct {
			Username string   `json:"username"`
			Devices  []string `json:"devices"`
		}

		data := binddata{}

		err := c.BindJSON(&data)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		defer db.Close()

		isAdmin := false

		if db.QueryRow("SELECT admin FROM account WHERE username = ?", data.Username).Scan(&isAdmin) == sql.ErrNoRows || isAdmin {
			c.Status(http.StatusOK)
			return
		}

		for _, devid := range data.Devices {
			db.Exec("UPDATE device SET username = ? WHERE id = ?", data.Username, devid)
		}

		c.Status(http.StatusOK)
	})

	r.POST("/unbind", func(c *gin.Context) {
		if cfg.LocalAuth || !isLocalRequest(c) {
			username := getLoginUsername(c)
			if !isAdminUsername(cfg, username) {
				c.Status(http.StatusUnauthorized)
				return
			}
		}

		type binddata struct {
			Devices []string `json:"devices"`
		}

		data := binddata{}

		err := c.BindJSON(&data)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		defer db.Close()

		for _, devid := range data.Devices {
			db.Exec("UPDATE device SET username = '' WHERE id = ?", devid)
		}

		c.Status(http.StatusOK)
	})

	r.POST("/delete", func(c *gin.Context) {
		type deldata struct {
			Devices []string `json:"devices"`
		}

		data := deldata{}

		err := c.BindJSON(&data)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		db, err := instanceDB(cfg.DB)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		defer db.Close()

		username := ""
		if cfg.LocalAuth || !isLocalRequest(c) {
			username = getLoginUsername(c)
			if isAdminUsername(cfg, username) {
				username = ""
			}
		}

		for _, devid := range data.Devices {
			if _, ok := br.devices[devid]; !ok {
				sql := fmt.Sprintf("DELETE FROM device WHERE id = '%s'", devid)

				if username != "" {
					sql += fmt.Sprintf(" AND username = '%s'", username)
				}

				db.Exec(sql)
			}
		}

		c.Status(http.StatusOK)
	})

	r.GET("/file/:sid", func(c *gin.Context) {
		sid := c.Param("sid")
		if fp, ok := br.fileProxy.Load(sid); ok {
			fp := fp.(*fileProxy)
			s := br.sessions[sid]
			fp.Ack(s.dev, sid)

			defer func() {
				if err := recover(); err != nil {
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								fp.reader.Close()
							}
						}
					}
				}
			}()

			c.DataFromReader(http.StatusOK, -1, "application/octet-stream", fp.reader, nil)
			br.fileProxy.Delete(sid)
		}
	})

	r.NoRoute(func(c *gin.Context) {
		fs, _ := fs.Sub(staticFs, "ui/dist")

		path := c.Request.URL.Path

		if path != "/" {
			f, err := fs.Open(path[1:])
			if err != nil {
				c.Request.URL.Path = "/"
				r.HandleContext(c)
				return
			}

			if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
				if strings.HasSuffix(path, "css") || strings.HasSuffix(path, "js") {
					magic := make([]byte, 2)
					f.Read(magic)
					if magic[0] == 0x1f && magic[1] == 0x8b {
						c.Writer.Header().Set("Content-Encoding", "gzip")
					}
				}
			}

			f.Close()
		}

		http.FileServer(http.FS(fs)).ServeHTTP(c.Writer, c.Request)
	})

	go func() {
		var err error

		if cfg.SslCert != "" && cfg.SslKey != "" {
			log.Info().Msgf("Listen user on: %s SSL on", cfg.AddrUser)
			err = r.RunTLS(cfg.AddrUser, cfg.SslCert, cfg.SslKey)
		} else {
			log.Info().Msgf("Listen user on: %s SSL off", cfg.AddrUser)
			err = r.Run(cfg.AddrUser)
		}

		log.Fatal().Err(err)
	}()
}
