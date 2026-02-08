package doctor

// CheckFederationRemotesAPI is a stub; Dolt backend has been removed.
func CheckFederationRemotesAPI(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Federation remotesapi",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryFederation,
	}
}

// CheckFederationPeerConnectivity is a stub; Dolt backend has been removed.
func CheckFederationPeerConnectivity(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Peer Connectivity",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryFederation,
	}
}

// CheckFederationSyncStaleness is a stub; Dolt backend has been removed.
func CheckFederationSyncStaleness(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Sync Staleness",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryFederation,
	}
}

// CheckFederationConflicts is a stub; Dolt backend has been removed.
func CheckFederationConflicts(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Federation Conflicts",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryFederation,
	}
}

// CheckDoltServerModeMismatch is a stub; Dolt backend has been removed.
func CheckDoltServerModeMismatch(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Mode",
		Status:   StatusOK,
		Message:  "N/A (Dolt backend removed)",
		Category: CategoryFederation,
	}
}


