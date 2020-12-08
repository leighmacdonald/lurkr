package sftp

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
)

type sftpTransport struct {
	cfg    config.SFTPConfig
	client *sftp.Client
}

func (t *sftpTransport) Send(reader io.Reader, path string) error {
	dst, err := t.client.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to create dest file")
	}
	if _, err := io.Copy(dst, reader); err != nil {
		return errors.Wrapf(err, "Failed to write file on remote host")
	}
	log.Infof("Send file successfully: %v", path)
	return nil
}

func NewSFTPTransport(cfg config.SFTPConfig) (*sftpTransport, error) {
	var authMethods []ssh.AuthMethod
	if cfg.Key != "" {
		key, err := ioutil.ReadFile(cfg.Key)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read private key")
		}
		var (
			signer ssh.Signer
			errKey error
		)
		if cfg.Password == "" {
			signer, errKey = ssh.ParsePrivateKey(key)
			if errKey != nil {
				return nil, errors.Wrapf(errKey, "Failed to parse private key")
			}
		} else {
			signer, errKey = ssh.ParsePrivateKeyWithPassphrase(key, []byte(cfg.Password))
			if errKey != nil {
				return nil, errors.Wrapf(errKey, "Failed to parse private key (passphrase)")
			}
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}
	conn, err := ssh.Dial("tcp", cfg.Address, &ssh.ClientConfig{
		User: cfg.Username,
		Auth: authMethods,
		// TODO: Check host key if requested
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Ciphers: []string{"3des-cbc", "aes256-cbc", "aes192-cbc", "aes128-cbc"},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to ssh server")
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create SFTP client")
	}
	return &sftpTransport{client: client, cfg: cfg}, nil
}
