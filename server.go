package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gliderssh "github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func buildModel(data ResumeData) model {
	m := model{
		viewport:    viewport.Model{},
		currentPage: 0,
		resumeData:  data,
	}

	m.pages = []string{
		m.buildProfilePage(),
		m.buildExperiencePage(),
		m.buildProjectsPage(),
		m.buildSkillsPage(),
		m.buildAchievementsPage(),
		m.buildContactPage(),
	}

	return m
}

func defaultHostKeyPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	return filepath.Join(filepath.Dir(execPath), "ssh_host_ed25519_key"), nil
}

func loadOrCreateHostSigner(keyPath string) (gossh.Signer, error) {
	privateKeyBytes, err := os.ReadFile(keyPath)
	if err == nil {
		signer, parseErr := gossh.ParsePrivateKey(privateKeyBytes)
		if parseErr == nil {
			return signer, nil
		}
	}

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate host key: %w", err)
	}

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal host key: %w", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})

	if writeErr := os.WriteFile(keyPath, pemBytes, 0o600); writeErr != nil {
		return nil, fmt.Errorf("failed to write host key: %w", writeErr)
	}

	signer, err := gossh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated host key: %w", err)
	}

	return signer, nil
}

func handleSSHSession(s gliderssh.Session, data ResumeData) {
	ptyReq, winCh, hasPty := s.Pty()
	if !hasPty {
		_, _ = s.Write([]byte("A TTY is required. Use: ssh -t <host>\n"))
		_ = s.Exit(1)
		return
	}

	p := tea.NewProgram(
		buildModel(data),
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		p.Send(tea.WindowSizeMsg{Width: ptyReq.Window.Width, Height: ptyReq.Window.Height})
		for win := range winCh {
			p.Send(tea.WindowSizeMsg{Width: win.Width, Height: win.Height})
		}
	}()

	if _, err := p.Run(); err != nil {
		_, _ = s.Write([]byte(fmt.Sprintf("Error running program: %v\n", err)))
		_ = s.Exit(1)
		return
	}

	_ = s.Exit(0)
}

func runSSHServer() error {
	data, err := loadResumeData()
	if err != nil {
		return fmt.Errorf("failed to load resume data: %w", err)
	}

	keyPath, err := defaultHostKeyPath()
	if err != nil {
		return err
	}

	if envPath := os.Getenv("TUI_CV_HOST_KEY"); envPath != "" {
		keyPath = envPath
	}

	signer, err := loadOrCreateHostSigner(keyPath)
	if err != nil {
		return err
	}

	addr := ":22"
	if envAddr := os.Getenv("TUI_CV_LISTEN_ADDR"); envAddr != "" {
		addr = envAddr
	}

	server := &gliderssh.Server{
		Addr:        addr,
		Handler:     func(s gliderssh.Session) { handleSSHSession(s, data) },
		HostSigners: []gliderssh.Signer{signer},
	}

	return server.ListenAndServe()
}
