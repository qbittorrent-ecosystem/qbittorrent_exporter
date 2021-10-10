package collector

import "github.com/prometheus/client_golang/prometheus"

func newMetric(namespace string, name string, helpString string, constLabels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+name, helpString, nil, constLabels)
}

var TORRENT_STATUSES = [6]string{
	"downloading",
	"uploading",
	"complete",
	"checking",
	"errored",
	"paused",
}

const (
	ERROR                = "error"
	MISSING_FILES        = "missingFiles"
	UPLOADING            = "uploading"
	PAUSED_UPLOAD        = "pausedUP"
	QUEUED_UPLOAD        = "queuedUP"
	STALLED_UPLOAD       = "stalledUP"
	CHECKING_UPLOAD      = "checkingUP"
	FORCED_UPLOAD        = "forcedUP"
	ALLOCATING           = "allocating"
	DOWNLOADING          = "downloading"
	METADATA_DOWNLOAD    = "metaDL"
	PAUSED_DOWNLOAD      = "pausedDL"
	QUEUED_DOWNLOAD      = "queuedDL"
	FORCED_DOWNLOAD      = "forcedDL"
	STALLED_DOWNLOAD     = "stalledDL"
	CHECKING_DOWNLOAD    = "checkingDL"
	CHECKING_RESUME_DATA = "checkingResumeData"
	MOVING               = "moving"
	UNKNOWN              = "unknown"
)

type stateFunc func(string) bool

var stateFuncs = map[string]stateFunc{
	"is_downloading": isDownloading,
	"is_uploading":   isUploading,
	"is_complete":    isComplete,
	"is_checking":    isChecking,
	"is_errored":     isErrored,
	"is_paused":      isPaused,
}

func isDownloading(state string) bool {
	return state == DOWNLOADING ||
		state == METADATA_DOWNLOAD ||
		state == PAUSED_DOWNLOAD ||
		state == QUEUED_DOWNLOAD ||
		state == FORCED_DOWNLOAD ||
		state == STALLED_DOWNLOAD ||
		state == CHECKING_DOWNLOAD
}

func isUploading(state string) bool {
	return state == UPLOADING ||
		state == STALLED_UPLOAD ||
		state == CHECKING_UPLOAD ||
		state == QUEUED_UPLOAD ||
		state == FORCED_UPLOAD
}

func isComplete(state string) bool {
	return state == UPLOADING ||
		state == STALLED_UPLOAD ||
		state == CHECKING_UPLOAD ||
		state == PAUSED_UPLOAD ||
		state == QUEUED_UPLOAD ||
		state == FORCED_UPLOAD
}

func isChecking(state string) bool {
	return state == CHECKING_UPLOAD ||
		state == CHECKING_DOWNLOAD ||
		state == CHECKING_RESUME_DATA
}

func isErrored(state string) bool {
	return state == MISSING_FILES || state == ERROR
}

func isPaused(state string) bool {
	return state == PAUSED_UPLOAD || state == PAUSED_DOWNLOAD
}
