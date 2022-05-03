package logicclipping

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Clip struct {
	s3.Object
	s3Api  *AwsS3ClientAPI
	bucket string
}

func (clip *Clip) GetData() ([]byte, error) {
	return clip.s3Api.GetObject(clip.bucket, *clip.Key)
}

func (clip *Clip) Download(dir string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	data, err := clip.GetData()
	if err != nil {
		return err
	}
	file := *clip.Key
	_, file, _ = strings.Cut(file, "/")
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, file), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (lc *LogicConnection) GetClips() ([]*Clip, error) {
	clips := []*Clip{}

	s3Clips, err := lc.s3.ListObjects(lc.bucketOutputName)
	if err != nil {
		return nil, err
	}
	for _, s3Clip := range s3Clips {
		clips = append(clips, &Clip{
			Object: *s3Clip,
			s3Api:  lc.s3,
			bucket: lc.bucketOutputName,
		})
	}
	return clips, nil
}

func (lc *LogicConnection) GetClipByAssetName(assetName string) (*Clip, error) {
	allClips, err := lc.GetClips()
	if err != nil {
		return nil, err
	}
	for _, clip := range allClips {
		// checks if clip starts with assetName
		if strings.HasPrefix(*clip.Key, fmt.Sprintf("clips/%s", assetName)) {
			return clip, nil
		}
	}

	return nil, errors.New("clip not found")
}
