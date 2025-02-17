package server

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/common/debug"
)

//go:embed templates/error.html
var errorTemplates embed.FS

type StackFrameContext struct {
	Frame      debug.StackFrame
	CodeLines  []CodeLine
	IsInternal bool
}

type CodeLine struct {
	Number    int
	Content   string
	IsCurrent bool
}

func ErrorHandler(isDev bool, c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if c.Get("X-Requested-With") != "XMLHttpRequest" {
		if isDev {
			if stackTrace, ok := c.Locals("stack_trace").([]debug.StackFrame); ok {
				html := errorPageWithTrace(stackTrace, err, code)
				return c.Type("html").Status(code).Send([]byte(html))
			}
		}

		html := errorPage(err, code)
		return c.Type("html").Status(code).Send([]byte(html))
	}

	detail := ""
	if isDev {
		var stackFrames []debug.StackFrame
		if stack, ok := c.Locals("stack_trace").([]debug.StackFrame); ok {
			stackFrames = stack
		} else if stack, err := debug.StackTrace(); err == nil {
			stackFrames = stack
		}

		for _, frame := range stackFrames {
			detail += fmt.Sprintf("%s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		}
	}

	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}

	response.Status = "error"
	response.Message = err.Error()
	response.Detail = detail

	return c.Status(code).JSON(response)
}

func getModulePath() string {
	file, err := os.Open("go.mod")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func isInternalFrame(frame debug.StackFrame, modulePath string) bool {
	if modulePath == "" {
		return false
	}
	return strings.Contains(frame.Function, modulePath)
}

func getFileContext(filename string, targetLine int, contextLines int) []CodeLine {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []CodeLine{}
	lineNum := 1
	start := targetLine - contextLines
	end := targetLine + contextLines

	for scanner.Scan() {
		if lineNum >= start && lineNum <= end {
			lines = append(lines, CodeLine{
				Number:    lineNum,
				Content:   scanner.Text(),
				IsCurrent: lineNum == targetLine,
			})
		}
		lineNum++
	}

	return lines
}

func errorPageWithTrace(stackTrace []debug.StackFrame, err error, code int) string {
	modulePath := getModulePath()
	frames := []StackFrameContext{}

	for _, frame := range stackTrace {
		frames = append(frames, StackFrameContext{
			Frame:      frame,
			CodeLines:  getFileContext(frame.File, frame.Line, 10),
			IsInternal: isInternalFrame(frame, modulePath),
		})
	}

	tmpl, _ := template.ParseFS(errorTemplates, "templates/error.html")
	data := struct {
		Code    int
		Message string
		Frames  []StackFrameContext
	}{
		Code:    code,
		Message: err.Error(),
		Frames:  frames,
	}

	var output strings.Builder
	tmpl.Execute(&output, data)
	return output.String()
}

func errorPage(err error, code int) string {
	return errorPageWithTrace(nil, err, code)
}
