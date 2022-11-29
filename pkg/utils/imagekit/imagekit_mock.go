package imagekit

import (
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"io"
)

type ImgKitServiceMock struct {
	mock.Mock
}

func (i *ImgKitServiceMock) UploadFile(ctx context.Context, file io.Reader, fileName string, folder string) (*uploader.UploadResult, error) {
	args := i.Called(ctx, file, fileName, folder)
	return args.Get(0).(*uploader.UploadResult), args.Error(1)
}
