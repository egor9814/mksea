package main

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"mksea/common"
	"mksea/input"
	"os"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
)

var errInvPass = errors.New("incorrect password")
var passwordAttempts = 5

func decodeEncoderKey(password []byte) bool {
	expected := []byte(common.PasswordTestTemplate())
	if len(expected) != len(input.Env.PasswordTest) {
		return false
	}
	passwordHash := md5.Sum(password)
	j, k := 0, 0
	for i, it := range input.Env.PasswordTest {
		actual := it ^ passwordHash[j] ^ input.Env.DecodeKey[k]
		if actual != expected[i] {
			return false
		}
		j = (j + 1) % len(passwordHash)
		k = (k + 1) % len(input.Env.DecodeKey)
	}
	j = 0
	for i, it := range input.Env.DecodeKey {
		input.Env.DecodeKey[i] = it ^ passwordHash[j]
		j = (j + 1) % len(passwordHash)
	}
	return true
}

func zenityPassword() error {
	for i := passwordAttempts; i > 0; i-- {
		title := "Type archive password"
		if i != passwordAttempts {
			title += ". Attempts left: " + strconv.Itoa(i)
		}
		_, p, err := zenity.Password(
			zenity.Title(title),
		)
		if err == zenity.ErrCanceled {
			return err
		}
		if err != nil {
			return err
		}
		if decodeEncoderKey([]byte(p)) {
			return nil
		}
	}
	return errInvPass
}

func scanPassword() error {
	reader := bufio.NewReader(os.Stdin)
	for i := passwordAttempts; i > 0; i-- {
		title := "type archive password"
		if i != passwordAttempts {
			title += " (attempts left: " + strconv.Itoa(i) + ")"
		}
		title += "> "
		fmt.Print(title)
		if text, err := reader.ReadString('\n'); err != nil {
			return common.NewContextError("cannot read password from terminal", err)
		} else if decodeEncoderKey([]byte(strings.TrimSpace(text))) {
			return nil
		}
	}
	return errInvPass
}

func testPassword(password []byte) error {
	if len(input.Env.PasswordTest) == 0 {
		return nil
	}
	if len(password) == 0 {
		if zenity.IsAvailable() {
			return zenityPassword()
		}
		return scanPassword()
	}
	if !decodeEncoderKey(password) {
		return errInvPass
	}
	return nil
}
