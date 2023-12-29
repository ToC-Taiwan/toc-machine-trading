package usecase

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	"tmt/internal/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/repo"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type SystemUseCase struct {
	repo    SystemRepo
	smtpCfg config.SMTP

	activationCodeMap     map[string]time.Time
	activationCodeMapLock sync.RWMutex
}

func NewSystem() *SystemUseCase {
	cfg := config.Get()
	return &SystemUseCase{
		repo:              repo.NewSystemRepo(cfg.GetPostgresPool()),
		activationCodeMap: make(map[string]time.Time),
		smtpCfg:           config.Get().SMTP,
	}
}

func (uc *SystemUseCase) AddUser(ctx context.Context, t *entity.User) error {
	allUser, err := uc.repo.QueryAllUser(ctx)
	if err != nil {
		return err
	}

	if len(allUser) >= 10 {
		return errors.New("user limit exceeded")
	}

	for _, user := range allUser {
		if user.Username == t.Username {
			return errors.New("username already exists")
		}
		if user.Email == t.Email {
			return errors.New("email already exists")
		}
	}

	if err := uc.repo.InsertUser(ctx, t); err != nil {
		return err
	}

	return uc.SendOTP(ctx, t)
}

func (uc *SystemUseCase) Login(ctx context.Context, username, password string) (bool, error) {
	user, err := uc.repo.QueryUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("username not found")
	}
	if user.Password != password {
		return false, errors.New("password not match")
	}
	if !user.EmailVerified {
		return false, errors.New("email not verified")
	}
	if !user.Activated {
		return false, errors.New("user not activated")
	}
	return true, nil
}

func (uc *SystemUseCase) SendOTP(ctx context.Context, t *entity.User) error {
	if uc.smtpCfg.Host == "" || uc.smtpCfg.Port == 0 || uc.smtpCfg.Username == "" || uc.smtpCfg.Password == "" {
		return errors.New("smtp config not set")
	}

	activationCode := uuid.NewString()
	uc.addActivationCode(activationCode)

	m := gomail.NewMessage()
	m.SetHeader("From", uc.smtpCfg.Username)
	m.SetHeader("To", t.Email)
	m.SetHeader("Subject", "Please verify your email address")
	m.SetBody(
		"text/html",
		fmt.Sprintf("Please click the following link in 30 minutes to verify your email address: <a href='https://trader.tocraw.com/tmt/v1/user/verify/%s/%s'>Verify</a>", t.Username, activationCode),
	)

	d := gomail.NewDialer(uc.smtpCfg.Host, 587, uc.smtpCfg.Username, uc.smtpCfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (uc *SystemUseCase) ActivateUser(ctx context.Context, username string) error {
	user, err := uc.repo.QueryUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("username not found")
	}
	if user.Activated {
		return errors.New("user already activated")
	}
	if !user.EmailVerified {
		return errors.New("email not verified")
	}
	return uc.repo.ActivateUser(ctx, username)
}

func (uc *SystemUseCase) VerifyEmail(ctx context.Context, username, code string) error {
	defer uc.activationCodeMapLock.Unlock()
	uc.activationCodeMapLock.Lock()

	expire, ok := uc.activationCodeMap[code]
	if !ok {
		return errors.New("invalid activation code")
	}

	if time.Now().After(expire.Add(30 * time.Minute)) {
		return errors.New("activation code expired")
	}

	delete(uc.activationCodeMap, code)
	return uc.repo.EmailVerification(ctx, username)
}

func (uc *SystemUseCase) addActivationCode(code string) {
	defer uc.activationCodeMapLock.Unlock()
	uc.activationCodeMapLock.Lock()
	uc.activationCodeMap[code] = time.Now()
}
