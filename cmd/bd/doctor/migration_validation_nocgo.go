package doctor

// MigrationValidationResult provides machine-parseable migration validation output.
// This stub exists because the Dolt backend has been removed.
type MigrationValidationResult struct {
	Phase           string   `json:"phase"`
	Ready           bool     `json:"ready"`
	Backend         string   `json:"backend"`
	JSONLCount      int      `json:"jsonl_count"`
	SQLiteCount     int      `json:"sqlite_count"`
	DoltCount       int      `json:"dolt_count"`
	MissingInDB     []string `json:"missing_in_db"`
	MissingInJSONL  []string `json:"missing_in_jsonl"`
	Errors          []string `json:"errors"`
	Warnings        []string `json:"warnings"`
	JSONLValid      bool     `json:"jsonl_valid"`
	JSONLMalformed  int      `json:"jsonl_malformed"`
	DoltHealthy     bool     `json:"dolt_healthy"`
	DoltLocked      bool     `json:"dolt_locked"`
	SchemaValid     bool     `json:"schema_valid"`
	RecommendedFix  string   `json:"recommended_fix"`
}

// CheckMigrationReadiness is a stub; Dolt backend has been removed.
func CheckMigrationReadiness(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
		Name:     "Migration Readiness",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryMaintenance,
	}, MigrationValidationResult{
		Phase:   "pre-migration",
		Ready:   false,
		Backend: "unknown",
		Errors:  []string{"Dolt backend has been removed; see github.com/BumpyClock/beads-dolt"},
	}
}

// CheckMigrationCompletion is a stub; Dolt backend has been removed.
func CheckMigrationCompletion(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
		Name:     "Migration Completion",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryMaintenance,
	}, MigrationValidationResult{
		Phase:   "post-migration",
		Ready:   false,
		Backend: "unknown",
		Errors:  []string{"Dolt backend has been removed; see github.com/BumpyClock/beads-dolt"},
	}
}

// CheckDoltLocks is a stub; Dolt backend has been removed.
func CheckDoltLocks(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Locks",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryMaintenance,
	}
}
