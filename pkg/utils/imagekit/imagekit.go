package imagekit

import (
	imagekit "github.com/imagekit-developer/imagekit-go"
)

type Client struct {
	privateKey     string
	publicKey      string
	URL            string
	ImageKitClient imagekit.ImageKit
}
