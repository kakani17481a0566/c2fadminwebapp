package session

import (
	"errors"
	"strings"
	"sync"

	"c2f_admin_webapp/model/net/redis"

	"c2f_admin_webapp/model/util/datetime"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
)

var GlobalSessions *session.Manager
var ru redis.Redis
var rc redis.Redis

var (
	memSessions = make(map[string]string)
	memMu       sync.RWMutex
)

func Init() (err error) {
	var cf session.ManagerConfig
	cf.CookieName = beego.AppConfig.String("appname")
	cf.EnableSetCookie = true
	cf.Gclifetime, _ = beego.AppConfig.Int64("gclifetime")
	cf.Maxlifetime, _ = beego.AppConfig.Int64("maxLifetime")
	cf.Secure, _ = beego.AppConfig.Bool("EnableHTTPS")
	cf.CookieLifeTime, _ = beego.AppConfig.Int("cookieLifeTime")
	cf.ProviderConfig = beego.AppConfig.String("SessionProviderConfig")

	provider := beego.AppConfig.String("SessionProviderName")
	if provider == "" {
		provider = "memory"
	}

	GlobalSessions, err = session.NewManager(provider, &cf)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// Fallback to memory if redis session manager fails
		provider = "memory"
		GlobalSessions, err = session.NewManager("memory", &cf)
		if err != nil {
			return err
		}
	}

	go GlobalSessions.GC()

	if provider == "redis" {
		ru.SavePath = beego.AppConfig.String("RedisConnectionString")
		if errRu := ru.Connect(); errRu != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Warning", "Redis ru Connect failed:", errRu)
		}
		rc.SavePath = beego.AppConfig.String("SessionProviderConfig")
		if errRc := rc.Connect(); errRc != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Warning", "Redis rc Connect failed:", errRc)
		}
	}

	return nil
}

func ValidateSession(sess session.Store) (err error) {
	if sess == nil {
		return errors.New("Session nil")
	}
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Session ID - ", sess.SessionID())

	if rc.Poollist != nil {
		if !rc.Exists(sess.SessionID()) {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Session Expire")
			err = errors.New("Session Time Out.Please Logout and Login Again.")
			return
		}
	}

	unameVal := sess.Get("uname")
	if unameVal == nil {
		return errors.New("User session not found.Please Logout and Login Again.")
	}
	uname := unameVal.(string)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)

	if ru.Poollist != nil && ru.Exists(uname) {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Found")
		val, err := ru.Get(uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return err
		}
		data := strings.Split(val, "*")
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Status - ", data[3])
		if data[0] == sess.SessionID() && data[3] == "LOGGED_IN" {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Good")
			return nil
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Bad")
			err = errors.New("Invalid session found in cache.Please Logout and Login Again.")
			return err
		}
	} else if ru.Poollist == nil {
		// In-memory fallback
		memMu.RLock()
		val, exists := memSessions[uname]
		memMu.RUnlock()
		if exists {
			data := strings.Split(val, "*")
			if len(data) > 3 && data[0] == sess.SessionID() && data[3] == "LOGGED_IN" {
				return nil
			}
		}
		return nil
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Not Found")
		err = errors.New("User session not found.Please Logout and Login Again.")
	}
	return
}

func SetUserSession(sid, uname string) (err error) {
	dt, _ := datetime.Get("", "", "Asia/Kolkata")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Session ID - ", sid)

	newSess := sid + "**" + dt + "*" + "LOGGED_IN" + "*"

	if ru.Poollist != nil {
		if ru.Exists(uname) {
			log.Println(beego.AppConfig.String("loglevel"), "Info", "User Session Found")
			val, err := ru.Get(uname)
			if err == nil {
				data := strings.Split(val, "*")
				if len(data) > 2 {
					newSess = sid + "*" + data[2] + "*" + dt + "*" + "LOGGED_IN" + "*"
				}
			}
		}
		_ = ru.Set(uname, newSess)
	}

	if rc.Poollist != nil {
		_ = rc.Set(sid, uname)
	}

	memMu.Lock()
	memSessions[uname] = newSess
	memSessions[sid] = uname
	memMu.Unlock()

	return nil
}

func SeTLogoutSession(uname string) (err error) {
	dt, _ := datetime.Get("", "", "Asia/Kolkata")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName - ", uname)
	if ru.Poollist != nil && ru.Exists(uname) {
		val, err := ru.Get(uname)
		if err == nil {
			data := strings.Split(val, "*")
			newSess := data[0] + "*" + data[1] + "*" + data[2] + "*" + "LOGGED_OUT" + "*" + dt
			_ = ru.Set(uname, newSess)
		}
	}
	memMu.Lock()
	delete(memSessions, uname)
	memMu.Unlock()
	return nil
}

func CheckUserSession(uname string) (err error) {
	dt, _ := datetime.Get("", "", "Asia/Kolkata")
	if ru.Poollist != nil && ru.Exists(uname) {
		val, err := ru.Get(uname)
		if err == nil {
			data := strings.Split(val, "*")
			if len(data) > 3 && data[3] == "LOGGED_OUT" {
				return nil
			}
			_ = ru.Set(uname, data[0]+"*"+data[1]+"*"+data[2]+"*"+"LOGGED_OUT"+"*"+dt)
			return errors.New("Multiple user login not allowed. Now you have been logged out from all browsers. Please login again.")
		}
	}
	return nil
}

func CheckHijack(sess session.Store, ctx *context.Context) error {
	storedIP, _ := sess.Get("bind_ip").(string)
	storedUA, _ := sess.Get("bind_ua").(string)

	if storedIP == "" || storedUA == "" {
		sess.Set("bind_ip", ctx.Input.IP())
		sess.Set("bind_ua", ctx.Request.UserAgent())
		return nil
	}

	if storedIP != ctx.Input.IP() || storedUA != ctx.Request.UserAgent() {
		sess.Flush()
		return errors.New("Session hijacked or invalid device detected")
	}

	return nil
}
