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
	"tmt/pkg/eventbus"
	"tmt/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type SystemUseCase struct {
	repo    SystemRepo
	smtpCfg config.SMTP

	activationCodeMap     map[string]time.Time
	activationCodeMapLock sync.RWMutex

	logger *log.Log
	bus    *eventbus.Bus
}

func NewSystem() *SystemUseCase {
	cfg := config.Get()
	uc := &SystemUseCase{
		repo:              repo.NewSystemRepo(cfg.GetPostgresPool()),
		activationCodeMap: make(map[string]time.Time),
		smtpCfg:           config.Get().SMTP,
		logger:            log.Get(),
		bus:               eventbus.Get(),
	}

	uc.UpdateAuthTradeUser()
	return uc
}

func (uc *SystemUseCase) UpdateAuthTradeUser() {
	allUser, err := uc.repo.QueryAllUser(context.Background())
	if err != nil {
		uc.logger.Fatal(err)
	}

	authUserName := []string{}
	for _, user := range allUser {
		if user.AuthTrade {
			authUserName = append(authUserName, user.Username)
		}
	}

	uc.bus.PublishTopicEvent(topicUpdateAuthTradeUser, authUserName)
}

func (uc *SystemUseCase) AddUser(ctx context.Context, t *entity.User) error {
	allUser, err := uc.repo.QueryAllUser(ctx)
	if err != nil {
		return err
	}

	// if len(allUser) >= 10 {
	// 	return errors.New("user limit exceeded")
	// }

	for _, user := range allUser {
		if user.Username == t.Username {
			return ErrUsernameAlreadyExists
		}
		if user.Email == t.Email {
			return ErrEmailAlreadyExists
		}
	}

	t.Password, err = uc.EncryptPassword(ctx, t.Password)
	if err != nil {
		return err
	}
	if err := uc.repo.InsertUser(ctx, t); err != nil {
		return err
	}

	return uc.SendOTP(ctx, t)
}

func (uc *SystemUseCase) Login(ctx *gin.Context, username, password string) error {
	user, err := uc.repo.QueryUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if user == nil {
		ctx.Set("USECASE_ERROR", ErrUserNotFound)
		return ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.Set("USECASE_ERROR", ErrPasswordNotMatch)
		return ErrPasswordNotMatch
	}
	if !user.EmailVerified {
		ctx.Set("USECASE_ERROR", ErrEmailNotVerified)
		return ErrEmailNotVerified
	}
	return nil
}

func (uc *SystemUseCase) SendOTP(ctx context.Context, t *entity.User) error {
	if uc.smtpCfg.Host == "" || uc.smtpCfg.Port == 0 || uc.smtpCfg.Username == "" || uc.smtpCfg.Password == "" {
		return errors.New("smtp config not set")
	}

	activationCode := uuid.NewString()
	uc.addActivationCode(activationCode)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("TMT <%s>", uc.smtpCfg.Username))
	m.SetHeader("To", t.Email)
	m.SetHeader("Subject", "Please verify your email address")
	m.SetBody(
		"text/html",
		fmt.Sprintf("Please click the following link in 30 minutes to verify your email address: <a href='https://tocraw.com/tmt/v1/user/verify/%s/%s'>Verify</a>", t.Username, activationCode),
	)

	d := gomail.NewDialer(uc.smtpCfg.Host, 587, uc.smtpCfg.Username, uc.smtpCfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
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

func (uc *SystemUseCase) EncryptPassword(ctx context.Context, password string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func (uc *SystemUseCase) IsPushTokenEnabled(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	dbToken, err := uc.repo.GetPushToken(ctx, token)
	if err != nil {
		return false, err
	}
	if dbToken == nil {
		return false, nil
	}
	return dbToken.Enabled, nil
}

func (uc *SystemUseCase) InsertPushToken(ctx context.Context, token, username string, enabled bool) error {
	if err := uc.repo.InsertOrUpdatePushToken(ctx, token, username, enabled); err != nil {
		return err
	}
	uc.bus.PublishTopicEvent(topicUpdatePushUser)
	return nil
}

func (uc *SystemUseCase) DeleteAllPushTokens(ctx context.Context) error {
	if err := uc.repo.DeleteAllPushTokens(ctx); err != nil {
		return err
	}
	uc.bus.PublishTopicEvent(topicUpdatePushUser)
	return nil
}
