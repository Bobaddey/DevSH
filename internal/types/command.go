package types

// Command represents a shell command resolved from natural language.
type Command struct {
	Tool        string  `json:"tool"`         // e.g., kubectl, aws, terraform, bash, docker, git
	Command     string  `json:"command"`      // The actual shell command to execute
	Confidence      float64  `json:"confidence"`      // 0.0 to 1.0 indicating model's confidence
	Explanation     string   `json:"explanation"`     // Explanation of what the command does
	RiskLevel       string   `json:"risk_level"`      // "low", "medium", "high"
	Recommendations []string `json:"recommendations"` // Optional troubleshooting steps or advice
	Insights        string   `json:"insights"`        // Detailed technical explanation of the issue
}
