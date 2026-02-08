package doctor

// MigrationValidationResult provides machine-parseable migration validation output.
// This stub exists because legacy backend migration support has been removed.
type MigrationValidationResult struct {
	Phase          string   `json:"phase"`
	Ready          bool     `json:"ready"`
	Backend        string   `json:"backend"`
	JSONLCount     int      `json:"jsonl_count"`
	SQLiteCount    int      `json:"sqlite_count"`
	LegacyCount    int      `json:"legacy_count"`
	MissingInDB    []string `json:"missing_in_db"`
	MissingInJSONL []string `json:"missing_in_jsonl"`
	Errors         []string `json:"errors"`
	Warnings       []string `json:"warnings"`
	JSONLValid     bool     `json:"jsonl_valid"`
	JSONLMalformed int      `json:"jsonl_malformed"`
	LegacyHealthy  bool     `json:"legacy_healthy"`
	LegacyLocked   bool     `json:"legacy_locked"`
	SchemaValid    bool     `json:"schema_valid"`
	RecommendedFix string   `json:"recommended_fix"`
}

// CheckMigrationReadiness is a stub; legacy backend migration has been removed.
func CheckMigrationReadiness(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
			Name:     "Migration Readiness",
			Status:   StatusOK,
			Message:  "N/A (legacy backend removed)",
			Category: CategoryMaintenance,
		}, MigrationValidationResult{
			Phase:   "pre-migration",
			Ready:   false,
			Backend: "unknown",
			Errors:  []string{"legacy backend migration has been removed"},
		}
}

// CheckMigrationCompletion is a stub; legacy backend migration has been removed.
func CheckMigrationCompletion(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
			Name:     "Migration Completion",
			Status:   StatusOK,
			Message:  "N/A (legacy backend removed)",
			Category: CategoryMaintenance,
		}, MigrationValidationResult{
			Phase:   "post-migration",
			Ready:   false,
			Backend: "unknown",
			Errors:  []string{"legacy backend migration has been removed"},
		}
}

// CheckLegacyLocks is a stub; legacy backend migration has been removed.
func CheckLegacyLocks(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Legacy Locks",
		Status:   StatusOK,
		Message:  "N/A (legacy backend removed)",
		Category: CategoryMaintenance,
	}
}
