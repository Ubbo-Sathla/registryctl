package main

import (
	"flag"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"io"
	"os"
)

var image string

func init() {
	flag.StringVar(&image, "image", "", "docker image")
	flag.Parse()
}
func pathOpener(path string) tarball.Opener {
	return func() (io.ReadCloser, error) {
		return os.Open(path)
	}
}
func main() {
	imgTar := image
	fmt.Println("Hello world!")
	manifests, err := tarball.LoadManifest(pathOpener(imgTar))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(manifests)
	o := crane.GetOptions(crane.Insecure, crane.WithAuth(authn.FromConfig(authn.AuthConfig{
		Username: "he7",
		Password: "Q6",
	})))

	for _, manifest := range manifests {
		for _, tag := range manifest.RepoTags {

			img, err := crane.LoadTag(imgTar, tag)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(tag, img)
			newTag := fmt.Sprintf("10.12.49.246/test/%s", tag)
			ref, _ := name.ParseReference(newTag, o.Name...)

			fmt.Printf("%#v\n", ref)

			var h v1.Hash
			switch t := img.(type) {
			case v1.Image:
				fmt.Println(ref, "v1.Image", o.Remote)
				if err := remote.Write(ref, t, o.Remote...); err != nil {
					fmt.Println(err)
				}
				if h, err = t.Digest(); err != nil {
					fmt.Println(err)
				}
			case v1.ImageIndex:
				fmt.Println(ref, "v1.ImageIndex", o.Remote)
				if err := remote.WriteIndex(ref, t, o.Remote...); err != nil {
					fmt.Println(err)
				}
				if h, err = t.Digest(); err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Errorf("cannot push type (%T) to registry", img)
			}
			digest := ref.Context().Digest(h.String())
			fmt.Println(digest)
		}
	}
}
