package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(1).
			Width(0)

	quoteStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#9F87FF")).
			MarginTop(1).
			MarginBottom(1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginTop(1).
			MarginBottom(0)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			MarginLeft(1)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Align(lipgloss.Center).
			MarginTop(1)

	containerStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)
)

type ResumeData struct {
	Personal struct {
		Name     string `yaml:"name"`
		Title    string `yaml:"title"`
		Quote    string `yaml:"quote"`
		Summary  string `yaml:"summary"`
		Location string `yaml:"location"`
	} `yaml:"personal"`
	Philosophy    []string `yaml:"philosophy"`
	CurrentFocus  []string `yaml:"current_focus"`
	Education     struct {
		Degree   string `yaml:"degree"`
		School   string `yaml:"school"`
		Period   string `yaml:"period"`
		Location string `yaml:"location"`
	} `yaml:"education"`
	Experience []struct {
		Title        string   `yaml:"title"`
		Organization string   `yaml:"organization"`
		Period       string   `yaml:"period"`
		Details      []string `yaml:"details"`
	} `yaml:"experience"`
	Github struct {
		Username     string `yaml:"username"`
		Repositories string `yaml:"repositories"`
		Status       string `yaml:"status"`
	} `yaml:"github"`
	Projects []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Tech        string `yaml:"tech"`
		Link        string `yaml:"link"`
	} `yaml:"projects"`
	Blog struct {
		Name         string `yaml:"name"`
		Description  string `yaml:"description"`
		Tagline      string `yaml:"tagline"`
		ContentFocus string `yaml:"content_focus"`
	} `yaml:"blog"`
	Skills struct {
		Categories []struct {
			Category    string   `yaml:"category"`
			Description string   `yaml:"description"`
			Skills      []string `yaml:"skills"`
		} `yaml:"categories"`
		Languages []string `yaml:"languages"`
	} `yaml:"skills"`
	Achievements struct {
		Hackathons    []string `yaml:"hackathons"`
		Leadership    []string `yaml:"leadership"`
		Presentations []string `yaml:"presentations"`
		Community     []string `yaml:"community"`
	} `yaml:"achievements"`
	Contact struct {
		Websites     []string `yaml:"websites"`
		Professional []string `yaml:"professional"`
		Info         []string `yaml:"info"`
		Note         string   `yaml:"note"`
		AboutResume  struct {
			Description string `yaml:"description"`
			Source      string `yaml:"source"`
		} `yaml:"about_resume"`
	} `yaml:"contact"`
}

type model struct {
	viewport    viewport.Model
	ready       bool
	currentPage int
	pages       []string
	resumeData  ResumeData
}

func loadResumeData() (ResumeData, error) {
	var data ResumeData
	
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		return data, fmt.Errorf("failed to get executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)
	
	// Construct the path to resume.yaml relative to the executable
	yamlPath := filepath.Join(execDir, "resume.yaml")
	
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return data, fmt.Errorf("failed to read resume.yaml from %s: %v", yamlPath, err)
	}
	
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return data, fmt.Errorf("failed to parse YAML: %v", err)
	}
	
	return data, nil
}

