# TUICV - Interactive Terminal Resume

A beautiful, interactive terminal-based resume application built with Go and Bubble Tea. Features a clean TUI interface with proper text wrapping, configurable content via YAML, and smooth navigation between sections.

## Features

✨ **Modern TUI Interface** - Clean, responsive terminal interface with Bubble Tea  
📱 **Responsive Design** - Automatically adapts to different terminal sizes  
⚙️ **YAML Configuration** - Easy content management without touching code  
🎨 **Beautiful Styling** - Professional color scheme and typography  
🔄 **Smart Text Wrapping** - Prevents text overflow on any screen size  
📄 **Multi-page Layout** - Organized sections for easy navigation  

## Quick Start

### Prerequisites

- Go 1.25.0 or later
- Terminal with color support

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd TUICV
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Customize your resume**
   Edit `resume.yaml` with your personal information:
   ```bash
   nano resume.yaml
   ```

4. **Build and run**
   ```bash
   go build -o resume-tui .
   ./resume-tui
   ```

## Configuration

The resume content is entirely configured through the `resume.yaml` file. This separation allows you to update your information without modifying any code.

### YAML Structure

```yaml
personal:
  name: "Your Name"
  title: "Your Title"
  quote: "Your favorite quote"
  summary: "Brief professional summary"
  location: "Your location"

philosophy:
  - "Your professional philosophy point 1"
  - "Your professional philosophy point 2"

current_focus:
  - "Current area of focus 1"
  - "Current area of focus 2"

education:
  degree: "Your degree"
  school: "Your school"
  period: "Time period"
  location: "School location"

experience:
  - title: "Job Title"
    organization: "Company Name"
    period: "Time period"
    details:
      - "Achievement or responsibility 1"
      - "Achievement or responsibility 2"

# ... and more sections
```

### Key Sections

- **Personal Info** - Name, title, summary, location
- **Philosophy** - Your professional beliefs and approach
- **Current Focus** - What you're working on now
- **Education** - Academic background
- **Experience** - Work history with details
- **Projects** - Significant projects with descriptions
- **Skills** - Technical skills organized by category
- **Achievements** - Awards, recognition, community impact
- **Contact** - How to reach you

## Usage

### Navigation

- **n** / **→** / **Space** - Next page
- **p** / **←** / **b** - Previous page  
- **q** / **Ctrl+C** - Quit application
- **Mouse scroll** - Scroll through content (if supported)

### Pages

1. **Profile** - Personal info, philosophy, current focus
2. **Experience** - Education and work history
3. **Projects** - Significant projects and contributions
4. **Skills** - Technical specializations and languages
5. **Achievements** - Awards, recognition, presentations
6. **Contact** - How to connect with you

## SSH Integration

The app now includes its own SSH server, so visitors can connect directly without a special username or SSH client config.

### Run as SSH Server

1. **Build the binary**
   ```bash
   go build -o resume-tui .
   ```

2. **Start SSH mode**
   ```bash
   ./resume-tui serve
   ```

   By default it listens on `:22`. You can override it with:
   ```bash
   TUI_CV_LISTEN_ADDR=:2322 ./resume-tui serve
   ```

3. **Connect from any machine**
   ```bash
   ssh your-domain.com
   ```

No `cv@` prefix is required. Any SSH username is accepted and routed to the resume TUI.

You can also force a stable host key path:
```bash
TUI_CV_HOST_KEY=/path/to/ssh_host_ed25519_key ./resume-tui serve
```

## Development

### Project Structure

```
TUICV/
├── main.go           # Bubble Tea TUI application
├── server.go         # Embedded SSH server entrypoint
├── resume.yaml       # Resume content configuration
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # This file
```

### Key Dependencies

- **Bubble Tea** - TUI framework for terminal applications
- **Lipgloss** - Styling and layout for terminal interfaces
- **Bubbles** - Common UI components (viewport)
- **Gliderlabs SSH** - Embedded SSH server
- **YAML v2** - YAML parsing for configuration

### Building

```bash
# Development build
go build -o resume-tui .

# Production build with optimizations
go build -ldflags="-s -w" -o resume-tui .
```

### Customization

The application uses a modular design where:

- **Styling** is defined in the style variables at the top of `main.go`
- **Content** is loaded from `resume.yaml`
- **Layout logic** handles responsive text wrapping and spacing
- **Navigation** provides smooth page transitions

## Technical Details

- **Text Wrapping** - Smart algorithm prevents overflow on any terminal size
- **Responsive Layout** - Adapts to terminal width changes in real-time
- **Memory Efficient** - Lazy loading of content with viewport optimization
- **Cross-platform** - Works on Linux, macOS, and Windows terminals

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Inspiration

This project demonstrates how traditional resume formats can be reimagined for the terminal, offering an interactive and memorable way to present professional information to technical audiences.

---

**Built with ❤️ using Go and Bubble Tea**
