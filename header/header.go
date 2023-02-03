package header

type Classification struct {
	// description: |
	//   CVE ID for the template
	// examples:
	//   - value: "\"CVE-2020-14420\""
	CVEID *Slice `yaml:"cve-id,omitempty"`
	// description: |
	//   CWE ID for the template.
	// examples:
	//   - value: "\"CWE-22\""
	CWEID *Slice `yaml:"cwe-id,omitempty"`
	// description: |
	//   CVSS Metrics for the template.
	// examples:
	//   - value: "\"3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H\""
	CVSSMetrics string `yaml:"cvss-metrics,omitempty"`
	// description: |
	//   CVSS Score for the template.
	// examples:
	//   - value: "\"9.8\""
	CVSSScore float64 `yaml:"cvss-score,omitempty"`
}

type Info struct {
	Name           string                 `yaml:"name" json:"name,omitempty"`
	Authors        *Slice                 `yaml:"authors" json:"authors,omitempty"`
	Tags           *Slice                 `yaml:"tags" json:"tags,omitempty"`
	Description    string                 `yaml:"description" json:"description,omitempty"`
	Classification *Classification        `yaml:"classification" json:"classification,omitempty"`
	Reference      []string               `yaml:"reference" json:"reference,omitempty"`
	Remediation    string                 `yaml:"remediation" json:"remediation,omitempty"`
	Metadata       map[string]interface{} `yaml:"metadata" json:"metadata,omitempty"`
}