func initialModel() model {
	data, err := loadResumeData()
	if err != nil {
		fmt.Printf("Error loading resume data: %v\n", err)
		os.Exit(1)
	}

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

func wrapText(text string, width int) string {
	// Ensure we have a reasonable minimum width
	if width <= 30 {
		width = 80 // Default to reasonable width if viewport is too small
	}

	// Account for minimal padding only
	effectiveWidth := width - 4 // Leave small room for padding
	if effectiveWidth <= 30 {
		effectiveWidth = width // Use full width if we're already tight on space
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= effectiveWidth {
			currentLine.WriteString(" " + word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

func (m model) buildProfilePage() string {
	var b strings.Builder

	// Build the header
	header := titleStyle.Render(m.resumeData.Personal.Name)
	b.WriteString(header + "\n")

	subtitle := m.resumeData.Personal.Title
	b.WriteString(itemStyle.Render(subtitle) + "\n")

	// Add quote
	b.WriteString("\n" + quoteStyle.Render(m.resumeData.Personal.Quote) + "\n")

	// Summary - wrap text to prevent overflow
	wrappedSummary := wrapText(m.resumeData.Personal.Summary, m.viewport.Width)
	b.WriteString("\n" + wrappedSummary + "\n")

	// Personal Philosophy
	b.WriteString("\n" + sectionStyle.Render("PERSONAL PHILOSOPHY") + "\n")
	for _, philosophy := range m.resumeData.Philosophy {
		wrappedPhilosophy := wrapText("• \""+philosophy+"\"", m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedPhilosophy) + "\n")
	}

	// Current Focus
	b.WriteString("\n" + sectionStyle.Render("CURRENT FOCUS") + "\n")
	for _, focus := range m.resumeData.CurrentFocus {
		b.WriteString(itemStyle.Render("• " + focus) + "\n")
	}

	// Location
	b.WriteString("\n" + sectionStyle.Render("LOCATION") + "\n")
	b.WriteString(itemStyle.Render(m.resumeData.Personal.Location) + "\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) buildExperiencePage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("EDUCATION & EXPERIENCE")
	b.WriteString(header + "\n\n")

	// Education
	b.WriteString(sectionStyle.Render("EDUCATION") + "\n")

	edu := fmt.Sprintf("%s, %s", highlightStyle.Render(m.resumeData.Education.Degree), m.resumeData.Education.School)
	eduPeriod := fmt.Sprintf("%s | %s", m.resumeData.Education.Period, m.resumeData.Education.Location)
	b.WriteString(itemStyle.Render(edu) + "\n")
	b.WriteString(itemStyle.Render(eduPeriod) + "\n\n")

	// Leadership Experience
	b.WriteString(sectionStyle.Render("LEADERSHIP") + "\n")

	for _, exp := range m.resumeData.Experience {
		job := fmt.Sprintf("%s, %s", highlightStyle.Render(exp.Title), exp.Organization)
		period := fmt.Sprintf("%s", exp.Period)
		b.WriteString(itemStyle.Render(job) + "\n")
		b.WriteString(itemStyle.Render(period) + "\n")

		for _, detail := range exp.Details {
			wrappedDetail := wrapText("• "+detail, m.viewport.Width)
			b.WriteString(itemStyle.Render(wrappedDetail) + "\n")
		}
		b.WriteString("\n")
	}

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) buildProjectsPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("PROJECTS & CONTRIBUTIONS")
	b.WriteString(header + "\n\n")

	// GitHub Overview
	b.WriteString(sectionStyle.Render("GITHUB PROFILE") + "\n")
	b.WriteString(itemStyle.Render("Username: " + m.resumeData.Github.Username) + "\n")
	b.WriteString(itemStyle.Render(m.resumeData.Github.Repositories) + "\n")
	b.WriteString(itemStyle.Render(m.resumeData.Github.Status) + "\n\n")

	// Significant Projects
	b.WriteString(sectionStyle.Render("SIGNIFICANT PROJECTS") + "\n")

	for _, project := range m.resumeData.Projects {
		proj := highlightStyle.Render(project.Name)
		b.WriteString(itemStyle.Render(proj) + "\n")
		wrappedDesc := wrapText("Description: "+project.Description, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedDesc) + "\n")
		wrappedTech := wrapText("Technologies: "+project.Tech, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedTech) + "\n")
		b.WriteString(itemStyle.Render("Link: " + project.Link) + "\n\n")
	}

	// Blog and Research
	b.WriteString(sectionStyle.Render("BLOG & RESEARCH") + "\n")
	wrappedBlogName := wrapText(m.resumeData.Blog.Name+" - "+m.resumeData.Blog.Description, m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedBlogName) + "\n")
	wrappedTagline := wrapText(m.resumeData.Blog.Tagline, m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedTagline) + "\n")
	wrappedContentFocus := wrapText(m.resumeData.Blog.ContentFocus, m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedContentFocus) + "\n\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) buildSkillsPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("TECHNICAL SPECIALIZATIONS")
	b.WriteString(header + "\n\n")

	// Core technical specializations
	for _, category := range m.resumeData.Skills.Categories {
		cat := sectionStyle.Render(category.Category)
		b.WriteString(cat + "\n")
		wrappedDesc := wrapText(category.Description, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedDesc) + "\n")
		wrappedSkills := wrapText("Key skills: "+strings.Join(category.Skills, ", "), m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedSkills) + "\n\n")
	}

	// Programming Languages
	b.WriteString(sectionStyle.Render("PROGRAMMING LANGUAGES") + "\n")
	wrappedLanguages := wrapText(strings.Join(m.resumeData.Skills.Languages, ", "), m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedLanguages) + "\n\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) buildAchievementsPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("ACHIEVEMENTS & RECOGNITION")
	b.WriteString(header + "\n\n")

	// Hackathon Victories
	b.WriteString(sectionStyle.Render("HACKATHON VICTORIES") + "\n")
	for _, hackathon := range m.resumeData.Achievements.Hackathons {
		wrappedHackathon := wrapText("• "+hackathon, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedHackathon) + "\n")
	}
	b.WriteString("\n")

	// Leadership Recognition
	b.WriteString(sectionStyle.Render("LEADERSHIP RECOGNITION") + "\n")
	for _, leadership := range m.resumeData.Achievements.Leadership {
		wrappedLeadership := wrapText("• "+leadership, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedLeadership) + "\n")
	}
	b.WriteString("\n")

	// Technical Presentations
	b.WriteString(sectionStyle.Render("TECHNICAL PRESENTATIONS") + "\n")
	for _, presentation := range m.resumeData.Achievements.Presentations {
		wrappedPresentation := wrapText("• "+presentation, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedPresentation) + "\n")
	}
	b.WriteString("\n")

	// Community Building
	b.WriteString(sectionStyle.Render("COMMUNITY IMPACT") + "\n")
	for _, community := range m.resumeData.Achievements.Community {
		wrappedCommunity := wrapText("• "+community, m.viewport.Width)
		b.WriteString(itemStyle.Render(wrappedCommunity) + "\n")
	}
	b.WriteString("\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) buildContactPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("CONNECT WITH DANIEL")
	b.WriteString(header + "\n\n")

	// Websites
	b.WriteString(sectionStyle.Render("PERSONAL WEBSITES") + "\n")
	for _, website := range m.resumeData.Contact.Websites {
		b.WriteString(itemStyle.Render("• " + website) + "\n")
	}
	b.WriteString("\n")

	// Social & Professional
	b.WriteString(sectionStyle.Render("PROFESSIONAL PROFILES") + "\n")
	for _, professional := range m.resumeData.Contact.Professional {
		b.WriteString(itemStyle.Render("• " + professional) + "\n")
	}
	b.WriteString("\n")

	// Contact
	b.WriteString(sectionStyle.Render("GET IN TOUCH") + "\n")
	for _, info := range m.resumeData.Contact.Info {
		b.WriteString(itemStyle.Render("• " + info) + "\n")
	}
	b.WriteString("\n")

	// Final Note
	b.WriteString(quoteStyle.Render(m.resumeData.Contact.Note) + "\n\n")

	// This Resume
	b.WriteString(sectionStyle.Render("ABOUT THIS RESUME") + "\n")
	wrappedDescription := wrapText(m.resumeData.Contact.AboutResume.Description, m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedDescription) + "\n")
	wrappedSource := wrapText(m.resumeData.Contact.AboutResume.Source, m.viewport.Width)
	b.WriteString(itemStyle.Render(wrappedSource) + "\n\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n\n" + footer)

	return b.String()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "right", "space":
			if m.currentPage < len(m.pages)-1 {
				m.currentPage++
				m.viewport.SetContent(m.pages[m.currentPage])
				m.viewport.GotoTop()
			}
		case "p", "left", "b":
			if m.currentPage > 0 {
				m.currentPage--
				m.viewport.SetContent(m.pages[m.currentPage])
				m.viewport.GotoTop()
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetContent(m.pages[m.currentPage])
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	pageIndicator := fmt.Sprintf("[Page %d/%d]", m.currentPage+1, len(m.pages))
	content := containerStyle.Render(m.viewport.View())
	return fmt.Sprintf("%s\n%s", content, pageIndicator)
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
