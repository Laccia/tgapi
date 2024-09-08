package sign

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

type Sign struct {
	PhoneNumber string
	Code2FA     string
}

func (Sign) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("signing up not implemented in Sign")
}

func (Sign) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (a Sign) Phone(_ context.Context) (string, error) {
	if a.PhoneNumber != "" {
		return a.PhoneNumber, nil
	}
	fmt.Print("Enter phone in international format (e.g. +1234567890): ")
	phone, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(phone), nil
}

func (Sign) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (Sign) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {

	code := os.Getenv("CODE")

	for code == "" {
		err := godotenv.Overload(".env")
		if err != nil {
			return "", err
		}
		code = os.Getenv("CODE")
		log.Info().Str("AUTH", "SignIn").Msg("Waiting 2FA code")
		time.Sleep(1 * time.Second)
		if code != "" {
			break
		}
	}

	return strings.TrimSpace(code), nil
}
