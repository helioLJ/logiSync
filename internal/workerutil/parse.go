package workerutil

import "github.com/google/uuid"

func ParseMessage(values map[string]any) (uuid.UUID, string, string, bool) {
	jobIDRaw, ok := values["jobId"]
	if !ok {
		return uuid.UUID{}, "", "", false
	}
	providerRaw, ok := values["provider"]
	if !ok {
		return uuid.UUID{}, "", "", false
	}
	trackingRaw, ok := values["trackingCode"]
	if !ok {
		return uuid.UUID{}, "", "", false
	}

	jobIDStr, ok := jobIDRaw.(string)
	if !ok {
		return uuid.UUID{}, "", "", false
	}
	provider, ok := providerRaw.(string)
	if !ok {
		return uuid.UUID{}, "", "", false
	}
	trackingCode, ok := trackingRaw.(string)
	if !ok {
		return uuid.UUID{}, "", "", false
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return uuid.UUID{}, "", "", false
	}
	return jobID, provider, trackingCode, true
}
