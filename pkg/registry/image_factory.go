package registry

import (
	"time"

	"github.com/buildpack/imgutil"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/pkg/errors"
)

type ImageFactory struct {
	KeychainFactory KeychainFactory
}

func (f *ImageFactory) NewRemote(imageRef ImageRef) (RemoteImage, error) {
	remote, err := imgutil.NewRemoteImage(imageRef.RepoName(), f.KeychainFactory.KeychainForImageRef(imageRef))
	return remote, errors.Wrapf(err, "could not create remote image from ref %s", imageRef.RepoName())
}

type KeychainFactory interface {
	KeychainForImageRef(ImageRef) authn.Keychain
}

type ImageRef interface {
	ServiceAccount() string
	Namespace() string
	RepoName() string
}

type noAuthImageRef struct {
	repoName string
}

func NewNoAuthImageRef(repoName string) *noAuthImageRef {
	return &noAuthImageRef{repoName: repoName}
}

func (na *noAuthImageRef) RepoName() string {
	return na.repoName
}

func (noAuthImageRef) ServiceAccount() string {
	return ""
}

func (noAuthImageRef) Namespace() string {
	return ""
}

type RemoteImage interface {
	CreatedAt() (time.Time, error)
	Digest() (string, error)
	Label(labelName string) (string, error)
	Env(key string) (string, error)
}

//go:generate counterfeiter . RemoteImageFactory
type RemoteImageFactory interface {
	NewRemote(imageRef ImageRef) (RemoteImage, error)
}
