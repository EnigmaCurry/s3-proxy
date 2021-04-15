package bucket

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/config"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/log"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/metrics"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/s3client"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/tracing"
)

// ErrRemovalFolder will be raised when end user is trying to delete a folder and not a file.
var ErrRemovalFolder = errors.New("can't remove folder")

// Client represents a client in order to GET, PUT or DELETE file on a bucket with a html output.
//go:generate mockgen -destination=./mocks/mock_Client.go -package=mocks github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/bucket Client
type Client interface {
	// Get allow to GET what's inside a request path
	Get(ctx context.Context, input *GetInput)
	// Put will put a file following input
	Put(ctx context.Context, inp *PutInput)
	// Delete will delete file on request path
	Delete(ctx context.Context, requestPath string)
	// Load file content. (Should be used internally only).
	LoadFileContent(ctx context.Context, path string) (string, error)
}

// GetInput represents Get input.
type GetInput struct {
	RequestPath       string
	IfModifiedSince   *time.Time
	IfMatch           string
	IfNoneMatch       string
	IfUnmodifiedSince *time.Time
	Range             string
}

// PutInput represents Put input.
type PutInput struct {
	RequestPath string
	Filename    string
	Body        io.ReadSeeker
	ContentType string
	ContentSize int64
}

// NewClient will generate a new client to do GET,PUT or DELETE actions.
func NewClient(
	tgt *config.TargetConfig, tplConfig *config.TemplateConfig, logger log.Logger,
	mountPath string,
	metricsCtx metrics.Client,
	parentTrace tracing.Trace,
	s3clientManager s3client.Manager,
) Client {
	return &requestContext{
		s3ClientManager: s3clientManager,
		targetCfg:       tgt,
		mountPath:       mountPath,
		tplConfig:       tplConfig,
	}
}
