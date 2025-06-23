package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(4).
			PaddingRight(4).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginTop(1).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Align(lipgloss.Center).
			MarginTop(1)
)

type model struct {
	viewport    viewport.Model
	ready       bool
	currentPage int
	pages       []string
}

func initialModel() model {
	return model{
		viewport:    viewport.Model{},
		currentPage: 0,
		pages: []string{
			buildResumePage(),
			buildProjectsPage(),
			buildSkillsPage(),
			buildContactPage(),
		},
	}
}

func buildResumePage() string {
	var b strings.Builder

	// CUSTOMIZE THIS SECTION WITH YOUR INFO
	name := "Daniel Alves"
	title := "Software Engineer"
	summary := "Passionate software engineer with expertise in systems programming, web development, and DevOps."

	// Build the header
	header := titleStyle.Render(fmt.Sprintf("%s - %s", name, title))
	b.WriteString(header + "\n\n")

	b.WriteString(summary + "\n\n")

	// Experience
	b.WriteString(sectionStyle.Render("EXPERIENCE") + "\n")

	// CUSTOMIZE THESE WITH YOUR ACTUAL EXPERIENCE
	experiences := []struct {
		title    string
		company  string
		period   string
		location string
		details  []string
	}{
		{
			title:    "Senior Software Engineer",
			company:  "Example Corp",
			period:   "2022 - Present",
			location: "Remote",
			details: []string{
				"Led development of microservices architecture",
				"Implemented CI/CD pipelines reducing deployment time by 70%",
				"Mentored junior developers on best practices",
			},
		},
		{
			title:    "Software Engineer",
			company:  "Tech Startup",
			period:   "2019 - 2022",
			location: "New York, NY",
			details: []string{
				"Developed backend services using Go and Rust",
				"Created RESTful APIs for mobile applications",
				"Optimized database queries improving performance by 40%",
			},
		},
	}

	for _, exp := range experiences {
		job := fmt.Sprintf("%s, %s", highlightStyle.Render(exp.title), exp.company)
		period := fmt.Sprintf("%s | %s", exp.period, exp.location)
		b.WriteString(itemStyle.Render(job) + "\n")
		b.WriteString(itemStyle.Render(period) + "\n")

		for _, detail := range exp.details {
			b.WriteString(itemStyle.Render("• " + detail) + "\n")
		}
		b.WriteString("\n")
	}

	// Education
	b.WriteString(sectionStyle.Render("EDUCATION") + "\n")

	// CUSTOMIZE THIS WITH YOUR EDUCATION
	education := struct {
		degree   string
		school   string
		period   string
		location string
	}{
		degree:   "Bachelor of Science in Computer Science",
		school:   "University of Technology",
		period:   "2015 - 2019",
		location: "San Francisco, CA",
	}

	edu := fmt.Sprintf("%s, %s", highlightStyle.Render(education.degree), education.school)
	eduPeriod := fmt.Sprintf("%s | %s", education.period, education.location)
	b.WriteString(itemStyle.Render(edu) + "\n")
	b.WriteString(itemStyle.Render(eduPeriod) + "\n\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n" + footer)

	return b.String()
}

func buildProjectsPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("PROJECTS")
	b.WriteString(header + "\n\n")

	// CUSTOMIZE THESE WITH YOUR ACTUAL PROJECTS
	projects := []struct {
		name        string
		description string
		tech        string
		link        string
	}{
		{
			name:        "SSH Resume TUI",
			description: "Interactive terminal-based resume accessible via SSH",
			tech:        "Go, Bubble Tea, SSH",
			link:        "cv.alves.world",
		},
		{
			name:        "Personal Website",
			description: "Portfolio website showcasing projects and blog",
			tech:        "React, Next.js, Tailwind CSS",
			link:        "github.com/yourusername/website",
		},
		{
			name:        "API Gateway",
			description: "High-performance API gateway with custom authentication",
			tech:        "Rust, Redis, PostgreSQL",
			link:        "github.com/yourusername/api-gateway",
		},
	}

	for _, project := range projects {
		proj := highlightStyle.Render(project.name)
		b.WriteString(itemStyle.Render(proj) + "\n")
		b.WriteString(itemStyle.Render("Description: " + project.description) + "\n")
		b.WriteString(itemStyle.Render("Technologies: " + project.tech) + "\n")
		b.WriteString(itemStyle.Render("Link: " + project.link) + "\n\n")
	}

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n" + footer)

	return b.String()
}

func buildSkillsPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("SKILLS & EXPERTISE")
	b.WriteString(header + "\n\n")

	// CUSTOMIZE THESE WITH YOUR ACTUAL SKILLS
	skillCategories := []struct {
		category string
		skills   []string
	}{
		{
			category: "Programming Languages",
			skills:   []string{"Go", "Rust", "JavaScript/TypeScript", "Python", "C/C++"},
		},
		{
			category: "Frameworks & Libraries",
			skills:   []string{"React", "Node.js", "Express", "Django", "Actix Web"},
		},
		{
			category: "Tools & Platforms",
			skills:   []string{"Docker", "Kubernetes", "AWS", "GitHub Actions", "Terraform"},
		},
		{
			category: "Databases",
			skills:   []string{"PostgreSQL", "MongoDB", "Redis", "Elasticsearch"},
		},
	}

	for _, category := range skillCategories {
		cat := sectionStyle.Render(category.category)
		b.WriteString(cat + "\n")
		b.WriteString(itemStyle.Render(strings.Join(category.skills, ", ")) + "\n\n")
	}

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n" + footer)

	return b.String()
}

func buildContactPage() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("CONTACT & LINKS")
	b.WriteString(header + "\n\n")

	// CUSTOMIZE THESE WITH YOUR ACTUAL CONTACT INFO
	contacts := []struct {
		method string
		value  string
	}{
		{
			method: "Email",
			value:  "you@example.com",
		},
		{
			method: "GitHub",
			value:  "github.com/yourusername",
		},
		{
			method: "LinkedIn",
			value:  "linkedin.com/in/yourusername",
		},
		{
			method: "Website",
			value:  "yourdomain.com",
		},
		{
			method: "Twitter",
			value:  "@yourusername",
		},
	}

	for _, contact := range contacts {
		method := highlightStyle.Render(contact.method)
		b.WriteString(itemStyle.Render(fmt.Sprintf("%s: %s", method, contact.value)) + "\n")
	}

	b.WriteString("\n\n")
	b.WriteString(itemStyle.Render("Thank you for checking out my resume!") + "\n")
	b.WriteString(itemStyle.Render("Feel free to reach out for opportunities or collaboration.") + "\n\n")

	// Footer with navigation help
	footer := footerStyle.Render("Press 'n' for next page, 'p' for previous, 'q' to quit")
	b.WriteString("\n" + footer)

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
	return fmt.Sprintf("%s\n%s", m.viewport.View(), pageIndicator)
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
