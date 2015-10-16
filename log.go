package cbes

import (
    "fmt"
    "log"
    "runtime"
    "strings"
    "time"
)

const (
    Gray = uint8(iota + 90)
    Red
    Green
    Yellow
    Blue
    Magenta
//NRed      = uint8(31) // Normal
    EndColor = "\033[0m"

    INFO = "INFO"
    TRAC = "TRAC"
    ERRO = "ERRO"
    WARN = "WARN"
    SUCC = "SUCC"
)

// ColorLog colors log and print to stdout.
// See color rules in function 'ColorLogS'.
func ColorLog(format string, a ...interface{}) {
    fmt.Print(ColorLogS(format, a...))
}

// ColorLogS colors log and return colored content.
// Log format: <level> <content [highlight][path]> [ error ].
// Level: TRAC -> blue; ERRO -> red; WARN -> Magenta; SUCC -> green; others -> default.
// Content: default; path: yellow; error -> red.
// Level has to be surrounded by "[" and "]".
// Highlights have to be surrounded by "# " and " #"(space), "#" will be deleted.
// Paths have to be surrounded by "( " and " )"(space).
// Errors have to be surrounded by "[ " and " ]"(space).
// Note: it hasn't support windows yet, contribute is welcome.
func ColorLogS(format string, a ...interface{}) string {
    log := fmt.Sprintf(format, a...)

    var clog string

    if runtime.GOOS != "windows" {
        // Level.
        i := strings.Index(log, "]")
        if log[0] == '[' && i > -1 {
            clog += "[" + getColorLevel(log[1:i]) + "]"
        }

        log = log[i+1:]

        // Error.
        log = strings.Replace(log, "[ ", fmt.Sprintf("[\033[%dm", Red), -1)
        log = strings.Replace(log, " ]", EndColor+"]", -1)

        // Path.
        log = strings.Replace(log, "( ", fmt.Sprintf("(\033[%dm", Yellow), -1)
        log = strings.Replace(log, " )", EndColor+")", -1)

        // Highlights.
        log = strings.Replace(log, "# ", fmt.Sprintf("\033[%dm", Gray), -1)
        log = strings.Replace(log, " #", EndColor, -1)

        log = clog + log

    } else {
        // Level.
        i := strings.Index(log, "]")
        if log[0] == '[' && i > -1 {
            clog += "[" + log[1:i] + "]"
        }

        log = log[i+1:]

        // Error.
        log = strings.Replace(log, "[ ", "[", -1)
        log = strings.Replace(log, " ]", "]", -1)

        // Path.
        log = strings.Replace(log, "( ", "(", -1)
        log = strings.Replace(log, " )", ")", -1)

        // Highlights.
        log = strings.Replace(log, "# ", "", -1)
        log = strings.Replace(log, " #", "", -1)

        log = clog + log
    }

    return time.Now().Format("2006/01/02 15:04:05 ") + log
}

// getColorLevel returns colored level string by given level.
func getColorLevel(level string) string {
    level = strings.ToUpper(level)
    switch level {
    case INFO:
        return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
    case TRAC:
        return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
    case ERRO:
        return fmt.Sprintf("\033[%dm%s\033[0m", Red, level)
    case WARN:
        return fmt.Sprintf("\033[%dm%s\033[0m", Magenta, level)
    case SUCC:
        return fmt.Sprintf("\033[%dm%s\033[0m", Green, level)
    default:
        return level
    }
    return level
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
    var response string
    _, err := fmt.Scanln(&response)
    if err != nil {
        log.Fatal(err)
    }
    okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
    nokayResponses := []string{"n", "N", "no", "No", "NO"}
    if containsString(okayResponses, response) {
        return true
    } else if containsString(nokayResponses, response) {
        return false
    } else {
        fmt.Println("Please type yes or no and then press enter:")
        return askForConfirmation()
    }
}

func containsString(slice []string, element string) bool {
    for _, elem := range slice {
        if elem == element {
            return true
        }
    }
    return false
}

// snake string, XxYy to xx_yy
func snakeString(s string) string {
    data := make([]byte, 0, len(s)*2)
    j := false
    num := len(s)
    for i := 0; i < num; i++ {
        d := s[i]
        if i > 0 && d >= 'A' && d <= 'Z' && j {
            data = append(data, '_')
        }
        if d != '_' {
            j = true
        }
        data = append(data, d)
    }
    return strings.ToLower(string(data[:len(data)]))
}

func camelString(s string) string {
    data := make([]byte, 0, len(s))
    j := false
    k := false
    num := len(s) - 1
    for i := 0; i <= num; i++ {
        d := s[i]
        if k == false && d >= 'A' && d <= 'Z' {
            k = true
        }
        if d >= 'a' && d <= 'z' && (j || k == false) {
            d = d - 32
            j = false
            k = true
        }
        if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
            j = true
            continue
        }
        data = append(data, d)
    }
    return string(data[:len(data)])
}

// The string flag list, implemented flag.Value interface
type strFlags []string

func (s *strFlags) String() string {
    return fmt.Sprintf("%d", *s)
}

func (s *strFlags) Set(value string) error {
    *s = append(*s, value)
    return nil
}