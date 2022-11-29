package imagekit

import (
	imagekit "github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"golang.org/x/net/context"
	"io"
)

type ImgKitService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string, folder string) (*uploader.UploadResult, error)
	DeleteFile(ctx context.Context, fileId string) error
}

type ImgKitServiceImpl struct {
	ImageKitClient *imagekit.ImageKit
}

func NewImgKitService(privateKey string, publicKey string, urlEndpoint string) ImgKitService {
	return &ImgKitServiceImpl{
		ImageKitClient: imagekit.NewFromParams(imagekit.NewParams{
			PrivateKey:  privateKey,
			PublicKey:   publicKey,
			UrlEndpoint: urlEndpoint,
		}),
	}
}

func (t *ImgKitServiceImpl) UploadFile(ctx context.Context, file io.Reader, fileName string, folder string) (*uploader.UploadResult, error) {
	uploadResp, err := t.ImageKitClient.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName: fileName,
		Folder:   folder,
	})
	if err != nil {
		return nil, err
	}

	return &uploadResp.Data, nil
}

func (t *ImgKitServiceImpl) DeleteFile(ctx context.Context, fileId string) error {
	_, err := t.ImageKitClient.Media.DeleteFile(ctx, fileId)
	return err
}
